package wideskyprovisionerapi

import (
	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/middleware"
	"github.com/grafana/grafana/pkg/services/wideskyprovisioner"
)

type WideSkyProvisionerAPI struct {
	wideSkyProvisionerService wideskyprovisioner.Service
}

func ProvideWideSkyProvisionerAPI(
	routeRegister routing.RouteRegister,
	wideSkyProvisionerService wideskyprovisioner.Service,
) *WideSkyProvisionerAPI {
	wspapi := &WideSkyProvisionerAPI{
		wideSkyProvisionerService: wideSkyProvisionerService,
	}

	wspapi.registerRoutes(routeRegister)
	return wspapi
}

func (wspapi *WideSkyProvisionerAPI) registerRoutes(router routing.RouteRegister) {
	router.Group("/api", func(apiRoute routing.RouteRegister) {
		apiRoute.Group("/widesky", func(wideSkyRoute routing.RouteRegister) {
			wideSkyRoute.Post("/teamPluginPerms", middleware.ReqOrgAdmin, routing.Wrap(wspapi.createPermission))
			wideSkyRoute.Get("/teamPluginPerms/:teamId", middleware.ReqOrgAdmin, routing.Wrap(wspapi.getPermissionsByTeam))
			wideSkyRoute.Delete("/teamPluginPerms", middleware.ReqOrgAdmin, routing.Wrap(wspapi.deletePermission))
		}, middleware.ReqSignedIn)
	})
}
