package models

type Place struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Address  Address `json:"address"`
}
