package services

import (
	"errors"

	"github.com/filz0r/jat/internal/database"
	"github.com/google/uuid"
)

func (sm *ServiceManager) CreateJobNote(
	jobID,
	statusID uint,
	note string,
	userID uuid.UUID,
) (database.ApplicationNote, error) {
	if sm.db == nil {
		return database.ApplicationNote{}, errors.New("database not initialized")
	}
	if note == "" {
		return database.ApplicationNote{}, errors.New("a note must have a value")
	}

	ownsStatus := sm.DoesUserOwnApplicationStatus(userID, statusID)
	if !ownsStatus {
		return database.ApplicationNote{}, errors.New("user does not own job application status")
	}

	ownsJob := sm.DoesUserOwnJobApplication(jobID, userID)
	if !ownsJob {
		return database.ApplicationNote{}, errors.New("this user does not own this job application")
	}

	dbNote := database.ApplicationNote{
		ApplicationID: jobID,
		UserID:        userID,
		Body:          note,
		StatusID:      statusID,
	}

	query := sm.db.Create(&dbNote)
	if query.Error != nil {
		return database.ApplicationNote{}, query.Error
	}

	result := query.Preload("Status").First(&dbNote)
	if result.Error != nil {
		return database.ApplicationNote{}, result.Error
	}

	return dbNote, nil
}

func (sm *ServiceManager) GetJobApplicationNoteByID(
	noteID,
	jobID uint,
	userID uuid.UUID,
) (database.ApplicationNote, error) {
	if sm.db == nil {
		return database.ApplicationNote{}, errors.New("database not initialized")
	}
	var dbNote database.ApplicationNote
	query := sm.db.Where("id = ? AND application_id = ?", noteID, jobID)

	if !sm.IsUserAdmin(userID) {
		query = query.Where("user_id = ?", userID)
	}

	result := query.First(&dbNote)
	if result.Error != nil {
		return database.ApplicationNote{}, result.Error
	}

	return dbNote, nil
}

func (sm *ServiceManager) GetJobApplicationNotesByID(
	jobID uint,
	userID uuid.UUID,
) ([]database.ApplicationNote, error) {
	if sm.db == nil {
		return []database.ApplicationNote{}, errors.New("database not initialized")
	}
	var dbNotes []database.ApplicationNote
	isAdmin := sm.IsUserAdmin(userID)
	if !isAdmin && !sm.DoesUserOwnJobApplication(jobID, userID) {
		return []database.ApplicationNote{}, errors.New("this user does not own this job application")
	}

	result := sm.db.
		Where("user_id = ? and application_id = ?", userID, jobID).
		Find(&dbNotes)

	if result.Error != nil {
		return []database.ApplicationNote{}, result.Error
	}

	return dbNotes, nil
}

func (sm *ServiceManager) UpdateJobNoteByID(
	noteID,
	jobID uint,
	userID uuid.UUID,
	updatedNote string,
) (database.ApplicationNote, error) {
	if sm.db == nil {
		return database.ApplicationNote{}, errors.New("database not initialized")
	}
	isAdmin := sm.IsUserAdmin(userID)
	if !isAdmin && !sm.DoesUserOwnJobApplication(jobID, userID) {
		return database.ApplicationNote{}, errors.New("this user does not own this job application")
	}
	ownsJob := sm.DoesUserOwnJobApplication(jobID, userID)
	if !ownsJob && !isAdmin {
		return database.ApplicationNote{}, errors.New("user does not own this job application")
	}
	noteData, err := sm.GetJobApplicationNoteByID(noteID, jobID, userID)
	if err != nil {
		return database.ApplicationNote{}, err
	}
	noteData.Body = updatedNote
	result := sm.db.Updates(&noteData)
	if result.Error != nil {
		return database.ApplicationNote{}, result.Error
	}
	return noteData, nil
}

func (sm *ServiceManager) DeleteApplicationNoteByID(
	noteID uint,
	jobID uint,
	userID uuid.UUID,
) error {
	if sm.db == nil {
		return errors.New("database not initialized")
	}
	dbNote, err := sm.GetJobApplicationNoteByID(noteID, jobID, userID)
	if err != nil {
		return err
	}
	result := sm.db.Delete(&dbNote)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
