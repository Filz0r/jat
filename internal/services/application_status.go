package services

import (
	"errors"
	"strconv"

	"github.com/filz0r/jat/internal/database"
	"github.com/filz0r/jat/internal/utils"
	"github.com/google/uuid"
)

func (sm *ServiceManager) CreateInitialApplicationStatus(userID uuid.UUID) error {
	initialStatuses := []string{
		"Didn't Apply",
		"Applied",
		"Ghosted",
		"Contacted",
		"1st Interview",
		"2nd Interview",
		"3rd Interview",
		"Rejected",
		"Accepted",
	}
	if sm.db == nil {
		return errors.New("database not initialized")
	}
	data := make([]database.ApplicationStatus, 0, len(initialStatuses))
	for _, initialStatus := range initialStatuses {
		temp := database.ApplicationStatus{
			Status: initialStatus,
			UserID: userID,
		}
		data = append(data, temp)
	}
	result := sm.db.Create(&data)
	return result.Error
}

func (sm *ServiceManager) CreateApplicationStatus(
	userID uuid.UUID,
	applicationStatus database.ApplicationStatus,
) (database.ApplicationStatus, error) {
	if sm.db == nil {
		return database.ApplicationStatus{}, errors.New("database not initialized")
	}
	applicationStatus.UserID = userID
	result := sm.db.Create(&applicationStatus)
	return applicationStatus, result.Error
}

func (sm *ServiceManager) UpdateApplicationStatus(
	userID uuid.UUID,
	applicationStatus database.ApplicationStatus,
) (database.ApplicationStatus, error) {
	if sm.db == nil {
		return database.ApplicationStatus{}, errors.New("database not initialized")
	}

	if applicationStatus.UserID != userID {
		return database.ApplicationStatus{}, errors.New("your user doesn't own this data")
	}
	applicationStatus.UserID = userID
	result := sm.db.Save(&applicationStatus)
	return applicationStatus, result.Error
}

func (sm *ServiceManager) GetApplicationStatus(
	userID uuid.UUID,
	name string,
) (database.ApplicationStatus, error) {
	if sm.db == nil {
		return database.ApplicationStatus{}, errors.New("database not initialized")
	}
	var applicationStatus database.ApplicationStatus
	result := sm.db.
		Where("status = ? and user_id = ?", name, userID).
		First(&applicationStatus)
	return applicationStatus, result.Error
}

func (sm *ServiceManager) FindApplicationStatusByName(
	userID uuid.UUID,
	name string,
) (database.ApplicationStatus, error) {
	if sm.db == nil {
		return database.ApplicationStatus{}, errors.New("database not initialized")
	}
	var applicationStatus database.ApplicationStatus
	result := sm.db.
		Where("lower(status) = lower(?) and user_id = ?", name, userID).
		First(&applicationStatus)
	return applicationStatus, result.Error
}

func (sm *ServiceManager) GetAllApplicationStatus(
	userID uuid.UUID,
) ([]database.ApplicationStatus, error) {
	if sm.db == nil {
		return []database.ApplicationStatus{}, errors.New("database not initialized")
	}
	var applicationStatus []database.ApplicationStatus
	result := sm.db.
		Where("user_id = ?", userID).
		Find(&applicationStatus)
	return applicationStatus, result.Error
}

func (sm *ServiceManager) DeleteApplicationStatus(
	userID uuid.UUID,
	name string,
) error {
	if sm.db == nil {
		return errors.New("database not initialized")
	}
	var applicationStatus database.ApplicationStatus
	result := sm.db.
		Delete(&applicationStatus, "status = ? and user_id = ?", name, userID)
	return result.Error
}

func (sm *ServiceManager) GetSuggestionForApplicationStatus(
	userID uuid.UUID,
) ([]utils.SuggestionRecord, error) {
	if sm.db == nil {
		return []utils.SuggestionRecord{}, errors.New("database not initialized")
	}
	var data []database.ApplicationStatus
	result := sm.db.Where("user_id = ?", userID).Find(&data)
	if result.Error != nil {
		return []utils.SuggestionRecord{}, result.Error
	}
	suggestionRecord := make([]utils.SuggestionRecord, 0, len(data))
	for _, d := range data {
		stringID := strconv.Itoa(int(d.ID))
		temp := utils.SuggestionRecord{

			Label: d.Status,
			Value: stringID,
			// Might be this still not sure so keep it here for a reminder
			//Value: utils.ToSnakeCase(d.Status),
		}
		suggestionRecord = append(suggestionRecord, temp)
	}
	return suggestionRecord, result.Error
}
