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
	GetById(id string) (*model.City, error)
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

func (s *MockCityService) GetById(id string) (*model.City, error) {
	return &model.City{ID: "mockID", Name: "mockName", Population: 0}, nil

}

type CityService struct {
	cityRep repository.CityRepository
}

func NewCityService(cityRep repository.CityRepository) ICityService {
	return &CityService{cityRep: cityRep}
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

func (s *CityService) GetById(id string) (*model.City, error) {
	city, err := s.cityRep.GetById(id)
	return city, err

}
