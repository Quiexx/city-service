package service

import (
	"github.com/Quiexx/city-service/internal/model"
	"github.com/Quiexx/city-service/internal/repository"
	"github.com/google/uuid"
)

type ICityService interface {
	Create(name string, population int) (string, error)
	Update(city *model.City) error
	Delete(id string) error
}

type MockCityService struct {
}

func NewMockCityService() ICityService {
	return &MockCityService{}
}

func (s *MockCityService) Create(name string, population int) (string, error) {
	return uuid.NewString(), nil
}

func (s *MockCityService) Update(city *model.City) error {
	return nil
}

func (s *MockCityService) Delete(id string) error {
	return nil
}

type CityService struct {
	cityRep repository.CityRepository
}

func NewCityService() ICityService {
	return &CityService{cityRep: repository.NewMemoryCityRepository()}
}

func (s *CityService) Create(name string, population int) (string, error) {
	id := uuid.NewString()
	city := model.City{ID: id, Name: name, Population: population}
	id, err := s.cityRep.Insert(&city)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *CityService) Update(city *model.City) error {
	err := s.cityRep.Update(city)
	return err
}

func (s *CityService) Delete(id string) error {
	err := s.cityRep.Delete(id)
	return err
}
