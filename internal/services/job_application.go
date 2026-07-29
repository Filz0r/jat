package services

import (
	"errors"

	"github.com/filz0r/jat/internal/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TODO: needs cleanup from AI SLOP

// DefaultApplicationStatus is the status every new job application starts
// with. It is part of the seeded set created per user, so the lookup is
// scoped to the owner. The UI uses it to pre-select the status picker too.
const DefaultApplicationStatus = "Applied"

// resolveCompany returns the company matching name (case-insensitive),
// creating it when none exists. Companies are global records shared by every
// user, so a mismatch in casing must never produce a duplicate row.
func (sm *ServiceManager) resolveCompany(name string) (database.Company, error) {
	company, err := sm.FindCompanyByName(name)
	if err == nil {
		return company, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return database.Company{}, err
	}
	company = database.Company{Name: name}
	if err := sm.db.Create(&company).Error; err != nil {
		return database.Company{}, err
	}
	return company, nil
}

// ownsApplication guards every per-user operation on a job application: the
// record must exist and belong to the caller.
func (sm *ServiceManager) ownsApplication(userID uuid.UUID, applicationID uint) error {
	var count int64
	result := sm.db.Model(&database.JobApplication{}).
		Where("id = ? AND user_id = ?", applicationID, userID).
		Count(&count)
	if result.Error != nil {
		return result.Error
	}
	if count == 0 {
		return errors.New("job application not found")
	}
	return nil
}

// CreateJobApplication creates the application, links its company (resolving
// or creating it) and writes the initial status-history entry, all in one
// transaction. When statusID is nil the user's default status applies.
func (sm *ServiceManager) CreateJobApplication(
	userID uuid.UUID,
	title, url, companyName string,
	statusID *uint,
) (database.JobApplication, error) {
	if sm.db == nil {
		return database.JobApplication{}, errors.New("database not initialized")
	}

	resolvedStatusID := statusID
	if resolvedStatusID == nil {
		status, err := sm.FindApplicationStatusByName(userID, DefaultApplicationStatus)
		if err != nil {
			return database.JobApplication{}, err
		}
		resolvedStatusID = &status.ID
	}

	var app database.JobApplication
	err := sm.db.Transaction(func(tx *gorm.DB) error {
		company, err := sm.resolveCompany(companyName)
		if err != nil {
			return err
		}
		app = database.JobApplication{
			Title:    title,
			Url:      url,
			UserID:   userID,
			StatusID: *resolvedStatusID,
		}
		if err := tx.Create(&app).Error; err != nil {
			return err
		}
		if err := tx.Model(&app).Association("Companies").Append(&company); err != nil {
			return err
		}
		return recordStatusChange(tx, app.ID, *resolvedStatusID)
	})
	if err != nil {
		return database.JobApplication{}, err
	}
	return sm.GetJobApplication(userID, app.ID)
}

// GetJobApplication loads one application with its status and companies,
// scoped to the owner.
func (sm *ServiceManager) GetJobApplication(
	userID uuid.UUID,
	applicationID uint,
) (database.JobApplication, error) {
	if sm.db == nil {
		return database.JobApplication{}, errors.New("database not initialized")
	}
	if err := sm.ownsApplication(userID, applicationID); err != nil {
		return database.JobApplication{}, err
	}
	var app database.JobApplication
	result := sm.db.
		Preload("Status").
		Preload("Companies").
		First(&app, applicationID)
	return app, result.Error
}

// GetJobApplications lists every application owned by userID, newest first.
func (sm *ServiceManager) GetJobApplications(
	userID uuid.UUID,
) ([]database.JobApplication, error) {
	if sm.db == nil {
		return []database.JobApplication{}, errors.New("database not initialized")
	}
	var applications []database.JobApplication
	result := sm.db.
		Where("user_id = ?", userID).
		Preload("Status").
		Preload("Companies").
		Order("created_at DESC").
		Find(&applications)
	return applications, result.Error
}

// UpdateJobApplication saves the editable fields, re-links the company and,
// when the status changed, appends a status-history entry — all atomically.
func (sm *ServiceManager) UpdateJobApplication(
	userID uuid.UUID,
	applicationID uint,
	title, url, companyName string,
	statusID uint,
) (database.JobApplication, error) {
	if sm.db == nil {
		return database.JobApplication{}, errors.New("database not initialized")
	}

	var app database.JobApplication
	err := sm.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("id = ? AND user_id = ?", applicationID, userID).First(&app)
		if result.Error != nil {
			return result.Error
		}

		app.Title = title
		app.Url = url
		statusChanged := app.StatusID != statusID
		app.StatusID = statusID
		// Save a copy with the associations stripped so GORM does not upsert
		// them; the company link is changed explicitly further down.
		record := app
		record.Status = database.ApplicationStatus{}
		record.Companies = nil
		if err := tx.Save(&record).Error; err != nil {
			return err
		}

		company, err := sm.resolveCompany(companyName)
		if err != nil {
			return err
		}
		// An application carries exactly one company in this UI; Clear +
		// Append is deterministic (unlike Replace, which diffs against the
		// in-memory association field that was never preloaded here).
		if err := tx.Model(&app).Association("Companies").Clear(); err != nil {
			return err
		}
		if err := tx.Model(&app).Association("Companies").Append(&company); err != nil {
			return err
		}
		if statusChanged {
			if err := recordStatusChange(tx, app.ID, statusID); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return database.JobApplication{}, err
	}
	return sm.GetJobApplication(userID, app.ID)
}

// DeleteJobApplication removes an application owned by userID.
func (sm *ServiceManager) DeleteJobApplication(userID uuid.UUID, applicationID uint) error {
	if sm.db == nil {
		return errors.New("database not initialized")
	}
	if err := sm.ownsApplication(userID, applicationID); err != nil {
		return err
	}
	result := sm.db.Delete(&database.JobApplication{}, applicationID)
	return result.Error
}
