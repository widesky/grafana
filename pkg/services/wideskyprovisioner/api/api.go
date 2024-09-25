package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/grafana/grafana/pkg/api/response"
	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/plugins"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	"github.com/grafana/grafana/pkg/services/pluginsintegration/pluginsettings"
	"github.com/grafana/grafana/pkg/services/pluginsintegration/pluginstore"
	"github.com/grafana/grafana/pkg/services/team"
	"github.com/grafana/grafana/pkg/services/wideskyprovisioner"
	"github.com/grafana/grafana/pkg/util"
	"github.com/grafana/grafana/pkg/web"
)

type API struct {
	teamService               team.Service
	wideSkyProvisionerService wideskyprovisioner.Service
	pluginSettings            pluginsettings.Service
	pluginStore               pluginstore.Store

	log log.Logger
}

func ProvideApi(
	routeRegister routing.RouteRegister,
	teamService team.Service,
	wideSkyProvisionerService wideskyprovisioner.Service,
	PluginSettings pluginsettings.Service,
	pluginStore pluginstore.Store,
) *API {
	api := &API{
		teamService:               teamService,
		wideSkyProvisionerService: wideSkyProvisionerService,
		pluginSettings:            PluginSettings,
		pluginStore:               pluginStore,
		log:                       log.New("wideskyprovisionerapi"),
	}

	api.log.Info("WideSky Provisioner API Registered")
	return api
}

// Using 'pluginSettings' from HTTPServer would create a circular dep
func pluginSettings(api *API, ctx context.Context, orgID int64) (map[string]*pluginsettings.InfoDTO, error) {
	pluginSettings := make(map[string]*pluginsettings.InfoDTO)

	// fill settings from database
	if pss, err := api.pluginSettings.GetPluginSettings(ctx, &pluginsettings.GetArgs{OrgID: orgID}); err != nil {
		return nil, err
	} else {
		for _, ps := range pss {
			pluginSettings[ps.PluginID] = ps
		}
	}

	// fill settings from app plugins
	for _, plugin := range api.pluginStore.Plugins(ctx, plugins.TypeApp) {
		// ignore settings that already exist
		if _, exists := pluginSettings[plugin.ID]; exists {
			continue
		}

		// add new setting which is enabled depending on if AutoEnabled: true
		pluginSetting := &pluginsettings.InfoDTO{
			PluginID:      plugin.ID,
			OrgID:         orgID,
			Enabled:       plugin.AutoEnabled,
			Pinned:        plugin.AutoEnabled,
			PluginVersion: plugin.Info.Version,
		}

		pluginSettings[plugin.ID] = pluginSetting
	}

	// fill settings from all remaining plugins (including potential app child plugins)
	for _, plugin := range api.pluginStore.Plugins(ctx) {
		// ignore settings that already exist
		if _, exists := pluginSettings[plugin.ID]; exists {
			continue
		}

		// add new setting which is enabled by default
		pluginSetting := &pluginsettings.InfoDTO{
			PluginID:      plugin.ID,
			OrgID:         orgID,
			Enabled:       true,
			Pinned:        false,
			PluginVersion: plugin.Info.Version,
		}

		// if plugin is included in an app, check app settings
		if plugin.IncludedInAppID != "" {
			// app child plugins are disabled unless app is enabled
			pluginSetting.Enabled = false
			if p, exists := pluginSettings[plugin.IncludedInAppID]; exists {
				pluginSetting.Enabled = p.Enabled
			}
		}
		pluginSettings[plugin.ID] = pluginSetting
	}

	return pluginSettings, nil
}

func isCurrentPluginId(api *API, c *contextmodel.ReqContext, pluginId string) (bool, error) {
	pluginSettingsMap, err := pluginSettings(api, c.Req.Context(), c.OrgID)
	if err != nil {
		return false, err
	}

	for _, plugin := range pluginSettingsMap {
		if pluginId == plugin.PluginID {
			return true, nil
		}
	}

	return false, nil
}

