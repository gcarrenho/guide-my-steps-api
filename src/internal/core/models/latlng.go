package models

type LatLng struct {
	Latitud  float64
	Longitud float64
}

func CastToLatLng(coordinate [][]float64) []LatLng {
	var latLngList []LatLng
	for _, c := range coordinate {
		latLng := LatLng{
			Latitud:  c[0],
			Longitud: c[1],
		}
		latLngList = append(latLngList, latLng)
	}
	return latLngList
}

func AddLatLng(list *[]LatLng, latLng []LatLng) {
	*list = append(*list, latLng...)
}
