package routing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gcarrenho/guidemysteps/internal/appuser"
	"github.com/gcarrenho/guidemysteps/internal/translator"
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

type OsmResponse struct {
	Code      string      `json:"code"`
	Routes    []OsmRoutes `json:"routes"`
	Waypoints []Waypoints `json:"waypoints"`
}

type OsmRoutes struct {
	Legs       []OsmLegs `json:"legs"`
	WeightName string    `json:"weight_name"`
	Weight     float64   `json:"weight"`
	Duration   float64   `json:"duration"`
	Distance   float64   `json:"distance"`
}

type OsmLegs struct {
	Steps    []Steps `json:"steps"`
	Summary  string  `json:"summary"`
	Weight   float64 `json:"weight"`
	Duration float64 `json:"duration"`
	Distance float64 `json:"distance"`
}

type Steps struct {
	Geometry      Geometry        `json:"geometry"`
	Maneuver      Maneuver        `json:"maneuver"`
	Mode          string          `json:"mode"`
	DrivingSide   string          `json:"driving_side"`
	Name          string          `json:"name"`
	Intersections []Intersections `json:"intersections"`
	Weight        float64         `json:"weight"`
	Duration      float64         `json:"duration"`
	Distance      float64         `json:"distance"`
}

type Geometry struct {
	Coordinate [][]float64 `json:"coordinates"`
	Type       string      `json:"type"`
}

type Maneuver struct {
	BearingAfter  int       `json:"bearing_after"`
	BearingBefore int       `json:"bearing_before"`
	Location      []float64 `json:"location"`
	Modifier      string    `json:"modifier"`
	Type          string    `json:"type"`
}

type Intersections struct {
	Out      int       `json:"out"`
	Entry    []bool    `json:"entry"`
	Bearings []int     `json:"bearings"`
	Location []float64 `json:"location"`
	In       int       `json:"in,omitempty"`
}

type Waypoints struct {
	Hint     string    `json:"hint"`
	Distance float64   `json:"distance"`
	Name     string    `json:"name"`
	Location []float64 `json:"location"`
}

var _ routingProvider = (*openStreetMapProvider)(nil)

type openStreetMapProvider struct {
	cli       *http.Client
	baseUrl   string
	translate *translator.TranslationService
}

func NewRoutingRepo(baseUrl string, translate *translator.TranslationService) *openStreetMapProvider {
	return &openStreetMapProvider{
		cli: &http.Client{
			Timeout: time.Second * 60,
		},
		baseUrl:   baseUrl,
		translate: translate,
	}
}

func (o *openStreetMapProvider) GetRouting(ctx context.Context, routesRequest RoutesRequest) (MySteps, error) {
	osmResponse, err := o.fetchOsmResponse(routesRequest)
	if err != nil {
		return MySteps{}, err
	}

	mySteps, err := o.convertToMySteps(osmResponse, appuser.User{})
	if err != nil {
		return MySteps{}, err
	}

	return mySteps, nil

}

func (o *openStreetMapProvider) fetchOsmResponse(routesRequest RoutesRequest) (OsmResponse, error) {
	const maxRetries = 3             // Definir un número máximo de reintentos
	var retryDelay = time.Second * 3 // Tiempo de espera entre reintentos

	var osmResponse OsmResponse
	// TODO it is mock know, but we need to recived in the request.
	routesRequest.DrivingMode = "routed-foot"
	url := o.baseUrl + "/" + routesRequest.DrivingMode + "/route/v1/driving/" +
		fmt.Sprintf("%f", routesRequest.Start.Longitud) + "," +
		fmt.Sprintf("%f", routesRequest.Start.Latitud) + ";" +
		fmt.Sprintf("%f", routesRequest.Destination.Longitud) + "," +
		fmt.Sprintf("%f", routesRequest.Destination.Latitud) +
		"?overview=false&geometries=geojson&steps=true"

	var resp *http.Response
	var err error
	fmt.Println(url)
	// Implementar un bucle de reintentos en caso de fallo
	for i := 0; i < maxRetries; i++ {
		resp, err = o.cli.Get(url)

		if err != nil {
			if errors.Is(err, io.EOF) {
				// Mostramos el error y esperamos antes de volver a intentar
				fmt.Printf("Intento %d fallido, reintentando... Error: %v\n", i, err)
				// Si se trata de un error de EOF, espera y reintenta
				time.Sleep(retryDelay)
				continue
			}
			return OsmResponse{}, err // Otros errores se devuelven directamente
		}
		if resp.StatusCode == 500 {
			fmt.Printf("Intento %d fallido, reintentando... Error: %v\n", i, err)

			time.Sleep(retryDelay)
			continue
		}
		fmt.Println("exito ", resp)
		break // Si la solicitud es exitosa, sale del bucle
	}

	if resp == nil {
		return OsmResponse{}, fmt.Errorf("respuesta nula del servidor después de %d intentos", maxRetries)
	}

	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return OsmResponse{}, err
	}

	err = json.Unmarshal(respBytes, &osmResponse)
	if err != nil {
		return OsmResponse{}, fmt.Errorf("error deserializando la respuesta: %w", err)
	}

	return osmResponse, nil
}

