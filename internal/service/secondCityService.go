package service

import (
	"github.com/Quiexx/city-service/internal/model"
	"github.com/Quiexx/city-service/internal/repository"
	"github.com/google/uuid"
)

type SecondCityService struct {
	cityRep repository.CityRepository
}

func NewSecondCityService(cityRep repository.CityRepository) ICityService {
	return &SecondCityService{cityRep: cityRep}
}

func (s *SecondCityService) Create(name string, population int) (string, error) {
	id := uuid.NewString()
	city := model.City{ID: id, Name: name, Population: population}
	id, err := s.cityRep.Insert(&city)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *SecondCityService) CreateWithId(id string, name string, population int) (string, error) {
	city := model.City{ID: id, Name: name, Population: population}
	id, err := s.cityRep.Insert(&city)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *SecondCityService) Update(city *model.City) error {
	err := s.cityRep.Update(city)
	return err
}

func (s *SecondCityService) Delete(id string) error {
	_, err := s.GetById(id)

	if err != nil {
		return err
	}

	err = s.cityRep.Delete(id)
	return err
}

func (s *SecondCityService) GetById(id string) (*model.City, error) {
	city, err := s.cityRep.GetById(id)
	return city, err
}

func (s *SecondCityService) GetAll() ([]*model.City, error) {
	cities, err := s.cityRep.GetAll()
	return cities, err
}
