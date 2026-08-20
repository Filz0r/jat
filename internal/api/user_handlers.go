package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/filz0r/jat/internal/auth"
	"github.com/google/uuid"
)

type userCreateRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Username string `json:"username" validate:"required"`
}

type userLoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	UserID       string `json:"user_id" validate:"required"`
	Email        string `json:"email" validate:"required"`
}

type userCreateResponse struct {
	UserID          uuid.UUID `json:"user_id" validate:"required"`
	Email           string    `json:"email" validate:"required"`
	CreatedAt       time.Time `json:"created_at" validate:"required"`
	UpdatedAt       time.Time `json:"updated_at" validate:"required"`
	Username        string    `json:"username" validate:"required"`
	IsAdmin         bool      `json:"is_admin,omitempty"`
	DefaultStatusID uint      `json:"default_status_id,omitempty"`
}

const jwtLifetime = time.Minute * 5              // 5 minutes
const refreshTokenLifetime = time.Hour * 24 * 60 // 60 days

// @Summary Create user
// @Description Creates the first user. If no admin exists and the service is not initialized, the new user becomes the first admin.
// @Tags users
// @Accept json
// @Produce json
// @Param request body userCreateRequest true "User creation payload"
// @Success 201 {object} apiResponse{data=userCreateResponse}
// @Failure 400 {object} apiResponse
// @Router /users [post]
func (s *Server) handleUserCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		user := userCreateRequest{}
		err := decoder.Decode(&user)
		if err != nil {
			s.respondWithError(w, 400, "error parsing json", err)
			return
		}
		if user.Username == "" || user.Password == "" || user.Email == "" {
			s.respondWithError(w, 400, "invalid payload", nil)
			return
		}
		dbUser, err := s.services.CreateUser(
			user.Username,
			user.Email,
			user.Password,
		)
		if err != nil {
			s.respondWithError(w, 400, "error creating user", err)
			return
		}
		err = s.services.CreateInitialApplicationStatus(dbUser.ID)
		if err != nil {
			s.respondWithError(w, 400, "error creating initial application status", err)
			return
		}
		// Create first admin if it doesn't exist AND the server isn't at a initialized state
		if s.services.GetFirstAdmin() == uuid.Nil && !s.services.IsInitialized() {
			set := s.services.SetFirstAdmin(dbUser.ID)
			if !set {
				s.respondWithError(w, 400, "error setting first admin", nil)
				return
			}
			err := s.services.ChangeUserAdminStatus(dbUser.ID, true)
			if err != nil {
				s.respondWithError(w, 400, "error changing user admin status", err)
				return
			}
		}
		response := userCreateResponse{
			UserID:    dbUser.ID,
			Email:     dbUser.Email,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Username:  dbUser.Username,
		}
		if dbUser.DefaultApplicationStatusID != nil {
			response.DefaultStatusID = *dbUser.DefaultApplicationStatusID
		}
		s.respondWithJSON(w, 201, apiResponse{
			Ok:      true,
			Data:    response,
			Message: "user created",
		})
	}
}

