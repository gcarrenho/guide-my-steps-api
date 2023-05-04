package models

type RoutesRequest struct {
	Version     string `json:"version"`
	Start       LatLng `json:"start"`
	Destination LatLng `json:"destination"`
	DrivingMode string `json:"driving_mode"`
	Language    string `json:"language"`
	UserEmail   string `json:"user_email"`
}
