package services

import (
	repositories "github.com/gcarrenho/guidemysteps/src/internal/adapters/repositories/open_street_map/models"
	"github.com/gcarrenho/guidemysteps/src/internal/core/models"
	"github.com/gcarrenho/guidemysteps/src/internal/core/ports"
)

const (
	// EarthRadius is about 6,371km according to Wikipedia
	EarthRadius = 6371
)

var modifierType = map[string]string{
	"depart_instruction":   "depart",
	"turn_instruction":     "turn",
	"continue_instruction": "turn",
	"new name_instruction": "turn",
	//"depart_pre":"depart_pre",

	"depart_pre":   "depart_pre",
	"turn_pre":     "turn_pre", //para distancia menor a 30 sino turn
	"continue_pre": "turn_pre", // para distancia menor a 30 sino turn
	"new name_pre": "turn_pre", // para distnacia menor a 30 sino turn

	"depart_alert":   "alert",
	"turn_alert":     "alert",
	"continue_alert": "alert",
	"new name_alert": "alert",

	"depart_post":   "depart_post",
	"turn_post":     "post",
	"continue_post": "post",
	"new name_post": "new_name_post",

	"end of road_instruction": "end_of_road",
	"end of road_alert":       "alert",
	"end of road_pre":         "end_of_road_pre", //si distancia es menor a 30
	"end of road_post":        "post",

	"arrive_instruction": "arrive",
	"arrive_alert":       "arrive_alert",
	"arrive_pre":         "arrive_pre",
	"arrive_post":        "post",
}

type routingSvc struct {
	routingRepo ports.RoutingRepository
	translate   ports.TranslateRepo
}

func NewRoutingSvc(routingRepo ports.RoutingRepository, translate ports.TranslateRepo) *routingSvc {
	return &routingSvc{
		routingRepo: routingRepo,
		translate:   translate,
	}
}

func (rout routingSvc) GetRouting(routesRequest models.RoutesRequest, user models.User) (*models.MySteps, error) {
	route, err := rout.routingRepo.GetRouting(routesRequest)
	mySteps := &models.MySteps{
		Version: "1.0.0",
		Status:  route.Code,
		Routes: []models.Routes{
			{
				Legs: []models.Legs{
					{},
				},
			},
		},
		Units: user.Config.Unit,
		//Waypoints: route.Waypoints,
		Language: user.Config.Language,
	}

	var allCoordinates []models.LatLng
	for _, r := range route.Routes {
		mySteps.Routes[len(mySteps.Routes)-1].Distance.Value = r.Distance // en metros
		mySteps.Routes[len(mySteps.Routes)-1].Duration.Value = r.Duration // en segundos
		allCoordinates = []models.LatLng{}

		for _, l := range r.Legs {
			for i := 0; i < len(l.Steps); i++ {

				nextInstruction := builNextInstruction(i, l.Steps, rout.translate, user)

				steps := builStep(i, l.Steps, rout.translate, user, nextInstruction)

				addNewStep(mySteps, steps)

				models.AddLatLng(&allCoordinates, models.CastToLatLng(l.Steps[i].Geometry.Coordinate))
			}
		}

		mySteps.Routes[len(mySteps.Routes)-1].Polypoints = allCoordinates
	}

	return mySteps, err
}

func builStep(i int, steps []repositories.Steps, translate ports.TranslateRepo, user models.User, nextInstruction string) models.Step {
	return models.Step{
		Distance: models.Distance{
			Value: steps[i].Distance,
		},
		Duration: models.Duration{
			Value: steps[i].Duration,
		},
		StartLocation: models.LatLngResponse{
			Latitud:  steps[i].Geometry.Coordinate[0][1],
			Longitud: steps[i].Geometry.Coordinate[0][0],
		},
		EndLocation: models.LatLngResponse{
			Latitud:  steps[i].Geometry.Coordinate[len(steps[i].Geometry.Coordinate)-1][1],
			Longitud: steps[i].Geometry.Coordinate[len(steps[i].Geometry.Coordinate)-1][0],
		},
		Intruction:                       translate.Translate(steps[i], modifierType[steps[i].Maneuver.Type+"_instruction"], nextInstruction, user),
		VerbalTransitionAlertInstruction: translate.Translate(steps[i], modifierType[steps[i].Maneuver.Type+"_alert"], nextInstruction, user),
		VerbalPreTransitionInstruction:   translate.Translate(steps[i], modifierType[steps[i].Maneuver.Type+"_pre"], nextInstruction, user),
		VerbalPostTransitionInstruction:  translate.Translate(steps[i], modifierType[steps[i].Maneuver.Type+"_post"], nextInstruction, user),
		StreetName:                       steps[i].Name,
		TravelMode:                       steps[i].Mode,
		DrivingSide:                      steps[i].DrivingSide,
		TravelType:                       "pediestran",
	}
}
func addNewStep(mySteps *models.MySteps, step models.Step) {
	mySteps.Routes[len(mySteps.Routes)-1].Legs[len(mySteps.Routes[len(mySteps.Routes)-1].Legs)-1].Steps = append(mySteps.Routes[len(mySteps.Routes)-1].Legs[len(mySteps.Routes[len(mySteps.Routes)-1].Legs)-1].Steps, step)
}

func builNextInstruction(i int, steps []repositories.Steps, translate ports.TranslateRepo, user models.User) string {
	var nextInstruction string
	if i != len(steps)-1 {
		nextInstruction = translate.Translate(steps[i+1], modifierType[steps[i+1].Maneuver.Type+"_instruction"], "", user)
	}
	return nextInstruction
}
