package routing

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
