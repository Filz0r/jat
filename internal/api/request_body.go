package api

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
