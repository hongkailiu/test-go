package db

import (
	"github.com/hongkailiu/test-go/pkg/http/model"
	"github.com/jinzhu/gorm"
)

type ServiceI interface {
	GetCities(limit, offset int) (*[]model.City, error)
}

// Service provides all functions with db
type Service struct {
	db *gorm.DB
}

// New constructs a db service
func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

// GetCities get cities
func (s *Service) GetCities(limit, offset int) (*[]model.City, error) {
	var cities []model.City
	err := s.db.Limit(limit).Offset(offset).Find(&cities).Error
	return &cities, err
}
