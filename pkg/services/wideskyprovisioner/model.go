package wideskyprovisioner

import (
	"time"

	"github.com/grafana/grafana/pkg/services/team"
)

type TeamPermissions struct {
	Id 			 int64     `json:"id"`
	TeamId       int64     `json:"teamId"`
	PluginId     string    `json:"pluginId"`
	LastModified time.Time `xorm:"updated"`
	PageAccess   string    `json:"pageAccess"`
}

// ---------------------
// COMMANDS
type CreateTeamPermissionsCommand struct {
	TeamId     int64    `json:"teamId"`
	PluginId   string   `json:"pluginId"`
	PageAccess []string `json:"pageAccess"`

	Result team.Team `json:"-"`
}

type PatchTeamPermissionsCommand struct {
	Id           int64    `json:"id"`
	LastModified int64    `json:"lastModified"`
	PageAccess   []string `json:"pageAccess"`

	Result team.Team `json:"-"`
}

type DeleteTeamPermissionsCommand struct {
	Id int64 `json:"id"`

	Result team.Team `json:"-"`
}
