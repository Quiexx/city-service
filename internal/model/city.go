package model

type CreateCityRequest struct {
	Name       string `json:"name"`
	Population int    `json:"population"`
}

type UpdateCityRequest struct {
	Name       string `json:"name"`
	Population int    `json:"population"`
}

type City struct {
	ID         string `json:"uuid"`
	Name       string `json:"name"`
	Population int    `json:"population"`
}
