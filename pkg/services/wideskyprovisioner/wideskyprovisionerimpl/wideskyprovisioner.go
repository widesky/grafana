package wideskyprovisionerimpl

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/services/pluginsintegration/pluginstore"
	"github.com/grafana/grafana/pkg/services/wideskyprovisioner"
	"github.com/grafana/grafana/pkg/setting"
)

type Service struct {
	store  store
	Cfg    *setting.Cfg
	Logger log.Logger

	pluginStore pluginstore.Store
}

var ALL_PLUGINS_ID int64 = -1
var NO_LIMIT string = "__NO_LIMIT__"

// Check if an item is in the Array. If it is, return the index it was found at.
// If not, return -1.
func itemInArray[K comparable, V string | int64](toFind V, array []V) int {
	for i, item := range array {
		if item == toFind {
			return i
		}
	}

	return -1
}

// Slice out an item out of an Array if it exists
func sliceOutItem[K comparable, V string | int64](item V, array []V) []V {
	if index := itemInArray[V, V](item, array); index != -1 {
		if len(array) == 1 {
			return []V{}
		} else if index == 0 {
			return array[1:]
		} else if index == (len(array) - 1) {
			return array[0:index]
		} else {
			newArray := array[0:index]
			return append(newArray, array[index+1:]...)
		}
	}
	return array
}

func ProvideService(cfg *setting.Cfg, db db.DB, pluginstore pluginstore.Store) *Service {
	logger := log.New("wideskyprovisionerservice")
	logger.Info("Created WideSky Provisioner Service")

	service := &Service{
		store:  &xormStore{db: db},
		Cfg:    cfg,
		Logger: logger,

		pluginStore: pluginstore,
	}

	service.LoadTeamPluginPermissions()
	service.LoadPluginPermissions()

	return service
}

// Initialise team-based plugin permissions from the database
func (s *Service) LoadTeamPluginPermissions() {
	s.Logger.Info("Loading in saved team plugin permissions...")
	perms, err := s.store.GetByTeam(ALL_PLUGINS_ID)
	if err != nil {
		s.Logger.Error("Failed to load in team plugin permissions", "err", err)
	}

	s.Logger.Info(fmt.Sprintf("Loading in %d permissions", len(perms)))
	for _, perm := range perms {
		s.AddNewPermission(perm.PluginId, perm.TeamId, strings.Split(perm.PageAccess, ","))
	}
}

// Initialise team based access to pages of non-core plugins
func (s *Service) LoadPluginPermissions() {
	s.Logger.Info("Loading in plugin page access based on team permissions...")
	for _, p := range s.pluginStore.Plugins(context.Background()) {
		if p.IsCorePlugin() {
			continue
		}

		// Link all dashboard uid back to widesky's permission map - if it is one with perm assigned
		if _, pluginAssigned := s.Cfg.WideSkyProvisioner.TeamPermissions.Plugins[p.ID]; pluginAssigned {
			s.AddNewPageAccess(&p)
			continue
		}

		// One of those app where there is no permission defined
		s.Logger.Warn("Permission not defined for plugin", "pluginId", p.ID)
		for _, include := range p.Includes {
			if include.Type == "dashboard" && len(include.UID) > 0 {
				// Every dashboard's uid is unique so create an entry for it
				// and an empty string array to represent no one can access it except Grafana admins
				s.Cfg.WideSkyProvisioner.TeamPermissions.Dashboards[include.UID] = make([]string, 0)
			}
		}
	}
}

// Add team access to plugin
func (s *Service) AddNewPageAccess(plugin *pluginstore.Plugin) {
	pageNames := s.Cfg.WideSkyProvisioner.TeamPermissions.Plugins[plugin.ID]

	for _, include := range plugin.Includes {
		if include.Type == "dashboard" && len(include.UID) > 0 {
			var toAdd []string
			if teams, found := pageNames[NO_LIMIT]; found {
				toAdd = teams
			} else if accessTeams, found := pageNames[include.Name]; found {
				toAdd = accessTeams
			} else {
				toAdd = []string{}
			}

			s.Cfg.WideSkyProvisioner.TeamPermissions.Dashboards[include.UID] = toAdd
		}
	}
}

// Remove a team's access to a page for a plugin
func (s *Service) RemovePageAccess(plugin *pluginstore.Plugin, teamId int64) {
	for _, include := range plugin.Includes {
		if include.Type == "dashboard" && len(include.UID) > 0 {
			s.Cfg.WideSkyProvisioner.TeamPermissions.Dashboards[include.UID] =
				sliceOutItem[string](strconv.FormatInt(teamId, 10), s.Cfg.WideSkyProvisioner.TeamPermissions.Plugins[plugin.ID][include.Name])
		}
	}
}

