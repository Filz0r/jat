package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Email    string `gorm:"uniqueIndex"`
	Username string `gorm:"uniqueIndex"`
	Password string
	IsAdmin  bool `gorm:"default:false"`

	StatusList      []ApplicationStatus `gorm:"foreignKey:UserID"`
	JobApplications []JobApplication    `gorm:"foreignKey:UserID"`
}

type ApplicationStatus struct {
	gorm.Model
	Status string `gorm:"index"`
	UserID uuid.UUID
	User   User `gorm:"constraints:OnDelete:CASCADE;foreignKey:UserID"`
}

type Company struct {
	gorm.Model
	Name    string `gorm:"uniqueIndex"`
	Website *string

	JobApplications []JobApplication `gorm:"constraints:OnDelete:CASCADE;many2many:company_application;joinForeignKey:company_id;joinReferences:application_id"`
}

type JobApplication struct {
	gorm.Model
	Title  string
	Url    string
	UserID uuid.UUID
	User   User `gorm:"foreignKey:UserID"`

	StatusID uint
	Status   ApplicationStatus `gorm:"foreignKey:StatusID"`

	Companies []Company `gorm:"constraints:OnDelete:CASCADE;many2many:company_application;joinForeignKey:application_id;joinReferences:company_id"`
}