// @Summary Login
// @Description Authenticates a user. For web clients the response includes two Set-Cookie headers: access_token (HttpOnly; Secure; SameSite=Strict; Path=/api; Max-Age=300) and refresh_token (HttpOnly; Secure; SameSite=Strict; Path=/api/auth; Max-Age=5184000). TUI clients receive the tokens in the response body. The X-Jat-Client-Type header is required by the server middleware but is set automatically by the generated web client.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body userLoginRequest true "Login credentials"
// @Success 200 {object} apiResponse{data=loginResponse}
// @Header 200 {string} Set-Cookie "access_token=<jwt>; HttpOnly; Secure; SameSite=Strict; Path=/api; Max-Age=300 AND refresh_token=<token>; HttpOnly; Secure; SameSite=Strict; Path=/api/auth; Max-Age=5184000"
// @Failure 400 {object} apiResponse
// @Failure 401 {object} apiResponse
// @Router /auth/login [post]
func (s *Server) handleUserLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientType, ok := clientTypeFromContext(r.Context())
		if !ok {
			s.respondWithError(w, 400, "missing client header", nil)
			return
		}
		decoder := json.NewDecoder(r.Body)
		user := userLoginRequest{}
		err := decoder.Decode(&user)
		if err != nil {
			s.respondWithError(w, 400, "error parsing json", err)
			return
		}
		dbUser, err := s.services.GetUserByEmail(user.Email)
		if err != nil {
			s.respondWithError(w, 401, "email or password are incorrect", err)
			return
		}
		success, err := auth.CheckPasswordHash(user.Password, dbUser.Password)
		if err != nil {
			s.respondWithError(w, 401, "email or password are incorrect", err)
			return
		}
		if !success {
			s.respondWithError(w, 401, "email or password are incorrect", nil)
			return
		}

		token, err := auth.MakeJWT(dbUser.ID, s.jwtSecret, jwtLifetime)
		if err != nil {
			s.respondWithError(w, 401, "email or password are incorrect", err)
			return
		}

		refreshToken, err := s.services.CreateRefreshToken(dbUser.ID, refreshTokenLifetime)
		if err != nil {
			s.respondWithError(w, 401, "email or password are incorrect", err)
			return
		}

		if clientType == webClient {
			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    token,
				Path:     "/api",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
				MaxAge:   int(jwtLifetime / time.Second),
			})

			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    refreshToken.Token,
				Path:     "/api/auth",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
				MaxAge:   int(refreshTokenLifetime / time.Second),
			})
			response := apiResponse{
				Data: loginResponse{
					UserID: dbUser.ID.String(),
					Email:  dbUser.Email,
				},
				Ok: true,
			}
			// Still return minimal info so the UI can know who logged in
			s.respondWithJSON(w, 200, response)
			return
		}

		response := apiResponse{
			Data: loginResponse{
				Token:        token,
				RefreshToken: refreshToken.Token,
				UserID:       dbUser.ID.String(),
				Email:        dbUser.Email,
			},
			Ok: true,
		}
		s.respondWithJSON(w, 200, response)
	}
}

// @Summary Refresh access token
// @Description Issues a new access token from a valid refresh token. Web clients get a new access_token cookie (path /api, MaxAge=300); TUI clients get the token in the response body. The X-Jat-Client-Type header is required by the server middleware but is set automatically by the generated web client.
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer refresh token (TUI clients)"
// @Success 201 {object} apiResponse{data=string}
// @Header 201 {string} Set-Cookie "access_token=<jwt>; HttpOnly; Secure; SameSite=Strict; Path=/api; Max-Age=300"
// @Failure 401 {object} apiResponse
// @Router /auth/refresh_token [get]
func (s *Server) handleUserTokenRefresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientType, ok := clientTypeFromContext(r.Context())
		if !ok {
			s.respondWithError(w, 401, "no valid token found", nil)
			return
		}
		refreshToken, ok := refreshTokenFromContext(r.Context())
		if !ok {
			s.respondWithError(w, 401, "no valid token found", nil)
			return
		}
		refreshRecord, err := s.services.GetValidRefreshToken(refreshToken)
		if err != nil {
			s.respondWithError(w, 401, "no valid token found", err)
			return
		}
		if time.Now().After(refreshRecord.ExpiresAt) {
			err := s.services.RevokeRefreshToken(refreshToken)
			if err != nil {
				s.respondWithError(w, 401, "no valid token found", err)
				return
			}
			s.respondWithError(w, 401, "refresh token is expired", nil)
			return
		}
		err = s.services.UpdateRefreshToken(refreshRecord.Token, time.Now())
		if err != nil {
			s.respondWithError(w, 401, "no valid token found", err)
			return
		}
		newToken, err := auth.MakeJWT(refreshRecord.UserID, s.jwtSecret, jwtLifetime)
		if err != nil {
			s.respondWithError(w, 401, "no valid token found", err)
			return

		}
		res := apiResponse{
			Ok:      true,
			Message: "access token updated",
		}
		if clientType == webClient {
			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    newToken,
				Path:     "/api",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
				MaxAge:   int(jwtLifetime / time.Second),
			})

		} else {
			res.Data = newToken
		}
		s.respondWithJSON(
			w,
			201,
			res,
		)
	}
}

