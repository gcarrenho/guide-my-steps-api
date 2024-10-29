package routinggrpc

import (
	"context"

	api "github.com/gcarrenho/guidemysteps/api/v1/mysteps"

	"github.com/gcarrenho/guidemysteps/internal/routing"
	"google.golang.org/grpc"
)

var _ api.MyStepsServer = (*GRPCServer)(nil)

type Config struct {
	FeatureSvc routing.RoutingComponent
}

type GRPCServer struct {
	api.UnimplementedMyStepsServer
	*Config
}

func newGRPCServer(config *Config) (srv *GRPCServer, err error) {
	srv = &GRPCServer{
		Config: config,
	}

	return srv, nil
}

func NewGRPCServer(config *Config, grpcOpts ...grpc.ServerOption) (
	*grpc.Server,
	error,
) {
	gsrv := grpc.NewServer(grpcOpts...) // Create the server gRPC with grpcOpts

	srv, err := newGRPCServer(config) // Create an instance of gRPC server with config
	if err != nil {
		return nil, err
	}
	api.RegisterMyStepsServer(gsrv, srv) // Registers  the our server(srv) in the gRPC server under th "LogServer" services defined in api.

	return gsrv, nil
}

func (s *GRPCServer) GetRoute(ctx context.Context, request *api.MyStepsRequest) (*api.MyStepsResponse, error) {
	routesModel := routing.RoutesRequest{
		Start: routing.LatLng{
			Latitud:  request.Start.Latitude,
			Longitud: request.Start.Longitude,
		},
		Destination: routing.LatLng{
			Latitud:  request.Destination.Latitude,
			Longitud: request.Destination.Longitude,
		},
		DrivingMode: request.DrivingMode,
		Language:    request.Language,
		UserEmail:   request.UserEmail,
	}

	mySteps, err := s.Config.FeatureSvc.GetRouting(ctx, routesModel)
	if err != nil {
		return nil, err
	}

	return mySteps.ConvertToProtoResponse(), nil
}

/*
func convertToProtoResponse(mySteps routing.MySteps) *api.MyStepsResponse {
	// Crear y mapear la lista de rutas (Routes)
	routesProto := make([]*api.Route, len(mySteps.Routes))
	for i, r := range mySteps.Routes {
		legsProto := make([]*api.Leg, len(r.Legs))
		for j, l := range r.Legs {
			stepsProto := make([]*api.Step, len(l.Steps))
			for k, s := range l.Steps {
				stepsProto[k] = &api.Step{
					StartLocation:                    convertLatLngResponseToProto(s.StartLocation),
					EndLocation:                      convertLatLngResponseToProto(s.EndLocation),
					Duration:                         convertDurationToProto(s.Duration),
					Distance:                         convertDistanceToProto(s.Distance),
					Intruction:                       s.Intruction,
					VerbalTransitionAlertInstruction: s.VerbalTransitionAlertInstruction,
					VerbalPreTransitionInstruction:   s.VerbalPreTransitionInstruction,
					VerbalPostTransitionInstruction:  s.VerbalPostTransitionInstruction,
					TravelMode:                       s.TravelMode,
					TravelType:                       s.TravelType,
					DrivingSide:                      s.DrivingSide,
					StreetName:                       s.StreetName,
				}
			}

			legsProto[j] = &api.Leg{
				Steps:   stepsProto,
				Summary: l.Summary,
			}
		}

		polyPointsProto := make([]*api.LatLng, len(r.Polypoints))
		for m, p := range r.Polypoints {
			polyPointsProto[m] = convertLatLngToProto(p)
		}

		routesProto[i] = &api.Route{
			Legs:       legsProto,
			Polypoints: polyPointsProto,
			Duration:   convertDurationToProto(r.Duration),
			Distance:   convertDistanceToProto(r.Distance),
		}
	}

	return &api.MyStepsResponse{
		//Version:   mySteps.Version,
		Status:    mySteps.Status,
		Routes:    routesProto,
		Units:     mySteps.Units,
		Waypoints: mySteps.Waypoints,
		Language:  mySteps.Language,
	}
}

func convertLatLngToProto(latlng routing.LatLng) *api.LatLng {
	return &api.LatLng{
		Latitude:  latlng.Latitud,
		Longitude: latlng.Longitud,
	}
}

func convertLatLngResponseToProto(latlngResp routing.LatLngResponse) *api.LatLng {
	return &api.LatLng{
		Latitude:  latlngResp.Latitud,
		Longitude: latlngResp.Longitud,
	}
}

func convertDurationToProto(duration routing.Duration) *api.Duration {
	return &api.Duration{
		Value: duration.Value,
		Text:  duration.Text,
	}
}

func convertDistanceToProto(distance routing.Distance) *api.Distance {
	return &api.Distance{
		Value: distance.Value,
		Text:  distance.Text,
	}
}*/

/*=======================================================================================================================================*/
/*SECOND VERSION WITH GORUTINES*/