package wideskyprovisioner

import "context"

type Service interface {
	AddTeamPluginPerm(ctx context.Context, pluginID string, teamID int64, pageAccess []string) (TeamPermissions, error)
	GetTeamPlugins(teamId int64) ([]TeamPermissions, error)
	GetTeamPlugin(id int64) (TeamPermissions, error)
	PatchTeamPlugin(ctx context.Context, id int64, teamId int64, pluginId string, pageAccess []string, oldPageAccess []string) (TeamPermissions, error)
	DeleteTeamPlugin(ctx context.Context, id int64, teamId int64, pluginId string, pageAccess []string) error
}
