package db

import (
	"github.com/hongkailiu/test-go/http/model"
	"github.com/jinzhu/gorm"
)

type Service struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{db}
}

func (s *Service) GetCities(limit, offset int, cities *[]model.City) error {
	return s.db.Limit(limit).Offset(offset).Find(cities).Error
}
