package routing

import (
	"sync"

	api "github.com/gcarrenho/guidemysteps/api/v1/mysteps"
)

type RoutesRequest struct {
	Version     string
	Start       LatLng
	Destination LatLng
	DrivingMode string
	Language    string
	UserEmail   string
}

type LatLng struct {
	Latitud  float64
	Longitud float64
}

type MySteps struct {
	Version   string
	Status    string
	Routes    []Route
	Units     string
	Waypoints string
	Language  string
}

type Route struct {
	Legs       []Leg
	Polypoints []LatLng
	Duration   Duration
	Distance   Distance
}

type Leg struct {
	Steps   []Step
	Summary string
}

type Duration struct {
	Value float64
	Text  string
}

type Distance struct {
	Value float64
	Text  string
}

type Step struct {
	StartLocation                    LatLngResponse
	EndLocation                      LatLngResponse
	Duration                         Duration
	Distance                         Distance
	Intruction                       string
	VerbalTransitionAlertInstruction string
	VerbalPreTransitionInstruction   string
	VerbalPostTransitionInstruction  string
	TravelMode                       string
	TravelType                       string
	DrivingSide                      string
	StreetName                       string
}

type LatLngResponse struct {
	Latitud  float64
	Longitud float64
}

func AddLatLng(list *[]LatLng, latLng []LatLng) {
	*list = append(*list, latLng...)
}

func CastToLatLng(coordinate [][]float64) []LatLng {
	var latLngList []LatLng
	for _, c := range coordinate {
		latLng := LatLng{
			Latitud:  c[1],
			Longitud: c[0],
		}
		latLngList = append(latLngList, latLng)
	}
	return latLngList
}

func (ms MySteps) ConvertToProtoResponse() *api.MyStepsResponse {
	routesProto := make([]*api.Route, len(ms.Routes))
	var wg sync.WaitGroup
	wg.Add(len(ms.Routes))

	for i, r := range ms.Routes {
		go func(i int, r Route) {
			defer wg.Done()
			routesProto[i] = convertRouteToProto(r)
		}(i, r)
	}

	wg.Wait()

	return &api.MyStepsResponse{
		//Version:   mySteps.Version,
		Status:    ms.Status,
		Routes:    routesProto,
		Units:     ms.Units,
		Waypoints: ms.Waypoints,
		Language:  ms.Language,
	}
}

func convertRouteToProto(route Route) *api.Route {
	legsProto := make([]*api.Leg, len(route.Legs))
	for j, l := range route.Legs {
		legsProto[j] = convertLegToProto(l)
	}

	polyPointsProto := make([]*api.LatLng, len(route.Polypoints))
	for m, p := range route.Polypoints {
		polyPointsProto[m] = convertLatLngToProto(p)
	}

	return &api.Route{
		Legs:       legsProto,
		Polypoints: polyPointsProto,
		Duration:   convertDurationToProto(route.Duration),
		Distance:   convertDistanceToProto(route.Distance),
	}
}

func convertLegToProto(leg Leg) *api.Leg {
	stepsProto := make([]*api.Step, len(leg.Steps))
	for i, step := range leg.Steps {
		stepsProto[i] = convertStepToProto(step)
	}

	return &api.Leg{
		Steps:   stepsProto,
		Summary: leg.Summary,
	}
}

func convertStepToProto(step Step) *api.Step {
	return &api.Step{
		StartLocation:                    convertLatLngToProto(LatLng(step.StartLocation)),
		EndLocation:                      convertLatLngToProto(LatLng(step.EndLocation)),
		Duration:                         convertDurationToProto(step.Duration),
		Distance:                         convertDistanceToProto(step.Distance),
		Intruction:                       step.Intruction,
		VerbalTransitionAlertInstruction: step.VerbalTransitionAlertInstruction,
		VerbalPreTransitionInstruction:   step.VerbalPreTransitionInstruction,
		VerbalPostTransitionInstruction:  step.VerbalPostTransitionInstruction,
		TravelMode:                       step.TravelMode,
		TravelType:                       step.TravelType,
		DrivingSide:                      step.DrivingSide,
		StreetName:                       step.StreetName,
	}
}

func convertLatLngToProto(latlng LatLng) *api.LatLng {
	return &api.LatLng{
		Latitude:  latlng.Latitud,
		Longitude: latlng.Longitud,
	}
}

func convertDurationToProto(duration Duration) *api.Duration {
	return &api.Duration{
		Value: duration.Value,
		Text:  duration.Text,
	}
}

func convertDistanceToProto(distance Distance) *api.Distance {
	return &api.Distance{
		Value: distance.Value,
		Text:  distance.Text,
	}
}
