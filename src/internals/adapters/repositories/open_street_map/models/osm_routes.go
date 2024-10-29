package repositories

type OsmResponse struct {
	Code      string      `json:"code"`
	Routes    []Routes    `json:"routes"`
	Waypoints []Waypoints `json:"waypoints"`
}

type Routes struct {
	Legs       []Legs  `json:"legs"`
	WeightName string  `json:"weight_name"`
	Weight     float64 `json:"weight"`
	Duration   float64 `json:"duration"`
	Distance   float64 `json:"distance"`
}

type Legs struct {
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
