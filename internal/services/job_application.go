package services

import (
	"errors"
	"time"

	"github.com/filz0r/jat/internal/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (sm *ServiceManager) GetUserJobApplications(userID uuid.UUID) ([]database.JobApplication, error) {
	if sm.db == nil {
		return nil, errors.New("database not initialized")
	}

	var applications []database.JobApplication

	result := sm.db.
		Where("user_id = ?", userID).
		Preload("Status").
		Preload("Company").
		Order("created_at desc").
		Find(&applications)

	if result.Error != nil {
		return nil, result.Error
	}

	return applications, nil
}

func (sm *ServiceManager) CreateUserJobApplication(
	userID uuid.UUID,
	statusID, companyID uint,
	title, url string,
	createdAt time.Time,
) (database.JobApplication, error) {
	if sm.db == nil {
		return database.JobApplication{}, errors.New("database not initialized")
	}

	if !sm.DoesUserOwnApplicationStatus(userID, statusID) {
		return database.JobApplication{}, errors.New("user does not own job application status")
	}

	if !sm.DoesCompanyExist(companyID) {
		return database.JobApplication{}, errors.New("company does not exist")
	}
	var application database.JobApplication
	err := sm.db.Transaction(func(tx *gorm.DB) error {
		application = database.JobApplication{
			CompanyID: companyID,
			StatusID:  statusID,
			Title:     title,
			Url:       url,
			UserID:    userID,
		}
		application.CreatedAt = createdAt
		if err := tx.Create(&application).Error; err != nil {
			return err
		}

		if err := sm.CreateApplicationStatusChange(
			tx,
			application.ID,
			statusID,
			nil,
		); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return database.JobApplication{}, err
	}

	res := sm.db.
		Where("id = ? and user_id = ?", application.ID, application.UserID).
		Preload("Status").
		Preload("Company").
		First(&application)
	if res.Error != nil {
		return database.JobApplication{}, res.Error
	}
	return application, nil
}

func (sm *ServiceManager) GetJobApplicationById(id uint) (database.JobApplication, error) {
	if sm.db == nil {
		return database.JobApplication{}, errors.New("database not initialized")
	}
	application := database.JobApplication{}
	result := sm.db.Where("id = ?", id).
		Preload("Status").
		Preload("Company").
		First(&application)
	if result.Error != nil {
		return database.JobApplication{}, result.Error
	}
	return application, nil
}

func (sm *ServiceManager) UpdateJobApplicationStatus(
	userID uuid.UUID,
	applicationID,
	newStatusID uint,
) (database.JobApplication, error) {
	if sm.db == nil {
		return database.JobApplication{}, errors.New("database not initialized")
	}
	if !sm.DoesUserOwnJobApplication(applicationID, userID) {
		return database.JobApplication{}, errors.New("user does not own job application")
	}
	err := sm.db.Transaction(func(tx *gorm.DB) error {
		var application database.JobApplication
		if err := tx.Where("id = ?", applicationID).First(&application).Error; err != nil {
			return err
		}
		oldStatus := application.StatusID
		application.StatusID = newStatusID
		if err := tx.Save(&application).Error; err != nil {
			return err
		}
		return sm.CreateApplicationStatusChange(tx, application.ID, newStatusID, &oldStatus)
	})
	if err != nil {
		return database.JobApplication{}, err
	}
	return sm.GetJobApplicationById(applicationID)
}

func (sm *ServiceManager) DeleteJobApplication(applicationID uint, userID uuid.UUID) error {
	if sm.db == nil {
		return errors.New("database not initialized")
	}
	if !sm.DoesUserOwnJobApplication(applicationID, userID) && !sm.IsUserAdmin(userID) {
		return errors.New("user does not own job application")
	}
	res := sm.db.
		Delete(&database.JobApplication{},
			"id = ?",
			applicationID,
		)
	if res.Error != nil {
		return res.Error
	} else if res.RowsAffected == 0 {
		return errors.New("job application does not exist")
	}
	return nil
}

func (sm *ServiceManager) DoesUserOwnJobApplication(jobID uint, userID uuid.UUID) bool {
	if sm.db == nil {
		return false
	}
	jobApplication, err := sm.GetJobApplicationById(jobID)
	if err != nil {
		return false
	}
	return jobApplication.UserID == userID
}
