package ports

import (
	osmmodels "project/guidemysteps/src/internal/adapters/repositories/open_street_map/models"
	"project/guidemysteps/src/internal/core/models"
)

type RoutingSvc interface {
	GetRouting(routesRequest models.RoutesRequest, user models.User) (*models.MySteps, error)
}

type RoutingRepository interface {
	GetRouting(routesRequest models.RoutesRequest) (*osmmodels.OsmResponse, error)
}
