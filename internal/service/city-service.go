package service

import (
	"github.com/Quiexx/city-service/internal/model"
	"github.com/google/uuid"
)

type ICityService interface {
	Create(name string, population int) (string, error)
	Update(city *model.City) error
	Delete(uuid string) error
}

type MockCityService struct {
}

func (s *MockCityService) Create(name string, population int) (string, error) {
	return uuid.NewString(), nil
}

func (s *MockCityService) Update(city *model.City) error {
	return nil
}

func (s *MockCityService) Delete(uuid string) error {
	return nil
}
