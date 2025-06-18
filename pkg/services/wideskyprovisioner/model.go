package wideskyprovisioner

type Permission struct {
	Id           int64  `json:"id" xorm:"pk autoincr 'id'"`
	TeamId       int64  `json:"teamId" xorm:"team_id"`
	PermissionId string `json:"permissionId" xorm:"permission_id"`
	Type         string `json:"type" xorm:"type"`
}

// ---------------------
// COMMANDS
type CreatePermissionCommand struct {
	TeamId       int64  `json:"teamId"`
	PermissionId string `json:"permissionId"`
	Type         string `json:"type"`
}

type DeletePermissionCommand struct {
	Id int64 `json:"id"`
}
