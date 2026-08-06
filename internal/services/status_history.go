package services

import (
	"errors"
	"time"

	"github.com/filz0r/jat/internal/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (sm *ServiceManager) CreateApplicationStatusChange(
	tx *gorm.DB,
	applicationID, newStatusID uint,
	oldStatusID *uint,
) error {
	if tx == nil {
		return errors.New("tx is nil")
	}
	entry := database.StatusHistory{
		ApplicationID: applicationID,
		NewStatusID:   newStatusID,
		OldStatusID:   oldStatusID,
		CreatedAt:     time.Now(),
	}
	result := tx.Create(&entry)
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return errors.New("status not found")
	}
	return nil
}

func (sm *ServiceManager) GetJobApplicationStatusChanges(
	userID uuid.UUID,
	jobID uint,
) ([]database.StatusHistory, error) {
	if sm.db == nil {
		return nil, errors.New("database is not initialized")
	}

	var data []database.StatusHistory
	result := sm.db.
		Preload("OldStatus").
		Preload("NewStatus").
		Find(&data, "application_id = ?", jobID)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, entry := range data {
		if !sm.DoesUserOwnJobApplication(entry.ApplicationID, userID) && !sm.IsUserAdmin(userID) {
			return nil, errors.New("user is not own job application")
		}
	}

	return data, nil

}
