package services

import (
	"github.com/filz0r/jat/internal/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func companyChanged(current, updated database.Company) bool {
	if current.Name != updated.Name {
		return true
	}
	switch {
	case current.Website == nil && updated.Website == nil:
		return false
	case current.Website == nil || updated.Website == nil:
		return true
	default:
		return *current.Website != *updated.Website
	}
}

func recordCompanyChange(
	tx *gorm.DB,
	current, updated database.Company,
	changedBy uuid.UUID,
) error {
	if !companyChanged(current, updated) {
		return nil
	}

	entry := database.CompanyChangeHistory{
		CompanyID:       current.ID,
		ChangedBy:       changedBy,
		OldNameValue:    &current.Name,
		NewNameValue:    &updated.Name,
		OldWebsiteValue: current.Website,
		NewWebsiteValue: updated.Website,
	}
	return tx.Create(&entry).Error
}
