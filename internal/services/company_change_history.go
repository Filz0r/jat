package services

import (
	"errors"

	"github.com/filz0r/jat/internal/database"
)

// TODO: This needs to be made so that only admins can track these changes
func (sm *ServiceManager) GetCompanyChangeHistory(
	companyID uint,
) ([]database.CompanyChangeHistory, error) {
	if sm.db == nil {
		return []database.CompanyChangeHistory{}, errors.New("database not initialized")
	}
	var history []database.CompanyChangeHistory
	result := sm.db.
		Where("company_id = ?", companyID).
		Preload("ChangedByUser").
		Order("created_at DESC").
		Find(&history)
	return history, result.Error
}
