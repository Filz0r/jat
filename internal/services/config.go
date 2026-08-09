package services

import (
	"errors"
	"time"

	"github.com/filz0r/jat/internal/database"
	"github.com/google/uuid"
)

func (sm *ServiceManager) getConfig(key string) (database.Config, error) {
	if sm.db == nil {
		return database.Config{}, errors.New("database not initialized")
	}
	data := database.Config{
		Key: key,
	}
	res := sm.db.Where(&data).First(&data)
	if res.Error != nil {
		return database.Config{}, res.Error
	}
	return data, nil
}

func (sm *ServiceManager) InitialConfigsExist() (bool, error) {
	if sm.db == nil {
		return false, errors.New("database not initialized")
	}
	var count int64
	res := sm.db.Model(&database.Config{}).
		Where("key IN ?", []string{"isInitialized", "applicationMode", "firstAdmin"}).
		Count(&count)
	if res.Error != nil {
		return false, res.Error
	}
	return count == 3, nil
}

func (sm *ServiceManager) CreateInitialConfigs() error {
	if sm.db == nil {
		return errors.New("database not initialized")
	}
	configs := []database.Config{
		{
			Key:   "isInitialized",
			Value: "false",
			Type:  database.BooleanConfig,
		},
		{
			Key:   "applicationMode",
			Value: database.ServerMode.String(),
			Type:  database.JatModeConfig,
		},
		{
			Key:   "firstAdmin",
			Value: uuid.Nil.String(),
			Type:  database.IDConfig,
		},
	}

	for _, config := range configs {
		config.CreatedAt = time.Now()
		config.UpdatedAt = time.Now()
		res := sm.db.Create(&config)
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}

func (sm *ServiceManager) IsInitialized() bool {
	if sm.db == nil {
		return false
	}
	data := database.Config{}
	res := sm.db.First(&data, "key = ?", "isInitialized")
	if res.Error != nil {
		return false
	}
	boolVal, err := data.Type.ToType(data.Value)
	if err != nil {
		return false
	}
	return boolVal.(bool)
}

func (sm *ServiceManager) SetInitialized() bool {
	if sm.db == nil {
		return false
	}
	if sm.IsInitialized() {
		return true
	}
	data := database.Config{}
	res := sm.db.First(&data, "key = ?", "isInitialized")
	if res.Error != nil {
		return false
	}
	data.UpdatedAt = time.Now()
	data.Value = "true"
	res = sm.db.Save(&data)
	if res.Error != nil {
		return false
	}
	return true
}

func (sm *ServiceManager) GetFirstAdmin() uuid.UUID {
	data, err := sm.getConfig("firstAdmin")
	if err != nil {
		return uuid.Nil
	}

	id, err := data.Type.ToType(data.Value)
	if err != nil {
		return uuid.Nil
	}
	res := id.(uuid.UUID)
	return res
}

func (sm *ServiceManager) SetFirstAdmin(userID uuid.UUID) bool {
	if sm.db == nil {
		return false
	}
	if sm.IsInitialized() {
		return false
	}
	data, err := sm.getConfig("firstAdmin")
	if err != nil {
		return false
	}
	_, err = sm.GetUserByID(userID)
	if err != nil {
		return false
	}
	data = database.Config{
		Key:       data.Key,
		Value:     userID.String(),
		Type:      data.Type,
		CreatedAt: data.CreatedAt,
		UpdatedAt: time.Now(),
	}
	res := sm.db.Save(&data)
	if res.Error != nil {
		return false
	}
	return true
}

func (sm *ServiceManager) IsFirstAdmin(userID uuid.UUID) bool {
	if sm.db == nil {
		return false
	}
	data := sm.GetFirstAdmin()
	if data != uuid.Nil {
		return false
	}
	return data == userID
}

func (sm *ServiceManager) SetServerMode() (bool, error) {
	if sm.db == nil {
		return false, errors.New("database not initialized")
	}
	exists, err := sm.GetServerMode()
	if err != nil || exists.Valid() {
		return false, err
	}
	config := database.Config{
		Key:       "applicationMode",
		Value:     database.ServerMode.String(),
		Type:      database.JatModeConfig,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	res := sm.db.Save(&config)
	if res.Error != nil {
		return false, res.Error
	}
	return true, nil
}

func (sm *ServiceManager) SetStandaloneMode() (bool, error) {
	if sm.db == nil {
		return false, errors.New("database not initialized")
	}
	exists, err := sm.GetServerMode()
	if err != nil || exists.Valid() {
		return false, err
	}
	config := database.Config{
		Key:       "applicationMode",
		Value:     database.StandaloneMode.String(),
		Type:      database.JatModeConfig,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	res := sm.db.Save(&config)
	if res.Error != nil {
		return false, res.Error
	}
	return true, nil
}

func (sm *ServiceManager) GetServerMode() (database.JatMode, error) {
	if sm.db == nil {
		return "", errors.New("database not initialized")
	}
	dbData, err := sm.getConfig("applicationMode")
	if err != nil {
		return "", err
	}
	converted, err := dbData.Type.ToType(dbData.Value)
	if err != nil {
		return "", err
	}
	return converted.(database.JatMode), nil
}
