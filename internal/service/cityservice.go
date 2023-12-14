package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Quiexx/city-service/internal/model"
	"github.com/Quiexx/city-service/internal/repository"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
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

type FirstCityService struct {
	kafkaWriter *kafka.Writer
	cityRep     repository.CityRepository
}

func NewFirstCityService(cityRep repository.CityRepository, kafkaWriter *kafka.Writer) ICityService {
	return &FirstCityService{cityRep: cityRep, kafkaWriter: kafkaWriter}
}

func (s *FirstCityService) Create(name string, population int) (string, error) {
	id := uuid.NewString()
	city := model.City{ID: id, Name: name, Population: population}
	id, err := s.cityRep.Insert(&city)
	if err != nil {
		return "", err
	}
	s.SendToKafka(&city, "INSERT")
	return id, nil
}

func (s *FirstCityService) Update(city *model.City) error {
	err := s.cityRep.Update(city)
	s.SendToKafka(city, "UPDATE")
	return err
}

func (s *FirstCityService) Delete(id string) error {
	city, err := s.GetById(id)

	if err != nil {
		return err
	}

	err = s.cityRep.Delete(id)
	s.SendToKafka(city, "DELETE")
	return err
}

func (s *FirstCityService) GetById(id string) (*model.City, error) {
	city, err := s.cityRep.GetById(id)
	return city, err

}

func (s *FirstCityService) SendToKafka(city *model.City, operation string) {
	mesId := uuid.NewString()
	update := model.CityUpdate{
		Operation: operation,
		City:      *city,
	}
	payload, err := json.Marshal(update)

	if err != nil {
		log.Fatal(err)
	}

	err = s.kafkaWriter.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(mesId),
			Value: []byte(payload),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
