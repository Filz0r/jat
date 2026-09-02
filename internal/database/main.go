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

	err = db.Callback().Query().Before("gorm:query").Register("auto_preload_user_settings", func(tx *gorm.DB) {
		if tx.Error != nil {
			return
		}
		switch tx.Statement.Model.(type) {
		case User, *User, []User, *[]User:
			if tx.Statement.Preloads == nil {
				tx.Statement.Preloads = make(map[string][]interface{})
			}
			if _, ok := tx.Statement.Preloads["UserSettings"]; !ok {
				tx.Statement.Preloads["UserSettings"] = nil
			}
		}
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
		&UserSettings{},
		&BanList{},
		&ApplicationStatus{},
		&Company{},
		&JobApplication{},
		&CompanyChangeHistory{},
		&RefreshToken{},
		&ApplicationNote{},
		&StatusHistory{},
		&Config{},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}
