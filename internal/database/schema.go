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

	UserSettingsID      uint                `gorm:"index"`
	UserSettings        UserSettings        `gorm:"foreignKey:UserSettingsID;constraints:OnDelete:CASCADE"`
	BanList             []BanList           `gorm:"foreignKey:UserID;constraints:OnDelete:CASCADE;-:migration"`
	StatusList          []ApplicationStatus `gorm:"foreignKey:UserID;-:migration"`
	JobApplications     []JobApplication    `gorm:"foreignKey:UserID;-:migration"`
	JobApplicationNotes []ApplicationNote   `gorm:"foreignKey:UserID;-:migration"`
}

type UserSettings struct {
	gorm.Model
	UserID                     uuid.UUID          `gorm:"type:uuid;not null"`
	SetupStep                  uint               `gorm:"default:1"` //0 means the account has been set up
	DefaultApplicationStatusID *uint              `gorm:"index"`
	DefaultApplicationStatus   *ApplicationStatus `gorm:"-:migration"`
	IsBanned                   bool               `gorm:"default:false"`
	IsEnabled                  bool               `gorm:"default:false"`
	IsAdmin                    bool               `gorm:"default:false"`
}

type BanList struct {
	gorm.Model
	UserID       uuid.UUID     `gorm:"type:uuid;not null"`
	BannedAt     time.Time     `gorm:"index;not null"`
	BanDuration  time.Duration `gorm:"not null"`
	BanReason    string        `gorm:"not null"`
	Active       bool          `gorm:"default:true"`
	BannedBy     uuid.UUID     `gorm:"not null"`
	BannedByUser User          `gorm:"foreignKey:BannedBy;constraints:OnDelete:CASCADE"`
}

type ApplicationStatus struct {
	gorm.Model
	Status   string `gorm:"index"`
	UserID   uuid.UUID
	User     User                  `gorm:"constraints:OnDelete:CASCADE;foreignKey:UserID"`
	Kind     ApplicationStatusKind `gorm:"type:string;not null"`
	Archived bool                  `gorm:"default:false;index;not null"`
}

type Company struct {
	gorm.Model
	Name    string `gorm:"uniqueIndex"`
	Website *string

	CreatedBy     uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedByUser User      `gorm:"foreignKey:CreatedBy"`

	EditedBy     uuid.UUID `gorm:"type:uuid;not null;index"`
	EditedByUser User      `gorm:"foreignKey:EditedBy"`

	JobApplications []JobApplication `gorm:"constraints:OnDelete:CASCADE;foreignKey:CompanyID;-:migration"`
}

type JobApplication struct {
	gorm.Model
	Title  string
	Url    string
	UserID uuid.UUID
	User   User `gorm:"foreignKey:UserID"`

	StatusID uint
	Status   ApplicationStatus `gorm:"foreignKey:StatusID"`

	CompanyID uint    `gorm:"not null,index"`
	Company   Company `gorm:"foreignKey:CompanyID;references:ID"`

	Notes []ApplicationNote `gorm:"foreignKey:ApplicationID"`
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

	StatusID uint
	Status   ApplicationStatus `gorm:"foreignKey:StatusID"`

	UserID uuid.UUID
	User   User `gorm:"foreignKey:UserID"`
}

type StatusHistory struct {
	ID            uint `gorm:"primaryKey"`
	ApplicationID uint
	Application   JobApplication `gorm:"constraints:OnDelete:CASCADE;foreignKey:ApplicationID"`
	NewStatusID   uint
	NewStatus     ApplicationStatus `gorm:"foreignKey:NewStatusID"`
	OldStatusID   *uint
	OldStatus     *ApplicationStatus `gorm:"foreignKey:OldStatusID"`
	CreatedAt     time.Time
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

type Config struct {
	Key       string     `gorm:"primaryKey;not null"`
	Value     string     `gorm:"not null"`
	Type      ConfigType `gorm:"type:string;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
