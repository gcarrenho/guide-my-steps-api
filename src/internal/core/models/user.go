package models

import "strconv"

type User struct {
	ID             int64   `json:"-"`
	Name           string  `json:"name"`
	LastName       string  `json:"last_name"`
	UserName       string  `json:"user_name"`
	Gender         string  `json:"gender"`
	Phone          Phone   `json:"phone"`
	Height         float64 `json:"height"` //Altura para el calculo del paso por defecto 170
	Address        Address `json:"address"`
	FavoritePlaces []Place `json:"favorite_places"`
	Config         Config  `json:"config"`
}

type Config struct {
	Theme    string `json:"theme"`
	Unit     string `json:"unit"`
	Language string `json:"language"` // es-ES, en-US en-GB fr-FR it--IT de-DE
}

////distancia/ altura*segun sexo 0.413 mujer 0.415 hombre
var LengthStep = map[string]float64{
	"m":  0.435,
	"f":  0.413,
	"sg": 0.435,
}

func NewUser(user User) User {
	return User{
		Name:     "",
		UserName: user.UserName, //email
		LastName: "",
		Gender:   "sg",
		Phone:    Phone{},
		Height:   170,
		Address: Address{
			StreetName: "",
			Location: LatLng{
				Latitud:  0,
				Longitud: 0,
			},
		},
		FavoritePlaces: []Place{},
		Config: Config{
			Theme:    "default",
			Unit:     "pasos",
			Language: "es-ES",
		},
	}
}

func GetLenghtStep(distance float64, user User) string {
	if user.Config.Unit == "metros" {
		return strconv.Itoa(int(distance))
	}

	return strconv.Itoa(int(distance * 100 / (user.Height * LengthStep[user.Gender])))
}
