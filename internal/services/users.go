package services

import (
	"errors"
	"time"

	"github.com/filz0r/jat/internal/auth"
	"github.com/filz0r/jat/internal/database"
	"github.com/google/uuid"
)

func (sm *ServiceManager) CreateUser(
	username,
	email,
	password string,
) (database.User, error) {
	if sm.db == nil {
		return database.User{}, errors.New("database not initialized")
	}

	user := database.User{
		ID:        uuid.New(),
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	hashedPW, err := auth.HashPassword(user.Password)

	if err != nil {
		return database.User{}, err
	}

	user.Password = hashedPW
	user.UserSettings = database.UserSettings{
		UserID:    user.ID,
		SetupStep: 1,
		IsEnabled: true,
	}

	result := sm.db.Create(&user)
	if result.Error != nil {
		return database.User{}, result.Error
	}
	return user, nil
}

func (sm *ServiceManager) UpdateUser(
	userID uuid.UUID,
	username,
	password,
	email string,
) (database.User, error) {
	if sm.db == nil {
		return database.User{}, errors.New("database not initialized")
	}

	user, err := sm.GetUserByID(userID)
	if err != nil {
		return database.User{}, err
	}

	if username != "" {
		user.Username = username
	}

	if password != "" {
		same, err := auth.CheckPasswordHash(password, user.Password)
		if err != nil {
			return database.User{}, err
		}
		if same {
			return database.User{}, errors.New("passwords are the same")
		}
		hashedPW, err := auth.HashPassword(password)
		if err != nil {
			return database.User{}, err
		}
		user.Password = hashedPW
	}
	if email != "" {
		user.Email = email
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
	return user.UserSettings.IsAdmin
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
	if sm.IsStatusArchived(statusID) {
		return errors.New("cannot set a status that is archived as default")
	}
	if !sm.DoesUserOwnApplicationStatus(id, statusID) {
		return errors.New("cannot set a status that is not owned as default")
	}
	result := sm.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return result.Error
	}
	result = sm.db.Model(&user.UserSettings).Update("default_application_status_id", statusID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (sm *ServiceManager) ChangeUserAdminStatus(userID uuid.UUID, give bool) error {
	if sm.db == nil {
		return errors.New("database not initialized")
	}
	user, err := sm.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user.UserSettings.IsAdmin == give {
		msg := "this user already "
		if !give {
			msg += "isn't an admin"
		} else {
			msg += "is an admin"
		}
		return errors.New(msg)
	}
	user.UserSettings.IsAdmin = give
	result := sm.db.Save(&user.UserSettings)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
