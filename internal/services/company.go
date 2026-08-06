package services

import (
	"errors"
	"strconv"
	"time"

	"github.com/filz0r/jat/internal/database"
	"github.com/filz0r/jat/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (sm *ServiceManager) FindCompanyByName(name string) (database.Company, error) {
	var company database.Company
	result := sm.db.Where("lower(name) = lower(?)", name).First(&company)
	if result.RowsAffected == 0 {
		return database.Company{}, errors.New("company not found")
	} else if result.Error != nil {
		return database.Company{}, result.Error
	}
	return company, nil
}

func (sm *ServiceManager) CreateCompany(userID uuid.UUID, companyName, websiteUrl string) (database.Company, error) {
	if sm.db == nil {
		return database.Company{}, errors.New("database not initialized")
	}

	var urlPtr *string
	if websiteUrl != "" {
		urlPtr = &websiteUrl
	} else {
		urlPtr = nil
	}

	company := database.Company{
		Name:      companyName,
		Website:   urlPtr,
		CreatedBy: userID,
		EditedBy:  userID,
	}

	result := sm.db.Create(&company)
	return company, result.Error
}

func (sm *ServiceManager) UpdateCompany(companyID uint, userID uuid.UUID, name, url string) (database.Company, error) {
	if sm.db == nil {
		return database.Company{}, errors.New("database not initialized")
	}
	var updated database.Company
	err := sm.db.Transaction(func(tx *gorm.DB) error {
		var current database.Company
		if err := tx.First(&current, companyID).Error; err != nil {
			return err
		}
		var urlPtr *string
		if url != "" {
			urlPtr = &url
		} else {
			urlPtr = nil
		}

		if current.Name == name && current.Website == urlPtr {
			return errors.New("nothing changed")
		}

		if err := recordCompanyChange(
			tx,
			companyID,
			name, current.Name,
			urlPtr, current.Website,
			userID,
		); err != nil {
			return err
		}

		current.EditedBy = userID
		current.Name = name
		current.Website = urlPtr
		current.UpdatedAt = time.Now()

		if err := tx.Save(&current).Error; err != nil {
			return err
		}

		updated = current
		return nil
	})

	if err != nil {
		return database.Company{}, err
	}

	return updated, nil
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

func (sm *ServiceManager) GetAllCompanies() ([]database.Company, error) {
	if sm.db == nil {
		return []database.Company{}, errors.New("database not initialized")
	}
	var companies []database.Company
	result := sm.db.Find(&companies)
	return companies, result.Error
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

func (sm *ServiceManager) GetCompanyByID(id uint) (database.Company, error) {
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
	company, err := sm.GetCompanyByID(id)

	if err != nil {
		return err
	}
	if !sm.IsUserAdmin(userID) {
		return errors.New("only admin users can delete companies")
	}

	company.EditedBy = userID

	result := sm.db.Delete(&company)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (sm *ServiceManager) DoesCompanyExist(companyID uint) bool {
	if sm.db == nil {
		return false
	}
	var company database.Company
	result := sm.db.First(&company, companyID)
	if result.Error != nil {
		return false
	}
	return true
}
