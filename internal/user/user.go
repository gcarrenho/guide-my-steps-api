package user

import "strconv"

type User struct {
	ID             int64
	Name           string
	LastName       string
	UserName       string
	Gender         string
	Phone          Phone
	Height         float64 //Altura para el calculo del paso por defecto 170
	Address        Address
	FavoritePlaces []Place
	Config         Config
}

type Phone struct {
	CodeCountry string
	CodeArea    string
	Number      string
}

type Address struct {
	ID         int64
	Country    string
	City       string
	StreetName string
	Floor      int
	Door       string
	Location   LatLng
}

type LatLng struct {
	Latitud  float64
	Longitud float64
}

type Place struct {
	ID       int64
	Name     string
	Category string
	Address  Address
}

type Config struct {
	Theme    string
	Unit     string
	Language string // es-ES, en-US en-GB fr-FR it--IT de-DE
}

// //distancia/ altura*segun sexo 0.413 mujer 0.415 hombre
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

func (u User) GetLengthStep(distance float64, user User) string {
	// mapea cada unidad de medida permitida a su función de conversión
	unitConversions := map[string]func(float64) int{
		"metros":     func(d float64) int { return int(d) },
		"kilometros": func(d float64) int { return int(d / 1000) },
		"millas":     func(d float64) int { return int(d * 0.00062137) },
		"yardas":     func(d float64) int { return int(d * 1.0936) },
		"braza":      func(d float64) int { return int(d * 0.546807) },
		"pasos":      func(d float64) int { return int(int(distance * 100 / (user.Height * LengthStep[user.Gender]))) },
	}

	if converter, ok := unitConversions[user.Config.Unit]; ok {
		// si la unidad de medida del usuario se encuentra en el mapa, utiliza su función de conversión correspondiente
		return strconv.Itoa(converter(distance))
	} else {
		// si no se encuentra, calcula la longitud de paso en función de la altura y el género del usuario
		return strconv.Itoa(int(distance * 100 / (user.Height * LengthStep[user.Gender])))
	}
}