// Validate the pageAccess Array does not consist of:
// - Illegal character ','
// - Incorrect usage of wildcard '*'
// - Unnecessary space around each page
func validatePageAccess(pageAccess []string) (bool, string) {
	if len(pageAccess) == 0 {
		return false, "No pages defined for access"
	}

	hasWildCard := false
	hasPageEntry := false
	for i, pAccess := range pageAccess {
		if strings.Contains(pAccess, ",") {
			return false,
				fmt.Sprintf("Payload Array item pageAccess has illegal char ',' in the %dth element", i)
		}
		pageAccess[i] = strings.Trim(pAccess, " ")

		// Check usage of wildcard
		if hasWildCard || (pAccess == "*" && hasPageEntry) {
			return false,
				"Items in pageAccess consist of wildcard '*'. Specifying pages is not required"

		} else if pAccess == "*" {
			hasWildCard = true
		} else {
			hasPageEntry = true
		}
	}

	return true, ""
}

// swagger:route POST /api/widesky/teamPluginPerms Create a team based plugin permission
//
// Create the plugin permission for a team.
//
// Responses:
// 200: created permission
// 400: Bad request
// 401: unauthorisedError
// 403: forbiddenError
// 404: Unknown teamId
// 500: internalServerError
func (api *API) CreateTeamPluginPerm(c *contextmodel.ReqContext) response.Response {
	api.log.Info("Incoming createTeamPluginPerms")
	cmd := wideskyprovisioner.CreateTeamPermissionsCommand{}
	if err := web.Bind(c.Req, &cmd); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}

	// Check teamId is valid
	query := team.GetTeamByIDQuery{
		OrgID: c.OrgID,
		ID:    cmd.TeamId,
	}

	if _, err := api.teamService.GetTeamByID(c.Req.Context(), &query); err != nil {
		if errors.Is(err, team.ErrTeamNotFound) {
			return response.Error(http.StatusNotFound, "Team not found", err)
		}

		return response.Error(500, "Failed to get Team", err)
	}

	// ensure pluginId is one of those installed
	if isValidPlugin, err := isCurrentPluginId(api, c, cmd.PluginId); err != nil {
		return response.Error(http.StatusInternalServerError, "Failed to get list of plugins", err)
	} else if !isValidPlugin {
		return response.Error(http.StatusNotFound, "Unknown pluginId", errors.New("invalid payload"))
	}

	if isValid, pageAccessErr := validatePageAccess(cmd.PageAccess); !isValid {
		return response.Error(http.StatusBadRequest, pageAccessErr, errors.New("invalid payload"))
	}

	// ensure teamId and pluginId entry are not duplicating. Users should patch entry if adding more or less pages
	foundPerms, foundPermErr := api.wideSkyProvisionerService.GetTeamPlugins(cmd.TeamId)
	if foundPermErr != nil {
		return response.Error(500, "Failed to check existing Team plugin permissions", foundPermErr)
	}

	for _, perm := range foundPerms {
		if perm.PluginId == cmd.PluginId {
			return response.Error(
				http.StatusBadRequest,
				"A permission already exists for the given teamId and pluginId. Consider patching the existing permission",
				errors.New("invalid payload"))
		}
	}


	_, err := api.wideSkyProvisionerService.AddTeamPluginPerm(c.Req.Context(), cmd.PluginId, cmd.TeamId, cmd.PageAccess)
	if err != nil {
		return response.Error(500, "Failed to create Team", err)
	}

	return response.JSON(http.StatusOK, &util.DynMap{
		"message": "Team plugin permission created",
	})
}

