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
//			Job Application Structs				//
// ---------------------------------------------//

type applicationResponse struct {
	ID        uint                      `json:"id" validate:"required"`
	Title     string                    `json:"title" validate:"required"`
	URL       string                    `json:"url" validate:"required"`
	UserID    uuid.UUID                 `json:"user_id" validate:"required"`
	CreatedAt time.Time                 `json:"created_at" validate:"required"`
	UpdatedAt time.Time                 `json:"updated_at" validate:"required"`
	Status    applicationStatusResponse `json:"status" validate:"required"`
	Company   companyResponse           `json:"company" validate:"required"`
}

// ---------------------------------------------//
//		Job Application Notes Structs			//
// ---------------------------------------------//

type noteDataResponse struct {
	ID        uint                      `json:"id" validate:"required"`
	Body      string                    `json:"body" validate:"required"`
	Status    applicationStatusResponse `json:"status" validate:"required"`
	UserID    uuid.UUID                 `json:"user_id" validate:"required"`
	JobID     uint                      `json:"job_id" validate:"required"`
	CreatedAt time.Time                 `json:"created_at" validate:"required"`
	UpdatedAt time.Time                 `json:"updated_at" validate:"required"`
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
//				Company Structs					//
// ---------------------------------------------//

type companyResponse struct {
	ID            uint              `json:"id,omitempty" validate:"required"`
	Name          string            `json:"name" validate:"required"`
	Website       string            `json:"website,omitempty"`
	CreatedAt     time.Time         `json:"created_at" validate:"required"`
	UpdatedAt     time.Time         `json:"updated_at" validate:"required"`
	UserCount     int64             `json:"user_count,omitempty"`
	TotalCount    int64             `json:"total_count,omitempty"`
	CreatedByUser *baseUserResponse `json:"created_by_user,omitempty"`
	UpdatedByUser *baseUserResponse `json:"updated_by_user,omitempty"`
}

type companyChangeHistoryResponse struct {
	ID         uint             `json:"id" validate:"required"`
	CompanyID  uint             `json:"company_id" validate:"required"`
	NewName    string           `json:"new_name,omitempty"`
	OldWebsite string           `json:"old_website,omitempty"`
	NewWebsite string           `json:"new_website,omitempty"`
	OldName    string           `json:"old_name,omitempty"`
	CreatedAt  time.Time        `json:"created_at" validate:"required"`
	ChangedBy  baseUserResponse `json:"changed_by" validate:"required"`
	Reverted   bool             `json:"reverted" validate:"required"`
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
//			Job Application Factories			//
// ---------------------------------------------//

func newApplicationResponse(row database.JobApplication) applicationResponse {
	res := applicationResponse{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		UserID:    row.UserID,
		Title:     row.Title,
		URL:       row.Url,
		Company:   newCompanyResponse(row.Company, 0, 0, false, nil, nil),
		Status:    newApplicationStatusResponse(row.Status),
	}
	if row.Company.Website != nil {
		res.Company.Website = *row.Company.Website
	}
	return res
}

// ---------------------------------------------//
//		Job Application Notes Factories			//
// ---------------------------------------------//

func newNoteDataResponse(data database.ApplicationNote) noteDataResponse {
	return noteDataResponse{
		ID:        data.ID,
		Body:      data.Body,
		Status:    newApplicationStatusResponse(data.Status),
		UserID:    data.UserID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		JobID:     data.ApplicationID,
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

// ---------------------------------------------//
//				Company Factories				//
// ---------------------------------------------//

func newCompanyResponse(data database.Company, totalCount, userCount int64, includeCounts bool, cbUser, ebUser *database.User) companyResponse {
	result := companyResponse{
		ID:        data.ID,
		Name:      data.Name,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
	if data.Website != nil {
		result.Website = *data.Website
	}
	if includeCounts {
		result.TotalCount = totalCount
		result.UserCount = userCount
	}
	if cbUser != nil {
		userResponse := newBaseUserResponse(*cbUser)
		result.CreatedByUser = &userResponse
	}
	if ebUser != nil {
		userResponse := newBaseUserResponse(*ebUser)
		result.UpdatedByUser = &userResponse
	}
	return result
}

func newCompanyChangeHistoryResponse(data database.CompanyChangeHistory) companyChangeHistoryResponse {
	result := companyChangeHistoryResponse{
		ID:        data.ID,
		CompanyID: data.CompanyID,
		CreatedAt: data.CreatedAt,
		Reverted:  data.Reverted,
		ChangedBy: newBaseUserResponse(data.ChangedByUser),
	}
	if data.OldNameValue != nil {
		result.OldName = *data.OldNameValue
	}
	if data.NewNameValue != nil {
		result.NewName = *data.NewNameValue
	}
	if data.OldWebsiteValue != nil {
		result.OldWebsite = *data.OldWebsiteValue
	}
	if data.NewWebsiteValue != nil {
		result.NewWebsite = *data.NewWebsiteValue
	}
	return result
}
