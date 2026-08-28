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

func (sm *ServiceManager) UpdateCompany(
	tx *gorm.DB,
	companyID uint,
	userID uuid.UUID,
	name, url string,
) (database.Company, error) {
	if sm.db == nil {
		return database.Company{}, errors.New("database not initialized")
	}
	db := sm.transactionOrDefault(tx)
	var updated database.Company
	err := db.Transaction(func(_tx *gorm.DB) error {
		var current database.Company
		if err := _tx.First(&current, companyID).Error; err != nil {
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
			_tx,
			companyID,
			name, current.Name,
			current.Website, urlPtr,
			userID,
		); err != nil {
			return err
		}

		current.EditedBy = userID
		current.Name = name
		current.Website = urlPtr
		current.UpdatedAt = time.Now()

		if err := _tx.Save(&current).Error; err != nil {
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

func (sm *ServiceManager) GetAllCompanies(
	tx *gorm.DB,
	userID uuid.UUID,
	includeTotal, includeUser bool,
) ([]database.Company, map[uint]utils.CompanyCounts, error) {
	if sm.db == nil {
		return []database.Company{}, map[uint]utils.CompanyCounts{}, errors.New("database not initialized")
	}

	db := sm.transactionOrDefault(tx)

	var companies []database.Company
	var counts map[uint]utils.CompanyCounts

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := db.Find(&companies).Error; err != nil {
			return err
		}
		if len(companies) == 0 || (!includeTotal && !includeUser) {
			return nil
		}
		counts = make(map[uint]utils.CompanyCounts, len(companies))
		for _, company := range companies {
			if includeUser {
				userCount, err := sm.CountCompanyApplications(tx, company.ID, userID, true)
				if err != nil {
					return err
				}
				c, ok := counts[company.ID]
				if !ok {
					c = utils.CompanyCounts{}
				}
				c.UserCount = userCount
				counts[company.ID] = c
			}
			if includeTotal {
				totalCount, err := sm.CountCompanyApplications(tx, company.ID, userID, false)
				if err != nil {
					return err
				}
				c, ok := counts[company.ID]
				if !ok {
					c = utils.CompanyCounts{}
				}
				c.TotalCount = totalCount
				counts[company.ID] = c
			}
		}
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return companies, counts, nil
}

func (sm *ServiceManager) GetAllCompaniesSuggestions() ([]utils.SuggestionRecord, error) {
	if sm.db == nil {
		return []utils.SuggestionRecord{}, errors.New("database not initialized")
	}
	result, _, err := sm.GetAllCompanies(nil, uuid.Nil, false, false)
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

func (sm *ServiceManager) GetCompanyByID(
	tx *gorm.DB,
	id uint,
	userID uuid.UUID,
	includeTotal,
	includeUser bool,
) (database.Company, utils.CompanyCounts, error) {
	if sm.db == nil {
		return database.Company{}, utils.CompanyCounts{}, errors.New("database not initialized")
	}
	db := sm.transactionOrDefault(tx)
	var company database.Company
	var counts utils.CompanyCounts

	err := db.Transaction(func(tx *gorm.DB) error {
		result := db.First(&company, id)
		if result.Error != nil {
			return result.Error
		}
		if includeTotal {
			totalCount, err := sm.CountCompanyApplications(tx, company.ID, userID, false)
			if err != nil {
				return err
			}
			counts.TotalCount = totalCount
		}
		if includeUser {
			userCount, err := sm.CountCompanyApplications(tx, company.ID, userID, true)
			if err != nil {
				return err
			}
			counts.UserCount += userCount
		}
		return nil
	})

	if err != nil {
		return database.Company{}, utils.CompanyCounts{}, err
	}

	return company, counts, nil
}

func (sm *ServiceManager) DeleteCompanyByID(tx *gorm.DB, id uint, userID uuid.UUID) error {
	if sm.db == nil {
		return errors.New("database not initialized")
	}
	db := sm.transactionOrDefault(tx)

	err := db.Transaction(func(tx *gorm.DB) error {
		company, _, err := sm.GetCompanyByID(tx, id, userID, false, false)

		if err != nil {
			return err
		}

		if !sm.IsUserAdmin(userID) {
			return errors.New("only admin users can delete companies")
		}

		company.EditedBy = userID

		result := tx.Delete(&company)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return err
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

func (sm *ServiceManager) CountCompanyApplications(tx *gorm.DB, companyID uint, userID uuid.UUID, includeUser bool) (int64, error) {
	if sm.db == nil {
		return 0, errors.New("database not initialized")
	}
	db := sm.transactionOrDefault(tx)
	var count int64
	query := db.Model(&database.JobApplication{}).Where("company_id = ?", companyID)
	if includeUser {
		query = query.Where("user_id = ?", userID)
	}
	res := query.Count(&count)
	if res.Error != nil {
		return 0, res.Error
	}
	return count, nil
}

func (sm *ServiceManager) GetAllCompanyApplications(tx *gorm.DB, companyID uint) ([]database.JobApplication, error) {
	if sm.db == nil {
		return nil, errors.New("database not initialized")
	}
	db := sm.transactionOrDefault(tx)
	var companies database.Company
	result := db.Preload("JobApplications").Where("id = ?", companyID).Find(&companies)
	if result.Error != nil {
		return nil, result.Error
	}
	return companies.JobApplications, nil
}
