package wideskyprovisionerimpl

import (
	"context"
	"strings"
	"time"

	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/services/wideskyprovisioner"
)

type store interface {
	Create(pluginID string, teamID int64, pageAccess []string) (wideskyprovisioner.TeamPermissions, error)
	GetByTeam(teamId int64) ([]wideskyprovisioner.TeamPermissions, error)
	GetByPerm(id int64) (wideskyprovisioner.TeamPermissions, error)
	Patch(id int64, teamID int64, pluginID string, pageAccess []string) (wideskyprovisioner.TeamPermissions, error)
	Delete(id int64) error
}

type xormStore struct {
	db db.DB
}

// Add a new team-based plugin permission entry
func (ss *xormStore) Create(pluginID string, teamID int64, pageAccess []string) (result wideskyprovisioner.TeamPermissions, err error) {
	err = ss.db.WithDbSession(context.Background(), func(dbSession *db.Session) error {
		session := dbSession.Table("ws_team_plugin_permissions")

		teamPluginPerm := wideskyprovisioner.TeamPermissions{
			TeamId:       teamID,
			PluginId:     pluginID,
			LastModified: time.Now(),
			PageAccess:   strings.Join(pageAccess, ","),
		}

		if _, err := session.Insert(&teamPluginPerm); err != nil {
			return err
		}

		result = teamPluginPerm
		return nil
	})

	return result, err
}

// Get all the stored plugin permission for a given teamId. Using teamId -1 will get all
// plugin permissions.
func (ss *xormStore) GetByTeam(teamId int64) ([]wideskyprovisioner.TeamPermissions, error) {
	result := make([]wideskyprovisioner.TeamPermissions, 0)

	err := ss.db.WithDbSession(context.Background(), func(dbSession *db.Session) error {
		sess := dbSession.Table("ws_team_plugin_permissions")
		sess.Cols(
			"id",
			"team_id",
			"plugin_id",
			"last_modified",
			"page_access",
		)
		// for internal use only to grab all perms
		if teamId != -1 {
			sess.Where("team_id = ?", teamId)
		}
		sess.Asc("team_id")

		if err := sess.Find(&result); err != nil {
			return err
		}
		return nil
	})

	return result, err
}

// Get a specific permission by its id
func (ss *xormStore) GetByPerm(id int64) (wideskyprovisioner.TeamPermissions, error) {
	result := make([]wideskyprovisioner.TeamPermissions, 0)
	err := ss.db.WithDbSession(context.Background(), func(dbSession *db.Session) error {
		sess := dbSession.Table("ws_team_plugin_permissions")
		sess.Cols(
			"id",
			"team_id",
			"plugin_id",
			"last_modified",
			"page_access",
		)
		sess.Where("id = ?", id)
		if err := sess.Find(&result); err != nil {
			return err
		}
		return nil
	})

	if len(result) > 0 {
		return result[0], err
	} else {
		return wideskyprovisioner.TeamPermissions{}, err
	}
}

// Overwrite a specific id with the new changes
func (ss *xormStore) Patch(id int64, teamID int64, pluginID string, pageAccess []string) (wideskyprovisioner.TeamPermissions, error) {
	teamPluginPerm := wideskyprovisioner.TeamPermissions{
		Id:            id,
		TeamId:       teamID,
		PluginId:     pluginID,
		LastModified: time.Now(),
		PageAccess:   strings.Join(pageAccess, ","),
	}
	err := ss.db.WithDbSession(context.Background(), func(dbSession *db.Session) error {
		sess := dbSession.Table("ws_team_plugin_permissions").
			ID(id).
			AllCols()
		if _, err := sess.Update(&teamPluginPerm); err != nil {
			return err
		}
		return nil
	})

	return teamPluginPerm, err
}

// Remove a specific permission from the database
func (ss *xormStore) Delete(id int64) error {
	return ss.db.WithDbSession(context.Background(), func(dbSession *db.Session) error {
		sess := dbSession.Table("ws_team_plugin_permissions").ID(id)
		if _, err := sess.Delete(&wideskyprovisioner.TeamPermissions{}); err != nil {
			return err
		}
		return nil
	})
}
