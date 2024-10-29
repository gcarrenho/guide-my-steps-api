package models

type Phone struct {
	CodeCountry string `json:"code_country"`
	CodeArea    string `json:"code_area"`
	Number      string `json:"number"`
}
