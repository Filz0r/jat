package services

import (
	"errors"
	"strconv"

	"github.com/filz0r/jat/internal/database"
	"github.com/filz0r/jat/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (sm *ServiceManager) FindCompanyByName(name string) (database.Company, error) {
	var company database.Company
	result := sm.db.Where("lower(name) = lower(?)", name).First(&company)
	return company, result.Error
}

func (sm *ServiceManager) CreateCompany(userID uuid.UUID, company database.Company) (database.Company, error) {
	if sm.db == nil {
		return database.Company{}, errors.New("database not initialized")
	}
	company.CreatedBy = userID
	company.EditedBy = userID

	result := sm.db.Create(&company)
	return company, result.Error
}

func (sm *ServiceManager) UpdateCompany(userID uuid.UUID, company database.Company) (database.Company, error) {
	if sm.db == nil {
		return database.Company{}, errors.New("database not initialized")
	}
	var updated database.Company
	err := sm.db.Transaction(func(tx *gorm.DB) error {
		var current database.Company
		if err := tx.First(&current, company.ID).Error; err != nil {
			return err
		}

		company.CreatedBy = current.CreatedBy
		company.EditedBy = userID

		if err := recordCompanyChange(tx, current, company, userID); err != nil {
			return err
		}

		if err := tx.Save(&company).Error; err != nil {
			return err
		}

		updated = company
		return nil
	})

	return updated, err
}

func (sm *ServiceManager) GetCompany(name string) (database.Company, error) {
	if sm.db == nil {
		return database.Company{}, errors.New("database not initialized")
	}
	var company database.Company
	result := sm.db.Where("name = ?", name).First(&company)
	return company, result.Error
}

func (sm *ServiceManager) FindCompany(name string) (database.Company, error) {
	if sm.db == nil {
		return database.Company{}, errors.New("database not initialized")
	}
	return sm.FindCompanyByName(name)
}

func (sm *ServiceManager) GetApplicationsForCompany(
	userID uuid.UUID,
	companyID uint,
) ([]database.JobApplication, error) {
	if sm.db == nil {
		return []database.JobApplication{}, errors.New("database not initialized")
	}
	var applications []database.JobApplication
	result := sm.db.
		Joins("JOIN company_application ON company_application.application_id = job_applications.id").
		Where("company_application.company_id = ? AND job_applications.user_id = ?", companyID, userID).
		Preload("Status").
		Preload("Companies").
		Order("job_applications.created_at DESC").
		Find(&applications)
	return applications, result.Error
}

func (sm *ServiceManager) GetAllCompanies() ([]database.Company, error) {
	if sm.db == nil {
		return []database.Company{}, errors.New("database not initialized")
	}
	var companies []database.Company
	result := sm.db.Find(&companies)
	return companies, result.Error
}

func (sm *ServiceManager) DeleteCompany(name string, user database.User) error {
	if sm.db == nil {
		return errors.New("database not initialized")
	}
	if !user.IsAdmin {
		return errors.New("only admin users can delete companies")
	}
	var company database.Company
	result := sm.db.Delete(&company, "name = ?", name)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("company not found")
	}
	return nil
}

func (sm *ServiceManager) GetAllCompaniesSuggestions() ([]utils.SuggestionRecord, error) {
	if sm.db == nil {
		return []utils.SuggestionRecord{}, errors.New("database not initialized")
	}
	result, err := sm.GetAllCompanies()
	if err != nil {
		return []utils.SuggestionRecord{}, err
	}
	suggestions := make([]utils.SuggestionRecord, 0, len(result))
	for _, company := range result {
		stringID := strconv.Itoa(int(company.ID))
		temp := utils.SuggestionRecord{
			Label: company.Name,
			Value: stringID,
		}
		suggestions = append(suggestions, temp)

	}
	return suggestions, nil
}

func (sm *ServiceManager) GetCompanyByID(id int) (database.Company, error) {
	if sm.db == nil {
		return database.Company{}, errors.New("database not initialized")
	}
	var company database.Company
	result := sm.db.First(&company, id)
	return company, result.Error
}

func (sm *ServiceManager) DeleteCompanyByID(id uint, userID uuid.UUID) error {
	if sm.db == nil {
		return errors.New("database not initialized")
	}
	var company database.Company
	result := sm.db.Delete(&company, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("company not found")
	}
	company.EditedBy = userID
	result = sm.db.Save(&company)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
