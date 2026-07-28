package services

import "gorm.io/gorm"

type ServiceManager struct {
	db *gorm.DB
}

func NewServiceManager(db *gorm.DB) *ServiceManager {
	return &ServiceManager{
		db: db,
	}
}
