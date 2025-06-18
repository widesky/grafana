package wideskyprovisionerimpl

import (
	"context"

	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/services/wideskyprovisioner"
)

type store interface {
	GetByTeam(ctx context.Context, teamId int64) ([]wideskyprovisioner.Permission, error)
	Create(ctx context.Context, cmd *wideskyprovisioner.CreatePermissionCommand) error
	Update(ctx context.Context, permission *wideskyprovisioner.Permission) error
	Delete(ctx context.Context, cmd *wideskyprovisioner.DeletePermissionCommand) error
}

type xormStore struct {
	db db.DB
}

// Get all the stored plugin permission for a given teamId. Using teamId -1 will get all plugin permissions.
func (ss *xormStore) GetByTeam(ctx context.Context, teamId int64) ([]wideskyprovisioner.Permission, error) {
	permissions := []wideskyprovisioner.Permission{}

	err := ss.db.WithDbSession(ctx, func(sess *db.Session) error {
		query := sess.Table("widesky_permissions").
			Cols("id", "team_id", "permission_id", "type").
			Asc("team_id")

		if teamId != -1 {
			query = query.Where("team_id = ?", teamId)
		}

		return query.Find(&permissions)
	})

	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// Add a new team-based plugin permission entry
func (ss *xormStore) Create(ctx context.Context, cmd *wideskyprovisioner.CreatePermissionCommand) error {
	return ss.db.WithDbSession(ctx, func(sess *db.Session) error {
		_, err := sess.Table("widesky_permissions").Insert(&wideskyprovisioner.Permission{
			TeamId:       cmd.TeamId,
			PermissionId: cmd.PermissionId,
			Type:         cmd.Type,
		})

		if err != nil {
			return err
		}

		return nil
	})
}

// Overwrite a specific id with the new changes
func (ss *xormStore) Update(ctx context.Context, permission *wideskyprovisioner.Permission) error {
	return ss.db.WithDbSession(ctx, func(sess *db.Session) error {
		update := map[string]any{"team_id": permission.TeamId, "permission_id": permission.PermissionId, "type": permission.Type}
		_, err := sess.Table("widesky_permissions").
			Where("id = ?", permission.Id).
			AllCols().
			Update(&update)

		if err != nil {
			return err
		}

		return nil
	})
}

// Remove a specific permission from the database
func (ss *xormStore) Delete(ctx context.Context, cmd *wideskyprovisioner.DeletePermissionCommand) error {
	return ss.db.WithDbSession(ctx, func(sess *db.Session) error {
		_, err := sess.Table("widesky_permissions").
			Where("id = ?", cmd.Id).Delete(&wideskyprovisioner.Permission{})

		if err != nil {
			return err
		}

		return nil
	})
}
