package sync

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/auth/identity"
	"github.com/grafana/grafana/pkg/services/authn"
	"github.com/grafana/grafana/pkg/services/team"
)

func ProvideTeamSync(teamService team.Service, teamPermissionsService accesscontrol.TeamPermissionsService) *TeamSync {
	return &TeamSync{teamService, teamPermissionsService, log.New("team.sync")}
}

type TeamSync struct {
	teamService            team.Service
	teamPermissionsService accesscontrol.TeamPermissionsService

	log log.Logger
}

type TeamPermission struct {
	OrgID      int64
	TeamID     int64
	Permission string
}

func (s *TeamSync) SyncTeamRolesHook(ctx context.Context, id *authn.Identity, _ *authn.Request) error {
	if !id.ClientParams.SyncTeams {
		return nil
	}

	ctxLogger := s.log.FromContext(ctx).New("id", id.ID, "login", id.Login)

	namespace, identifier := id.GetNamespacedID()
	if namespace != authn.NamespaceUser {
		ctxLogger.Warn("Failed to sync teams, invalid namespace for identity", "namespace", namespace)
		return nil
	}

	userID, err := identity.IntIdentifier(namespace, identifier)
	if err != nil {
		ctxLogger.Warn("Failed to sync teams, invalid ID for identity", "namespace", namespace, "err", err)
		return nil
	}

	// Unpack permissions
	var permissions []TeamPermission
	for _, access := range id.Access {
		accessParts := strings.SplitN(access, ":", 3)

		orgID, err := strconv.ParseInt(accessParts[0], 10, 64)
		if err != nil {
			ctxLogger.Error("Failed to parse OrgID", "error", err)
			return nil
		}

		teamID, err := strconv.ParseInt(accessParts[1], 10, 64)
		if err != nil {
			ctxLogger.Error("Failed to parse teamID", "error", err)
			return nil
		}

		teamPermission := TeamPermission{orgID, teamID, accessParts[2]}
		permissions = append(permissions, teamPermission)
	}

	// Remove user from all teams
	s.teamService.RemoveUsersMemberships(ctx, userID)

	// Add user to teams
	for _, permission := range permissions {
		teamIDString := strconv.FormatInt(permission.TeamID, 10)
		if _, err := s.teamPermissionsService.SetUserPermission(ctx, permission.OrgID, accesscontrol.User{ID: userID}, teamIDString, permission.Permission); err != nil {
			errStr := fmt.Sprintf("failed setting permissions for user %d in team %d: %v", userID, permission.TeamID, err.Error())
			ctxLogger.Error("Failed to add user to team", "err", errStr)
		} else {
			ctxLogger.Info("Added user to team", "permission", permission)
		}
	}

	return nil
}