// @Summary List all users
// @Description Admin-only endpoint that returns every user.
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} apiResponse{data=[]userCreateResponse}
// @Failure 400 {object} apiResponse
// @Failure 403 {object} apiResponse
// @Router /admin/users [get]
func (s *Server) handleGetAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := userAdminFromContext(r.Context())
		if !ok || !isAdmin {
			s.respondWithError(w, 403, "forbidden", nil)
			return
		}
		users, err := s.services.GetAllUsers()
		if err != nil {
			s.respondWithError(w, 400, "could not fetch users", err)
			return
		}
		response := make([]userCreateResponse, 0, len(users))
		for _, user := range users {
			temp := userCreateResponse{
				UserID:    user.ID,
				Email:     user.Email,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Username:  user.Username,
				IsAdmin:   user.IsAdmin,
			}
			response = append(response, temp)
		}
		s.respondWithJSON(w, 200, apiResponse{
			Ok:   true,
			Data: response,
		})

	}
}

// @Summary Logout
// @Description Revokes the current refresh token. For web clients the refresh_token cookie is cleared by setting Max-Age=-1 on the response. The X-Jat-Client-Type header is required by the server middleware but is set automatically by the generated web client.
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer refresh token (TUI clients)"
// @Success 200 {object} apiResponse
// @Header 200 {string} Set-Cookie "refresh_token=; HttpOnly; Secure; SameSite=Strict; Path=/api/auth; Max-Age=-1"
// @Failure 401 {object} apiResponse
// @Router /auth/logout [post]
func (s *Server) handleUserLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientType, ok := clientTypeFromContext(r.Context())
		if !ok {
			s.respondWithError(w, 401, "no client type found", nil)
			return
		}
		refreshToken, ok := refreshTokenFromContext(r.Context())
		if !ok {
			s.respondWithError(w, 401, "invalid refresh token found", nil)
			return
		}
		_, err := s.services.GetValidRefreshToken(refreshToken)
		if err != nil {
			s.respondWithError(w, 401, "invalid refresh token found", err)
			return
		}
		err = s.services.RevokeRefreshToken(refreshToken)
		if err != nil {
			s.respondWithError(w, 401, "invalid refresh token found", err)
			return
		}

		if clientType == webClient {
			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    "",
				Path:     "/api/auth",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})

			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    "",
				Path:     "/api",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})
		}

		s.respondWithJSON(w, 200, apiResponse{Ok: true, Message: "refresh token revoked"})
	}
}

// @Summary Revoke refresh token
// @Description Revokes the refresh token passed in the request. For web clients the refresh_token cookie is cleared by setting Max-Age=-1 on the response. The X-Jat-Client-Type header is required by the server middleware but is set automatically by the generated web client.
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer refresh token (TUI clients)"
// @Success 200 {object} apiResponse
// @Header 200 {string} Set-Cookie "refresh_token=; HttpOnly; Secure; SameSite=Strict; Path=/api/auth; Max-Age=-1"
// @Failure 401 {object} apiResponse
// @Router /auth/revoke_token [get]
func (s *Server) handleUserRevokeToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientType, ok := clientTypeFromContext(r.Context())
		if !ok {
			s.respondWithError(w, 401, "no client type found", nil)
			return
		}
		refreshToken, ok := refreshTokenFromContext(r.Context())
		if !ok {
			s.respondWithError(w, 401, "no token found", nil)
			return
		}
		refreshRecord, err := s.services.GetValidRefreshToken(refreshToken)
		if err != nil {
			s.respondWithError(w, 401, "no token found", err)
			return
		}
		err = s.services.RevokeRefreshToken(refreshRecord.Token)
		if err != nil {
			s.respondWithError(w, 401, "no token found", err)
			return
		}
		if clientType == webClient {
			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    "",
				Path:     "/api/auth",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})
		}
		s.respondWithJSON(w, 200, apiResponse{Ok: true, Message: "refresh token revoked"})
	}
}

