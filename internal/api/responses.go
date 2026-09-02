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
