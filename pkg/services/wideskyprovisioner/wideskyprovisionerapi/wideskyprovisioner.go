package wideskyprovisionerapi

import (
	"net/http"
	"strconv"

	"github.com/grafana/grafana/pkg/api/response"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	"github.com/grafana/grafana/pkg/services/wideskyprovisioner"
	"github.com/grafana/grafana/pkg/web"
)

// swagger:route POST /api/widesky/teamPluginPerms
//
// Create the plugin permission for a team.
//
// Responses:
// 200: okResponse
// 400: badRequestError
// 500: internalServerError
func (wspapi *WideSkyProvisionerAPI) createPermission(c *contextmodel.ReqContext) response.Response {
	cmd := wideskyprovisioner.CreatePermissionCommand{}
	if err := web.Bind(c.Req, &cmd); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}

	permissions, err := wspapi.wideSkyProvisionerService.GetPermissionsByTeam(c.Req.Context(), cmd.TeamId)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "Failed to check existing permissions", err)
	}

	// Update existing permission
	for _, permission := range permissions {
		if permission.PermissionId != cmd.PermissionId {
			continue
		}

		cmd := wideskyprovisioner.Permission{
			Id:           permission.Id,
			TeamId:       cmd.TeamId,
			PermissionId: cmd.PermissionId,
			Type:         cmd.Type,
		}

		if err := wspapi.wideSkyProvisionerService.UpdatePermission(c.Req.Context(), &cmd); err != nil {
			return response.Error(http.StatusInternalServerError, "Failed to patch permission", err)
		}

		return response.Success("Permission updated")
	}

	// Create new permission
	if err := wspapi.wideSkyProvisionerService.CreatePermission(c.Req.Context(), &cmd); err != nil {
		return response.Error(http.StatusInternalServerError, "Failed to create permission", err)
	}

	return response.Success("Permission created")
}

// swagger:route GET /api/widesky/teamPluginPerms
//
// Get the registered plugin permissions for a team.
//
// Responses:
// 200: list of permissions found
// 400: badRequestError
// 500: internalServerError
func (wspapi *WideSkyProvisionerAPI) getPermissionsByTeam(c *contextmodel.ReqContext) response.Response {
	teamID, err := strconv.ParseInt(web.Params(c.Req)[":teamId"], 10, 64)
	if err != nil {
		return response.Error(http.StatusBadRequest, "teamId is invalid", err)
	}

	permissions, err := wspapi.wideSkyProvisionerService.GetPermissionsByTeam(c.Req.Context(), teamID)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "Failed to get permissions", err)
	}

	return response.JSON(http.StatusOK, permissions)
}

// swagger:route DELETE /api/widesky/teamPluginPerms
//
// Delete an existing permission.
//
// Responses:
// 200: okResponse
// 400: badRequestError
// 500: internalServerError
func (wspapi *WideSkyProvisionerAPI) deletePermission(c *contextmodel.ReqContext) response.Response {
	cmd := wideskyprovisioner.DeletePermissionCommand{}
	if err := web.Bind(c.Req, &cmd); err != nil {
		return response.Error(http.StatusBadRequest, "bad request data", err)
	}

	if err := wspapi.wideSkyProvisionerService.DeletePermission(c.Req.Context(), &cmd); err != nil {
		return response.Error(http.StatusInternalServerError, "Failed to delete permission", err)
	}

	return response.Success("Permission removed")
}
