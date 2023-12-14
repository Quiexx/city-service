package service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Quiexx/city-service/internal/model"
)

type ThirdCityService struct {
	citiesUrl string
}

func NewThirdCityService(citiesUrl string) *ThirdCityService {
	return &ThirdCityService{citiesUrl}
}

func (s *ThirdCityService) GetById(id string) (*model.City, error) {
	utlWithId := s.citiesUrl + "/" + id
	resp, err := http.Get(utlWithId)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var city *model.City
	err = json.NewDecoder(resp.Body).Decode(&city)
	if err != nil {
		log.Fatal("Couldn't parse json:", err)
		return nil, err
	}
	return city, err
}
