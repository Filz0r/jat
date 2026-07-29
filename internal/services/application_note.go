package services

import (
	"errors"

	"github.com/filz0r/jat/internal/database"
	"github.com/google/uuid"
)

func (sm *ServiceManager) AddNote(
	userID uuid.UUID,
	applicationID uint,
	body string,
	statusID *uint,
) (database.ApplicationNote, error) {
	if sm.db == nil {
		return database.ApplicationNote{}, errors.New("database not initialized")
	}
	if err := sm.ownsApplication(userID, applicationID); err != nil {
		return database.ApplicationNote{}, err
	}
	note := database.ApplicationNote{
		Body:          body,
		ApplicationID: applicationID,
		StatusID:      statusID,
	}
	result := sm.db.Create(&note)
	return note, result.Error
}

func (sm *ServiceManager) GetNotes(
	userID uuid.UUID,
	applicationID uint,
) ([]database.ApplicationNote, error) {
	if sm.db == nil {
		return []database.ApplicationNote{}, errors.New("database not initialized")
	}
	if err := sm.ownsApplication(userID, applicationID); err != nil {
		return []database.ApplicationNote{}, err
	}
	var notes []database.ApplicationNote
	result := sm.db.
		Where("application_id = ?", applicationID).
		Order("created_at ASC").
		Find(&notes)
	return notes, result.Error
}

func (sm *ServiceManager) loadOwnedNote(
	userID uuid.UUID,
	noteID uint,
) (database.ApplicationNote, error) {
	var note database.ApplicationNote
	result := sm.db.
		Joins("JOIN job_applications ON job_applications.id = application_notes.application_id").
		Where("application_notes.id = ? AND job_applications.user_id = ?", noteID, userID).
		First(&note)
	return note, result.Error
}

func (sm *ServiceManager) UpdateNote(
	userID uuid.UUID,
	noteID uint,
	body string,
) (database.ApplicationNote, error) {
	if sm.db == nil {
		return database.ApplicationNote{}, errors.New("database not initialized")
	}
	note, err := sm.loadOwnedNote(userID, noteID)
	if err != nil {
		return database.ApplicationNote{}, err
	}
	note.Body = body
	result := sm.db.Save(&note)
	return note, result.Error
}

func (sm *ServiceManager) DeleteNote(userID uuid.UUID, noteID uint) error {
	if sm.db == nil {
		return errors.New("database not initialized")
	}
	note, err := sm.loadOwnedNote(userID, noteID)
	if err != nil {
		return err
	}
	result := sm.db.Delete(&note)
	return result.Error
}
