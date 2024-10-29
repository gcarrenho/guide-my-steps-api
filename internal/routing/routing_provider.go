package routing

import "context"

type routingProvider interface {
	GetRouting(ctx context.Context, myStepResquest RoutesRequest) (MySteps, error)
}
