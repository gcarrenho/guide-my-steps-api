package models

type RoutesRequest struct {
	Version     string `json:"version"`
	Start       LatLng `json:"start"`
	Destination LatLng `json:"destination"`
	DrivingMode string `json:"driving_mode"`
	Language    string `json:"language"`
	UserEmail   string `json:"user_email"`
}

type LatLng struct {
	Latitud  float64 `json:"latitude"`
	Longitud float64 `json:"longitude"`
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

func AddLatLng(list *[]LatLng, latLng []LatLng) {
	*list = append(*list, latLng...)
}