// @Summary Get a user
// @Description Returns a single user. Users can read their own record; admins can read any record.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userID path string true "User UUID"
// @Success 200 {object} apiResponse{data=userCreateResponse}
// @Failure 401 {object} apiResponse
// @Failure 403 {object} apiResponse
// @Failure 404 {object} apiResponse
// @Router /users/{userID} [get]
func (s *Server) handleGetSingleUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromContext(r.Context())
		if !ok {
			s.respondWithError(w, 401, "no user id found", nil)
			return
		}
		param := r.PathValue("userID")
		paramUUID, err := uuid.Parse(param)
		if err != nil {
			s.respondWithError(w, 404, "user not found", nil)
			return
		}
		// only admins can check other users
		isAdmin := s.services.IsUserAdmin(userID)
		if paramUUID != userID && !isAdmin {
			s.respondWithError(w, 403, "forbidden", nil)
			return
		}
		user, err := s.services.GetUserByID(paramUUID)
		if err != nil {
			s.respondWithError(w, 404, "user not found", err)
			return
		}
		response := apiResponse{
			Ok:      true,
			Message: "user found",
			Data: userCreateResponse{
				UserID:    user.ID,
				IsAdmin:   user.IsAdmin,
				Email:     user.Email,
				Username:  user.Username,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			},
		}
		s.respondWithJSON(w, 200, response)
	}
}

// @Summary Update current user
// @Description Updates the authenticated user's profile.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userCreateRequest true "Updated user fields"
// @Success 200 {object} apiResponse{data=userCreateResponse}
// @Failure 400 {object} apiResponse
// @Failure 401 {object} apiResponse
// @Router /users [put]
func (s *Server) handleUserUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		user := userCreateRequest{}
		err := decoder.Decode(&user)
		if err != nil {
			s.respondWithError(w, 400, "invalid body", err)
			return
		}
		userID, ok := userIDFromContext(r.Context())
		if !ok {
			s.respondWithError(w, 401, "no user id found", nil)
			return
		}

		dbUser, err := s.services.UpdateUser(
			userID,
			user.Username,
			user.Password,
			user.Email,
		)
		if err != nil {
			s.respondWithError(w, 400, "could not update user", err)
			return
		}
		response := apiResponse{
			Ok:      true,
			Message: "user updated",
			Data: userCreateResponse{
				UserID:   dbUser.ID,
				Email:    dbUser.Email,
				Username: dbUser.Username,
			},
		}
		s.respondWithJSON(w, 200, response)
	}
}

// @Summary Toggle admin status
// @Description Admin-only endpoint to grant (GET) or revoke (DELETE) admin privileges for a user.
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param userID path string true "User UUID"
// @Success 200 {object} apiResponse
// @Failure 400 {object} apiResponse
// @Failure 403 {object} apiResponse
// @Failure 404 {object} apiResponse
// @Router /admin/users/{userID} [get]
func (s *Server) handleMakeUserAdmin(give bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := r.PathValue("userID")
		paramUUID, err := uuid.Parse(param)
		if err != nil {
			s.respondWithError(w, 404, "user not found", nil)
			return
		}
		err = s.services.ChangeUserAdminStatus(paramUUID, give)
		if err != nil {
			s.respondWithError(w, 400, "could not change admin status for this user", err)
			return
		}
		message := "User with Id: " + param + " is now admin"
		if !give {
			message = "User with Id: " + param + " is no longer admin"
		}
		s.respondWithJSON(w, 200, apiResponse{Ok: true, Message: message})
	}
}
