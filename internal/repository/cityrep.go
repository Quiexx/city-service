package repository

import "github.com/Quiexx/city-service/internal/model"

type CityRepository interface {
	Insert(city *model.City) (string, error)
	Update(city *model.City) error
	Delete(id string) error
	GetById(id string) (*model.City, error)
}
