package database

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ConnectResult struct {
	DB  *gorm.DB
	Err error
}

func ConnectDb(uri string, server bool) (*gorm.DB, error) {
	lg, err := loggerFor(server)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{
		Logger: lg,
	})
	if err != nil {
		return nil, err
	}
	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(time.Hour)
	err = db.AutoMigrate(
		&User{},
		&ApplicationStatus{},
		&Company{},
		&JobApplication{},
		&CompanyChangeHistory{},
		&RefreshToken{},
		&ApplicationNote{},
		&StatusHistory{},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}
