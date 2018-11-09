package model

import "github.com/jinzhu/gorm"

type Order struct {
	gorm.Model
	Name   string
	CityID int
}

type City struct {
	gorm.Model
	Name string
}
