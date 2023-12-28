package repository

import (
	"fmt"

	"github.com/Quiexx/city-service/internal/model"
)

type MemoryCityRepository struct {
	data map[string]*model.City
}

func NewMemoryCityRepository() CityRepository {
	return &MemoryCityRepository{data: make(map[string]*model.City)}
}

func (r *MemoryCityRepository) Insert(city *model.City) (string, error) {
	r.data[city.ID] = city
	return city.ID, nil
}

func (r *MemoryCityRepository) Update(city *model.City) error {
	_, ok := r.data[city.ID]
	if !ok {
		return fmt.Errorf("city with id %s not found", city.ID)
	}
	r.data[city.ID] = city
	return nil
}

func (r *MemoryCityRepository) Delete(id string) error {
	delete(r.data, id)
	return nil
}

func (r *MemoryCityRepository) GetById(id string) (*model.City, error) {
	city, ok := r.data[id]
	if !ok {
		return nil, fmt.Errorf("city with id %s not found", id)
	}

	return city, nil
}

func (r *MemoryCityRepository) GetAll() ([]*model.City, error) {
	cities := make([]*model.City, 0, len(r.data))

	for _, city := range r.data {
		cities = append(cities, city)
	}

	return cities, nil
}