// swagger:route GET /api/widesky/teamPluginPerms
//
// # Get the registered plugin permissions for a team
//
// Responses:
// 200: Valid teamId and list of permissions found
// 401: unauthorisedError
// 403: forbiddenError
// 404: teamId not found
// 500: internalServerError
func (api *API) GetTeamPluginPerms(c *contextmodel.ReqContext) response.Response {
	api.log.Info("Incoming GET team plugin permissions")

	teamId := c.QueryInt64("teamId")
	if teamId == 0 {
		return response.Error(
			http.StatusBadRequest,
			"Missing required parameter teamId or invalid teamId 0",
			errors.New("bad request"))
	}

	if teamId != -1 {
		query := team.GetTeamByIDQuery{
			OrgID: c.OrgID,
			ID:    teamId,
		}

		if _, err := api.teamService.GetTeamByID(c.Req.Context(), &query); err != nil {
			if errors.Is(err, team.ErrTeamNotFound) {
				return response.Error(http.StatusNotFound, "Team not found", err)
			}

			return response.Error(500, "Failed to get Team", err)
		}
	}

	foundPerms, err := api.wideSkyProvisionerService.GetTeamPlugins(teamId)
	if err != nil {
		return response.Error(500, "Failed to get Team plugin permissions", err)
	}

	return response.JSON(http.StatusOK, foundPerms)
}

// swagger:route PATCH /api/widesky/teamPluginPerms
//
// # Override an existing permission
//
// Responses:
// 200: Permission successfully overridden
// 401: unauthorisedError
// 403: forbiddenError
// 404: id not found
// 500: internalServerError
func (api *API) PatchTeamPluginPerms(c *contextmodel.ReqContext) response.Response {
	api.log.Info("Incoming PATCH team plugin permissions")

	cmd := wideskyprovisioner.PatchTeamPermissionsCommand{}
	if err := web.Bind(c.Req, &cmd); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}

	foundPerm, err := api.wideSkyProvisionerService.GetTeamPlugin(cmd.Id)
	if err != nil {
		return response.Error(500, "Failed to get Team-based permission", err)
	} else if time.Unix(cmd.LastModified, 0).Unix() <= foundPerm.LastModified.Unix() {
		return response.Error(
			http.StatusBadRequest,
			"patch request is not later than the latest db entry",
			errors.New("invalid payload"))
	}

	if isValid, pageAccessErr := validatePageAccess(cmd.PageAccess); !isValid {
		return response.Error(http.StatusBadRequest, pageAccessErr, errors.New("invalid payload"))
	}

	_, patchErr := api.wideSkyProvisionerService.PatchTeamPlugin(c.Req.Context(),
		cmd.Id, foundPerm.TeamId, foundPerm.PluginId, cmd.PageAccess, strings.Split(foundPerm.PageAccess, ","))
	if patchErr != nil {
		return response.Error(500, "Failed to patch Team-based permission", err)
	}

	return response.JSON(http.StatusOK, &util.DynMap{
		"message": "Team-based permission patched",
	})
}

// swagger:route DELETE /api/widesky/teamPluginPerms
//
// # Delete an existing permission
//
// Responses:
// 200: Permission successfully removed
// 401: unauthorisedError
// 403: forbiddenError
// 404: id not found
// 500: internalServerError
func (api *API) DeleteTeamPluginPerm(c *contextmodel.ReqContext) response.Response {
	api.log.Info("Incoming DELETE team plugin permissions")

	cmd := wideskyprovisioner.DeleteTeamPermissionsCommand{}
	if err := web.Bind(c.Req, &cmd); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}
	api.log.Debug("Delete with args", "id", cmd.Id)

	foundPerm, err := api.wideSkyProvisionerService.GetTeamPlugin(cmd.Id)
	if err != nil {
		return response.Error(500, "Failed to get Team-based permission", err)
	}

	deleteError := api.wideSkyProvisionerService.DeleteTeamPlugin(c.Req.Context(),
		cmd.Id, foundPerm.TeamId, foundPerm.PluginId, strings.Split(foundPerm.PageAccess, ","))
	if deleteError != nil {
		return response.Error(500, "Failed to delete Team-based permission", deleteError)
	}

	return response.JSON(http.StatusOK, &util.DynMap{
		"message": "Team-based permission removed",
	})
}
