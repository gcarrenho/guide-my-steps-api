package routing

import "context"

type RoutingComponent interface {
	GetRouting(ctx context.Context, myStepResquest RoutesRequest) (MySteps, error)
	//TranslateSteps(steps Step, stepType string, nextInstruction string, user appuser.User) string
}
