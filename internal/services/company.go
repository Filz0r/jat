package services

import (
	"errors"
	"strconv"

	"github.com/filz0r/jat/internal/database"
	"github.com/filz0r/jat/internal/utils"
)

func (sm *ServiceManager) CreateCompany(company database.Company) (database.Company, error) {
	if sm.db == nil {
		return database.Company{}, errors.New("database not initialized")
	}
	result := sm.db.Create(&company)
	return company, result.Error
}

func (sm *ServiceManager) UpdateCompany(company database.Company) (database.Company, error) {
	if sm.db == nil {
		return database.Company{}, errors.New("database not initialized")
	}
	result := sm.db.Save(&company)
	return company, result.Error
}

func (sm *ServiceManager) GetCompany(name string) (database.Company, error) {
	if sm.db == nil {
		return database.Company{}, errors.New("database not initialized")
	}
	var company database.Company
	result := sm.db.Where("name = ?", name).First(&company)
	return company, result.Error
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
	return result.Error
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
