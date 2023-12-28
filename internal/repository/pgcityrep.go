package repository

import (
	"github.com/Quiexx/city-service/internal/model"
	"gorm.io/gorm"
)

type PGCityRepository struct {
	db *gorm.DB
}

func NewPGCityRepository(db *gorm.DB) CityRepository {
	return &PGCityRepository{db: db}
}

func (r *PGCityRepository) Insert(city *model.City) (string, error) {
	res := r.db.Create(&city)
	return city.ID, res.Error
}

func (r *PGCityRepository) Update(city *model.City) error {
	res := r.db.Save(&city)
	return res.Error
}

func (r *PGCityRepository) Delete(id string) error {
	return r.db.Delete(&model.City{}, "id = ?", id).Error
}

func (r *PGCityRepository) GetById(id string) (*model.City, error) {
	var city model.City
	res := r.db.First(&city, "id = ?", id)

	return &city, res.Error
}

func (r *PGCityRepository) GetAll() ([]*model.City, error) {
	var cities []*model.City
	res := r.db.Find(&cities)

	return cities, res.Error
}
