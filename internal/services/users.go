package services

import (
	"errors"

	"github.com/filz0r/jat/internal/auth"
	"github.com/filz0r/jat/internal/database"
)

func (sm *ServiceManager) CreateUser(user database.User) (database.User, error) {
	if sm.db == nil {
		return database.User{}, errors.New("database not initialized")
	}
	password, err := auth.HashPassword(user.Password)

	if err != nil {
		return database.User{}, err
	}

	user.Password = password

	result := sm.db.Create(&user)
	return user, result.Error
}

func (sm *ServiceManager) UpdateUser(user database.User) (database.User, error) {
	if sm.db == nil {
		return database.User{}, errors.New("database not initialized")
	}
	if user.Password != "" {
		password, err := auth.HashPassword(user.Password)
		if err != nil {
			return database.User{}, err
		}
		user.Password = password
	}
	result := sm.db.Updates(&user)
	return user, result.Error
}
