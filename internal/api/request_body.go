package api

import "time"

// ---------------------------------------------//
//					 User Structs				//
// ---------------------------------------------//

type userCreateRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Username string `json:"username" validate:"required"`
}

type userLoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// ---------------------------------------------//
//			Job Application Structs				//
// ---------------------------------------------//

type applicationRequest struct {
	Title     string    `json:"title" validate:"required"`
	URL       string    `json:"url" validate:"required"`
	StatusID  int       `json:"status_id" validate:"required"`
	CompanyID int       `json:"company_id" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
}

// ---------------------------------------------//
//		Job Application Notes Structs			//
// ---------------------------------------------//

type noteCreateRequest struct {
	Body string `json:"body" validate:"required"`
}

// ---------------------------------------------//
//			Application Status Structs			//
// ---------------------------------------------//

type applicationStatusRequest struct {
	Status string `json:"status" validate:"required"`
	Kind   string `json:"kind" validate:"required"`
}

// ---------------------------------------------//
//				Company Structs					//
// ---------------------------------------------//

type companyBodyRequest struct {
	Name    string `json:"name" validate:"required"`
	Website string `json:"website,omitempty"`
}
