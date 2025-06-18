package wideskyprovisioner

import (
	"context"

	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
)

type Service interface {
	GetPermissionsByTeam(ctx context.Context, teamId int64) ([]Permission, error)
	CreatePermission(ctx context.Context, cmd *CreatePermissionCommand) error
	UpdatePermission(ctx context.Context, permission *Permission) error
	DeletePermission(ctx context.Context, cmd *DeletePermissionCommand) error

	WideSkyTeamHasAccess(c *contextmodel.ReqContext, permissionID string, isDashboard bool) bool
	IsProvisioned(c *contextmodel.ReqContext, permissionID string) bool
}