// Add a new team-based plugin permission into memory
func (s *Service) AddNewPermission(pluginID string, teamID int64, pageAccess []string) {
	if _, found := s.Cfg.WideSkyProvisioner.TeamPermissions.Plugins[pluginID]; !found {
		s.Cfg.WideSkyProvisioner.TeamPermissions.Plugins[pluginID] = make(map[string][]string)
	}

	for _, pgName := range pageAccess {
		if pgName == "*" {
			pgName = NO_LIMIT
		}

		s.Cfg.WideSkyProvisioner.TeamPermissions.Plugins[pluginID][pgName] =
			append(s.Cfg.WideSkyProvisioner.TeamPermissions.Plugins[pluginID][pgName], strconv.FormatInt(teamID, 10))
	}
}

// Remove a team based plugin permission from memory
func (s *Service) RemovePermission(pluginId string, teamId int64, pageNames []string) {
	for pgName, teamIds := range s.Cfg.WideSkyProvisioner.TeamPermissions.Plugins[pluginId] {
		if itemInArray[string](pgName, pageNames) != -1 {
			s.Cfg.WideSkyProvisioner.TeamPermissions.Plugins[pluginId][pgName] =
				sliceOutItem[string](strconv.FormatInt(teamId, 10), teamIds)
		}
	}
}

// Add a new team-based plugin permission to the database
func (s *Service) AddTeamPluginPerm(ctx context.Context, pluginID string, teamID int64, pageAccess []string) (wideskyprovisioner.TeamPermissions, error) {
	resp, err := s.store.Create(pluginID, teamID, pageAccess)

	if err != nil {
		return resp, err
	}

	plugin, pluginExists := s.pluginStore.Plugin(ctx, pluginID)

	if !pluginExists {
		return resp, nil
	}

	s.AddNewPermission(pluginID, teamID, pageAccess)
	s.AddNewPageAccess(&plugin)
	return resp, nil
}

// Get all the plugin permissions for a given teamId
func (s *Service) GetTeamPlugins(teamID int64) ([]wideskyprovisioner.TeamPermissions, error) {
	perms, err := s.store.GetByTeam(teamID)
	formattedRes := make([]wideskyprovisioner.TeamPermissions, 0)
	if err != nil {
		return formattedRes, err
	}

	for _, perm := range perms {
		t := wideskyprovisioner.TeamPermissions{
			Id: perm.Id,
			TeamId: perm.TeamId,
			PluginId:     perm.PluginId,
			LastModified: perm.LastModified,
			PageAccess:   perm.PageAccess,
		}
		formattedRes = append(formattedRes, t)
	}

	return formattedRes, err
}

// Get a specific plugin by its id
func (s *Service) GetTeamPlugin(id int64) (wideskyprovisioner.TeamPermissions, error) {
	perm, err := s.store.GetByPerm(id)
	if err != nil {
		return wideskyprovisioner.TeamPermissions{}, err
	}

	return perm, nil
}

// Patch a specific permission
func (s *Service) PatchTeamPlugin(ctx context.Context, id int64, teamId int64, pluginId string, pageAccess []string, oldPageAccess []string) (wideskyprovisioner.TeamPermissions, error) {
	resp, err := s.store.Patch(id, teamId, pluginId, pageAccess)
	if err != nil {
		return resp, err
	}

	plugin, pluginExists := s.pluginStore.Plugin(ctx, pluginId)

	if !pluginExists {
		return resp, nil
	}

	s.RemovePermission(pluginId, teamId, oldPageAccess)
	s.RemovePageAccess(&plugin, teamId)
	s.AddNewPermission(pluginId, teamId, pageAccess)
	s.AddNewPageAccess(&plugin)

	return resp, nil
}

// Delete a specific permission
func (s *Service) DeleteTeamPlugin(ctx context.Context, id int64, teamId int64, pluginId string, pageAccess []string) error {
	err := s.store.Delete(id)
	if err != nil {
		return err
	}

	plugin, pluginExists := s.pluginStore.Plugin(ctx, pluginId)

	if !pluginExists {
		return nil
	}

	s.RemovePermission(pluginId, teamId, pageAccess)
	s.RemovePageAccess(&plugin, teamId)
	return nil
}
