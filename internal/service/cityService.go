package service

import (
	"github.com/Quiexx/city-service/internal/model"
	"github.com/google/uuid"
)

type ICityService interface {
	Create(name string, population int) (string, error)
	CreateWithId(id string, name string, population int) (string, error)
	Update(city *model.City) error
	Delete(id string) error
	GetById(id string) (*model.City, error)
	GetAll() ([]*model.City, error)
}

type MockCityService struct {
}

func NewMockCityService() ICityService {
	return &MockCityService{}
}

func (s *MockCityService) Create(name string, population int) (string, error) {
	return uuid.NewString(), nil
}

func (s *MockCityService) CreateWithId(id string, name string, population int) (string, error) {
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

func (s *MockCityService) GetAll() ([]*model.City, error) {
	return []*model.City{{ID: "mockID", Name: "mockName", Population: 0}}, nil
}
