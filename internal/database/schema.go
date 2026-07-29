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

	CreatedBy     uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedByUser User      `gorm:"foreignKey:CreatedBy"`

	EditedBy     uuid.UUID `gorm:"type:uuid;not null;index"`
	EditedByUser User      `gorm:"foreignKey:EditedBy"`

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

type RefreshToken struct {
	Token  string    `gorm:"primaryKey"`
	UserID uuid.UUID `gorm:"index"`
	User   User      `gorm:"constraints:OnDelete:CASCADE;foreignKey:UserID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
	ExpiredAt *time.Time
}

type ApplicationNote struct {
	gorm.Model
	Body          string
	ApplicationID uint
	Application   JobApplication `gorm:"constraints:OnDelete:CASCADE;foreignKey:ApplicationID"`

	StatusID *uint
	Status   *ApplicationStatus `gorm:"foreignKey:StatusID"`
}

type StatusHistory struct {
	ID            uint `gorm:"primaryKey"`
	ApplicationID uint
	Application   JobApplication `gorm:"constraints:OnDelete:CASCADE;foreignKey:ApplicationID"`

	StatusID uint
	Status   ApplicationStatus `gorm:"foreignKey:StatusID"`

	CreatedAt time.Time
}

type CompanyChangeHistory struct {
	ID        uint    `gorm:"primaryKey"`
	CompanyID uint    `gorm:"not null;index"`
	Company   Company `gorm:"constraints:OnDelete:CASCADE;foreignKey:CompanyID"`

	OldNameValue    *string // nullable, old value may be null
	NewNameValue    *string // nullable, new value may be null
	OldWebsiteValue *string // nullable, old value may be null
	NewWebsiteValue *string // nullable, new value may be null

	ChangedBy     uuid.UUID `gorm:"type:uuid;not null;index"`
	ChangedByUser User      `gorm:"foreignKey:ChangedBy"`

	CreatedAt time.Time
	Reverted  bool `gorm:"default:false"`
}
