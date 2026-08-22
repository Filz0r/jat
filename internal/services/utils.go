package services

import (
	"github.com/filz0r/jat/internal/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func websiteEqual(a, b *string) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil || b == nil:
		return false
	default:
		return *a == *b
	}
}

func recordCompanyChange(
	tx *gorm.DB,
	companyID uint,
	oldName, newName string,
	oldURL, newURL *string,
	changedBy uuid.UUID,
) error {
	if oldName == newName && websiteEqual(oldURL, newURL) {
		return nil
	}

	entry := database.CompanyChangeHistory{
		CompanyID:       companyID,
		ChangedBy:       changedBy,
		OldNameValue:    &oldName,
		NewNameValue:    &newName,
		OldWebsiteValue: oldURL,
		NewWebsiteValue: newURL,
	}
	res := tx.Create(&entry)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (sm *ServiceManager) transactionOrDefault(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return sm.db
}
