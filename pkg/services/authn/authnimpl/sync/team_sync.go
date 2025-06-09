package sync

import (
	"context"
	"fmt"
	"strconv"

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

	// Remove user from all teams
	s.teamService.RemoveUsersMemberships(ctx, userID)

	// Add user to teams
	for _, permission := range id.Access {
		teamUIDString := strconv.FormatInt(permission.TeamUID, 10)
		if _, err := s.teamPermissionsService.SetUserPermission(ctx, permission.OrgUID, accesscontrol.User{ID: userID}, teamUIDString, permission.Role); err != nil {
			errStr := fmt.Sprintf("failed setting permissions for user %d in team %d: %v", userID, permission.TeamUID, err.Error())
			ctxLogger.Error("Failed to add user to team", "err", errStr)
		} else {
			ctxLogger.Info("Added user to team", "permission", permission)
		}
	}

	return nil
}
