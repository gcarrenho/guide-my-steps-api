package models

type Address struct {
	ID         int64  `json:"id"`
	Country    string `json:"country"`
	City       string `json:"city"`
	StreetName string `json:"street_name"`
	Floor      int    `json:"floor"`
	Door       string `json:"door"`
	Location   LatLng `json:"location"`
}
