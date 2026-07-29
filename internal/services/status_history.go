package services

import (
	"errors"

	"github.com/filz0r/jat/internal/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TODO: needs cleanup from AI SLOP

// recordStatusChange appends one entry to an application's history. It is
// always called inside the transaction of the operation that changed the
// status, so a failed history write rolls the change back too.
func recordStatusChange(tx *gorm.DB, applicationID, statusID uint) error {
	entry := database.StatusHistory{
		ApplicationID: applicationID,
		StatusID:      statusID,
	}
	return tx.Create(&entry).Error
}

// GetStatusHistory returns an application's status changes, oldest first,
// scoped to the owning user.
func (sm *ServiceManager) GetStatusHistory(
	userID uuid.UUID,
	applicationID uint,
) ([]database.StatusHistory, error) {
	if sm.db == nil {
		return []database.StatusHistory{}, errors.New("database not initialized")
	}
	if err := sm.ownsApplication(userID, applicationID); err != nil {
		return []database.StatusHistory{}, err
	}
	var history []database.StatusHistory
	result := sm.db.
		Where("application_id = ?", applicationID).
		Preload("Status").
		Order("created_at ASC").
		Find(&history)
	return history, result.Error
}
