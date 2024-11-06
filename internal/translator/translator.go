package translator

type TranslatorStep struct {
	StepType        string
	Name            string
	Distance        float64
	Modifier        string
	StreetName      string
	MessageID       string
	NextInstruction string
	BearingAfter    int
	User
}

type User struct {
	/*ID             int64   `json:"-"`
	Name           string  `json:"name"`
	LastName       string  `json:"last_name"`
	UserName       string  `json:"user_name"`
	Gender         string  `json:"gender"`
	Phone          Phone   `json:"phone"`
	Height         float64 `json:"height"` //Altura para el calculo del paso por defecto 170
	Address        Address `json:"address"`
	FavoritePlaces []Place `json:"favorite_places"`*/
	Config     Config
	LengthStep string
}

type Config struct {
	//Theme    string `json:"theme"`
	Unit     string
	Language string // es-ES, en-US en-GB fr-FR it--IT de-DE
}
