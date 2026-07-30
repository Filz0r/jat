package services

import (
	"errors"
	"strconv"
	"time"

	"github.com/filz0r/jat/internal/database"
	"github.com/filz0r/jat/internal/utils"
	"github.com/google/uuid"
)

type applicationKindMapping struct {
	kind database.ApplicationStatusKind
	name string
}

func (sm *ServiceManager) CreateInitialApplicationStatus(userID uuid.UUID) error {
	initialStatuses := []applicationKindMapping{
		{
			kind: database.Irrelevant,
			name: "Didn't Apply",
		},
		{
			kind: database.Applied,
			name: "Applied",
		},
		{
			kind: database.Ghosted,
			name: "Ghosted",
		},
		{
			kind: database.Interviewed,
			name: "Contacted",
		},
		{
			kind: database.Interviewed,
			name: "1st Interview",
		},
		{
			kind: database.Interviewed,
			name: "2nd Interview",
		},
		{
			kind: database.Interviewed,
			name: "3rd Interview",
		},
		{
			kind: database.Accepted,
		},
	}
	if sm.db == nil {
		return errors.New("database not initialized")
	}
	data := make([]database.ApplicationStatus, 0, len(initialStatuses))
	for _, initialStatus := range initialStatuses {
		temp := database.ApplicationStatus{
			Status: initialStatus.name,
			UserID: userID,
			Kind:   initialStatus.kind,
		}
		data = append(data, temp)
	}
	result := sm.db.Create(&data)
	if result.Error != nil {
		return result.Error
	}
	// this is kinda garbage but it should work
	err := sm.SetUserDefaultApplicationStatus(userID, data[1].ID)
	return err
}

func (sm *ServiceManager) CreateApplicationStatus(applicationStatus database.ApplicationStatus) (database.ApplicationStatus, error) {
	if sm.db == nil {
		return database.ApplicationStatus{}, errors.New("database not initialized")
	}
	result := sm.db.Create(&applicationStatus)
	return applicationStatus, result.Error
}

func (sm *ServiceManager) UpdateApplicationStatus(applicationStatus database.ApplicationStatus) (database.ApplicationStatus, error) {
	if sm.db == nil {
		return database.ApplicationStatus{}, errors.New("database not initialized")
	}
	applicationStatus.UpdatedAt = time.Now()
	result := sm.db.Save(&applicationStatus)
	return applicationStatus, result.Error
}

func (sm *ServiceManager) GetApplicationStatus(id uint) (database.ApplicationStatus, error) {
	if sm.db == nil {
		return database.ApplicationStatus{}, errors.New("database not initialized")
	}
	var applicationStatus database.ApplicationStatus
	result := sm.db.
		Where("id = ?", id).
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

func (sm *ServiceManager) GetAllApplicationStatusAdmin() ([]database.ApplicationStatus, error) {
	if sm.db == nil {
		return []database.ApplicationStatus{}, errors.New("database not initialized")
	}
	var applicationStatus []database.ApplicationStatus
	data := sm.db.Find(&applicationStatus)
	if data.Error != nil {
		return []database.ApplicationStatus{}, data.Error
	}
	return applicationStatus, nil
}
