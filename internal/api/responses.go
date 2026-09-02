package api

import (
	"time"

	"github.com/filz0r/jat/internal/database"
	"github.com/google/uuid"
)

// ---------------------------------------------//
//				User Structs					//
// ---------------------------------------------//

type baseUserResponse struct {
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Username  string    `json:"username" validate:"required,username"`
	CreatedAt time.Time `json:"create_at" validate:"required"`
	UpdatedAt time.Time `json:"update_at" validate:"required"`
}

type meUserResponse struct {
	baseUserResponse
	IsAdmin         bool `json:"is_admin,omitempty"`
	DefaultStatusID uint `json:"default_status_id" validate:"required"`
	IsSetup         bool `json:"is_setup" validate:"required"`
	SetupStep       uint `json:"setup_step" validate:"required"`
}

type adminUserResponse struct {
	baseUserResponse
	IsAdmin   bool `json:"is_admin" validate:"required"`
	IsBanned  bool `json:"is_banned" validate:"required"`
	IsEnabled bool `json:"is_enabled" validate:"required"`
	IsSetup   bool `json:"is_setup" validate:"required"`
}

type loginResponse struct {
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// ---------------------------------------------//
//			Application Status Structs			//
// ---------------------------------------------//

type applicationStatusResponse struct {
	ID        uint      `json:"id" validate:"required"`
	Status    string    `json:"status" validate:"required"`
	Kind      string    `json:"kind" validate:"required"`
	UserID    uuid.UUID `json:"user_id,omitempty"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	Archived  bool      `json:"archived" validate:"required"`
}

type applicationStatusHistoryResponse struct {
	ID            uint                       `json:"id" validate:"required"`
	ApplicationID uint                       `json:"application_id" validate:"required"`
	OldStatus     *applicationStatusResponse `json:"old_status,omitempty"`
	NewStatus     applicationStatusResponse  `json:"new_status" validate:"required"`
	CreatedAt     time.Time                  `json:"created_at" validate:"required"`
}

// ---------------------------------------------//
//					 User Factories				//
// ---------------------------------------------//

func newBaseUserResponse(data database.User) baseUserResponse {
	return baseUserResponse{
		UserID:    data.ID,
		Email:     data.Email,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Username:  data.Username,
	}
}

func newMeUserResponse(data database.User) meUserResponse {
	result := meUserResponse{
		baseUserResponse: newBaseUserResponse(data),
		IsAdmin:          data.UserSettings.IsAdmin,
		SetupStep:        data.UserSettings.SetupStep,
	}

	if data.UserSettings.SetupStep == 0 {
		result.IsSetup = true
	} else {
		result.IsSetup = false
	}

	if data.UserSettings.DefaultApplicationStatusID != nil {
		result.DefaultStatusID = *data.UserSettings.DefaultApplicationStatusID
	}
	return result
}

func newAdminUserResponse(data database.User) adminUserResponse {
	result := adminUserResponse{
		baseUserResponse: newBaseUserResponse(data),
		IsAdmin:          data.UserSettings.IsAdmin,
		IsEnabled:        data.UserSettings.IsEnabled,
		IsSetup:          data.UserSettings.SetupStep == 0,
		IsBanned:         data.UserSettings.IsBanned,
	}
	return result
}

func newLoginResponse(token, refreshToken string) loginResponse {
	return loginResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}
}

// ---------------------------------------------//
//			Application Status Factories		//
// ---------------------------------------------//

func newApplicationHistoryResponse(d database.StatusHistory) applicationStatusHistoryResponse {
	response := applicationStatusHistoryResponse{
		ID:            d.ID,
		ApplicationID: d.ApplicationID,
		NewStatus:     newApplicationStatusResponse(d.NewStatus),
		CreatedAt:     d.CreatedAt,
	}

	if d.OldStatus != nil {
		old := newApplicationStatusResponse(*d.OldStatus)
		response.OldStatus = &old
	}

	return response
}

func newApplicationStatusResponse(d database.ApplicationStatus) applicationStatusResponse {
	return applicationStatusResponse{
		ID:        d.ID,
		Status:    d.Status,
		Kind:      d.Kind.String(),
		UpdatedAt: d.UpdatedAt,
		CreatedAt: d.CreatedAt,
		UserID:    d.UserID,
		Archived:  d.Archived,
	}
}
