package ports

import (
	osmmodels "github.com/gcarrenho/guidemysteps/src/internals/adapters/repositories/open_street_map/models"
	"github.com/gcarrenho/guidemysteps/src/internals/core/models"
)

type RoutingSvc interface {
	GetRouting(routesRequest models.RoutesRequest, user models.User) (*models.MySteps, error)
}

type RoutingRepository interface {
	GetRouting(routesRequest models.RoutesRequest) (*osmmodels.OsmResponse, error)
}
