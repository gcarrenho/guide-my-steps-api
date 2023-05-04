package models

type MySteps struct {
	Version   string   `json:"version"`
	Status    string   `json:"status"`
	Routes    []Routes `json:"routes"`
	Units     string   `json:"units"`
	Waypoints string   `json:"waypoints"`
	Language  string   `json:"language"`
}

type Routes struct {
	Legs       []Legs   `json:"legs"`
	Polypoints []LatLng `json:"polypoints"`
	Duration   Duration `json:"duration"`
	Distance   Distance `json:"distance"`
}

type Legs struct {
	Steps   []Step `json:"steps"`
	Summary string `json:"summary"`
}

type Duration struct {
	Value float64 `json:"value"`
	Text  string  `json:"text"`
}

type Distance struct {
	Value float64 `json:"value"`
	Text  string  `json:"text"`
}

type Step struct {
	StartLocation                    LatLng   `json:"start_location"`
	EndLocation                      LatLng   `json:"end_location"`
	Duration                         Duration `json:"duration"`
	Distance                         Distance `json:"distance"`
	Intruction                       string   `json:"intruction"`
	VerbalTransitionAlertInstruction string   `json:"verbal_transition_alert_instruction,omitempty"`
	VerbalPreTransitionInstruction   string   `json:"verbal_pre_transition_instruction"`
	VerbalPostTransitionInstruction  string   `json:"verbal_post_transition_instruction"`
	TravelMode                       string   `json:"travel_mode"`
	TravelType                       string   `json:"travel_type"`
	DrivingSide                      string   `json:"driving_side"`
	StreetName                       string   `json:"street_name"`
}

/*func BuilVerbalPostTransitionInstruction(t string, d float64) string {

	str := "Continue for "
	if t == "arrive" {
		str = ""
	} else {
		pasos := strconv.Itoa(int(d / 0.762))
		if true { // config es pasos
			str += pasos + " steps"
		}
	}

	return str
}*/
