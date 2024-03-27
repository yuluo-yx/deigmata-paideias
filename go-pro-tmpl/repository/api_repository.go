package repository

import "gorm.io/gorm"

type ApiRepository struct {
	ApiRepository gorm.DB
}

func (ar ApiRepository) Set() {
	
}