func (o *openStreetMapProvider) convertToMySteps(osmResponse OsmResponse, user appuser.User) (MySteps, error) {
	mySteps := MySteps{
		Version: "1.0.0",
		Status:  osmResponse.Code,
		Routes: []Route{
			{
				Legs: []Leg{
					{},
				},
			},
		},
	}

	var allCoordinates []LatLng
	for _, r := range osmResponse.Routes {
		mySteps.Routes[len(mySteps.Routes)-1].Distance.Value = r.Distance // en metros
		mySteps.Routes[len(mySteps.Routes)-1].Duration.Value = r.Duration // en segundos
		allCoordinates = []LatLng{}

		for _, l := range r.Legs {
			for i := 0; i < len(l.Steps); i++ {

				nextInstruction := builNextInstruction(i, l.Steps, o.translate, user)

				steps := builStep(i, l.Steps, o.translate, user, nextInstruction)

				addNewStep(&mySteps, steps)

				AddLatLng(&allCoordinates, CastToLatLng(l.Steps[i].Geometry.Coordinate))
			}
		}

		mySteps.Routes[len(mySteps.Routes)-1].Polypoints = allCoordinates
	}

	return mySteps, nil
}

func builStep(i int, steps []Steps, translate *translator.TranslationService, user appuser.User, nextInstruction string) Step {
	step := steps[i]

	// Construimos los parámetros para las traducciones usando `buildTransParam`
	instruction := translate.TranslateStep(buildTransParam(step, modifierType[step.Maneuver.Type+"_instruction"], nextInstruction, user))
	transitionAlert := translate.TranslateStep(buildTransParam(step, modifierType[step.Maneuver.Type+"_alert"], nextInstruction, user))
	preTransition := translate.TranslateStep(buildTransParam(step, modifierType[step.Maneuver.Type+"_pre"], nextInstruction, user))
	postTransition := translate.TranslateStep(buildTransParam(step, modifierType[step.Maneuver.Type+"_post"], nextInstruction, user))

	return Step{
		Distance: Distance{
			Value: steps[i].Distance,
		},
		Duration: Duration{
			Value: steps[i].Duration,
		},
		StartLocation: LatLngResponse{
			Latitud:  steps[i].Geometry.Coordinate[0][1],
			Longitud: steps[i].Geometry.Coordinate[0][0],
		},
		EndLocation: LatLngResponse{
			Latitud:  steps[i].Geometry.Coordinate[len(steps[i].Geometry.Coordinate)-1][1],
			Longitud: steps[i].Geometry.Coordinate[len(steps[i].Geometry.Coordinate)-1][0],
		},
		Intruction:                       instruction,
		VerbalTransitionAlertInstruction: transitionAlert,
		VerbalPreTransitionInstruction:   preTransition,
		VerbalPostTransitionInstruction:  postTransition,
		StreetName:                       steps[i].Name,
		TravelMode:                       steps[i].Mode,
		DrivingSide:                      steps[i].DrivingSide,
		TravelType:                       "pediestran",
	}
}

func addNewStep(mySteps *MySteps, step Step) {
	mySteps.Routes[len(mySteps.Routes)-1].Legs[len(mySteps.Routes[len(mySteps.Routes)-1].Legs)-1].Steps = append(mySteps.Routes[len(mySteps.Routes)-1].Legs[len(mySteps.Routes[len(mySteps.Routes)-1].Legs)-1].Steps, step)
}

func builNextInstruction(i int, steps []Steps, translate *translator.TranslationService, user appuser.User) string {
	var nextInstruction string
	if i != len(steps)-1 {
		//user.GetLengthStep(distance float64, user User)
		nextInstruction = translate.TranslateStep(buildTransParam(steps[i+1], modifierType[steps[i+1].Maneuver.Type+"_instruction"], "", user))
	}
	return nextInstruction
}

func buildTransParam(step Steps, stepType string, nextInstruction string, user appuser.User) translator.TranslatorStep {
	return translator.TranslatorStep{
		StepType:        stepType,
		Name:            step.Name,
		Distance:        step.Distance,
		Modifier:        step.Maneuver.Modifier,
		StreetName:      step.Name,
		MessageID:       step.Maneuver.Type,
		NextInstruction: nextInstruction,
		BearingAfter:    step.Maneuver.BearingAfter,
		User: translator.User{
			Config: translator.Config{
				Unit:     user.Config.Unit,
				Language: "es-ES",
			},
		},
	}
}

/*var drivingMode = map[string]string{
	"foot":             "routed-foot",
	"public-transport": "",
	"car":              "routed-car",
}*/

/*
func (r *openStreetMapProvider) GetRouting(routesRequest RoutesRequest) (MySteps, error) {
    osmResponse, err := r.fetchOsmResponse(routesRequest)
    if err != nil {
        return MySteps{}, err
    }

    mySteps, err := r.convertToMySteps(osmResponse, routesRequest.User)
    if err != nil {
        return MySteps{}, err
    }

    return mySteps, nil
}

// fetchOsmResponse maneja la comunicación con el servicio de OpenStreetMap.
func (r *openStreetMapProvider) fetchOsmResponse(routesRequest RoutesRequest) (OsmResponse, error) {
    // Lógica existente para realizar la solicitud y manejar errores.
}

// convertToMySteps convierte la respuesta de OSM a MySteps.
func (r *openStreetMapProvider) convertToMySteps(osmResponse OsmResponse, user models.User) (MySteps, error) {
    mySteps := MySteps{
        Version: "1.0.0",
        Status:  osmResponse.Code,
    }

    // Aquí va la lógica existente para llenar mySteps desde osmResponse.

    return mySteps, nil
}*/
