package services

import (
	"errors"

	"github.com/filz0r/jat/internal/auth"
	"github.com/filz0r/jat/internal/database"
	"github.com/google/uuid"
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
	result := sm.db.Save(&user)
	return user, result.Error
}

func (sm *ServiceManager) GetUserByEmail(email string) (database.User, error) {
	if sm.db == nil {
		return database.User{}, errors.New("database not initialized")
	}
	user := database.User{}
	result := sm.db.Where("email = ?", email).First(&user)
	return user, result.Error
}

func (sm *ServiceManager) GetUserByID(id uuid.UUID) (database.User, error) {
	if sm.db == nil {
		return database.User{}, errors.New("database not initialized")
	}
	user := database.User{}
	result := sm.db.Where("id = ?", id).First(&user)
	return user, result.Error
}

func (sm *ServiceManager) IsUserAdmin(id uuid.UUID) bool {
	if sm.db == nil {
		return false
	}
	user := database.User{}
	result := sm.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return false
	}
	return user.IsAdmin
}

func (sm *ServiceManager) GetAllUsers() ([]database.User, error) {
	if sm.db == nil {
		return []database.User{}, errors.New("database not initialized")
	}
	var users []database.User
	result := sm.db.Find(&users)
	return users, result.Error
}

func (sm *ServiceManager) SetUserDefaultApplicationStatus(id uuid.UUID, statusID uint) error {
	if sm.db == nil {
		return errors.New("database not initialized")
	}
	user := database.User{}
	result := sm.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return result.Error
	}
	result = sm.db.Model(&user).Update("default_application_status_id", statusID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
