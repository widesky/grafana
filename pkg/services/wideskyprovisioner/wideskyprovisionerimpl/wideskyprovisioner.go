package wideskyprovisionerimpl

import (
	"context"

	"github.com/grafana/grafana/pkg/infra/db"
	contextmodel "github.com/grafana/grafana/pkg/services/contexthandler/model"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/wideskyprovisioner"
	"github.com/grafana/grafana/pkg/setting"
)

type Service struct {
	store      store
	features   featuremgmt.FeatureToggles
	primaryOrg int64
}

func ProvideService(db db.DB, features featuremgmt.FeatureToggles, settings *setting.Cfg) wideskyprovisioner.Service {
	return &Service{store: &xormStore{db: db}, features: features, primaryOrg: settings.WideSkyProvisioner.PrimaryOrg}
}

func (s *Service) GetPermissionsByTeam(ctx context.Context, teamId int64) ([]wideskyprovisioner.Permission, error) {
	return s.store.GetByTeam(ctx, teamId)
}

func (s *Service) CreatePermission(ctx context.Context, cmd *wideskyprovisioner.CreatePermissionCommand) error {
	return s.store.Create(ctx, cmd)
}

func (s *Service) UpdatePermission(ctx context.Context, permission *wideskyprovisioner.Permission) error {
	return s.store.Update(ctx, permission)
}

func (s *Service) DeletePermission(ctx context.Context, cmd *wideskyprovisioner.DeletePermissionCommand) error {
	return s.store.Delete(ctx, cmd)
}

// Determine if the team which the user resides, has access to an app's navlink
func (s *Service) WideSkyTeamHasAccess(c *contextmodel.ReqContext, permissionID string, isDashboard bool) bool {
	// Check Is Grafana admin, administrator has access to everything
	if c.SignedInUser.IsGrafanaAdmin {
		return true
	}

	// Check if the feature is enabled
	if !s.features.IsEnabledGlobally(featuremgmt.FlagWsProvisioner) {
		return true
	}

	// Check if this org is the provisioned org
	if s.primaryOrg != c.SignedInUser.GetOrgID() {
		return true
	}

	// Dashboards are not required to be checked for existence in the wider permissions context as
	// they can only be searched when using endpoints that have already gated the specific inclusion
	// of the given dashboard.
	if !isDashboard && !s.IsProvisioned(c, permissionID) {
		return true
	}

	// Check if this users set of teams is permitted to see the page
	for _, teamId := range c.SignedInUser.Teams {
		permissions, err := s.GetPermissionsByTeam(context.Background(), teamId)

		if err != nil {
			return false
		}

		for _, permission := range permissions {
			if permissionID == permission.PermissionId {
				return true
			}
		}
	}

	return false
}

func (s *Service) IsProvisioned(c *contextmodel.ReqContext, permissionID string) bool {
	allPermissions, err := s.GetPermissionsByTeam(context.Background(), -1)

	if err != nil {
		return false
	}

	for _, permission := range allPermissions {
		if permissionID == permission.PermissionId {
			return true
		}

	}

	return false
}
