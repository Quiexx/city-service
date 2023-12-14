package model

type CityUpdate struct {
	Operation string `json:"operation"`
	City      City   `json:"city"`
}
