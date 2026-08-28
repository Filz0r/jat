package services

import (
	"errors"

	"github.com/filz0r/jat/internal/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func recordCompanyChange(
	tx *gorm.DB,
	companyID uint,
	oldName, newName string,
	oldURL, newURL *string,
	changedBy uuid.UUID,
) error {
	if oldName == newName && websiteEqual(oldURL, newURL) {
		return errors.New("nothing changed")
	}

	entry := database.CompanyChangeHistory{
		CompanyID: companyID,
		ChangedBy: changedBy,
	}
	if newName != oldName {
		entry.OldNameValue = &oldName
		entry.NewNameValue = &newName
	}
	if !websiteEqual(oldURL, newURL) {
		entry.OldWebsiteValue = oldURL
		entry.NewWebsiteValue = newURL
	}
	res := tx.Create(&entry)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (sm *ServiceManager) GetCompanyHistoryEntry(
	tx *gorm.DB,
	userID uuid.UUID,
	historyID uint,
) (database.CompanyChangeHistory, error) {
	if sm.db == nil {
		return database.CompanyChangeHistory{}, errors.New("database not initialized")
	}
	db := sm.transactionOrDefault(tx)
	if !sm.IsUserAdmin(userID) {
		return database.CompanyChangeHistory{}, errors.New("user not admin")
	}
	var entry database.CompanyChangeHistory
	result := db.Preload("Company").
		Where("id = ?", historyID).
		Find(&entry)
	if result.Error != nil {
		return database.CompanyChangeHistory{}, result.Error
	}
	if result.RowsAffected == 0 {
		return database.CompanyChangeHistory{}, errors.New("history not found")
	}
	return entry, nil
}

func (sm *ServiceManager) GetCompanyHistory(
	tx *gorm.DB,
	userID uuid.UUID,
	companyID uint,
) ([]database.CompanyChangeHistory, error) {
	if sm.db == nil {
		return nil, errors.New("database not initialized")
	}

	if !sm.IsUserAdmin(userID) {
		return nil, errors.New("only admin users can change companies")
	}
	db := sm.transactionOrDefault(tx)
	var history []database.CompanyChangeHistory
	result := db.
		Preload("ChangedByUser").
		Where("company_id = ?", companyID).
		Order("created_at DESC").
		Find(&history)
	if result.Error != nil {
		return nil, result.Error
	}
	return history, nil
}

func (sm *ServiceManager) RevertCompanyChange(
	tx *gorm.DB,
	userID uuid.UUID,
	changeID,
	companyID uint,
) error {
	if sm.db == nil {
		return errors.New("database not initialized")
	}
	db := sm.transactionOrDefault(tx)
	if !sm.IsUserAdmin(userID) {
		return errors.New("only admin users can revert company changes")
	}
	err := db.Transaction(func(_tx *gorm.DB) error {
		entry, err := sm.GetCompanyHistoryEntry(_tx, userID, changeID)
		if err != nil {
			return err
		}
		if entry.Reverted {
			return errors.New("this change is already reverted")
		}
		if entry.CompanyID != companyID {
			return errors.New("company id does not match")
		}
		company := entry.Company
		if entry.NewNameValue != nil && entry.OldNameValue != nil {
			company.Name = *entry.OldNameValue
		}
		if entry.NewWebsiteValue != nil {
			// OldWebsiteValue can be restored to nil if the website is previously nil, this is by design
			company.Website = entry.OldWebsiteValue
		}

		var websiteURL string
		if company.Website != nil {
			websiteURL = *company.Website
		} else {
			websiteURL = ""
		}
		_, err = sm.UpdateCompany(_tx, company.ID, userID, company.Name, websiteURL)

		if err != nil {
			return err
		}

		err = _tx.Model(&entry).Update("reverted", true).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
