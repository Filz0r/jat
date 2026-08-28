package services

import (
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

func (sm *ServiceManager) transactionOrDefault(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return sm.db
}
