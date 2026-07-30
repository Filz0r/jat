package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/filz0r/jat/internal/auth"
	"github.com/filz0r/jat/internal/database"
	"github.com/google/uuid"
)

type userCreateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type userLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
}

type userCreateResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
	IsAdmin   bool      `json:"is_admin,omitempty"`
}

const jwtLifetime = time.Minute * 5              // 5 minutes
const refreshTokenLifetime = time.Hour * 24 * 60 // 60 days

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
		dbUser, err := s.services.CreateUser(database.User{
			Email:    user.Email,
			Username: user.Username,
			Password: user.Password,
		})
		if err != nil {
			s.respondWithError(w, 400, "error creating user", err)
			return
		}
		err = s.services.CreateInitialApplicationStatus(dbUser.ID)
		if err != nil {
			s.respondWithError(w, 400, "error creating initial application status", err)
		}
		response := userCreateResponse{
			UserID:    dbUser.ID,
			Email:     dbUser.Email,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Username:  dbUser.Username,
		}
		s.respondWithJSON(w, 201, apiResponse{
			Ok:      true,
			Data:    response,
			Message: "user created",
		})
	}
}

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

		token, err := auth.MakeJWT(dbUser.ID, *s.cfg.SecretJWT, jwtLifetime)
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
		}
		err = s.services.UpdateRefreshToken(refreshRecord.Token, time.Now())
		if err != nil {
			s.respondWithError(w, 401, "no valid token found", err)
			return
		}
		newToken, err := auth.MakeJWT(refreshRecord.UserID, *s.cfg.SecretJWT, jwtLifetime)
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
		return
	}
}

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
		for i := range users {
			users[i].Password = ""
		}
		s.respondWithJSON(w, 200, apiResponse{
			Ok:   true,
			Data: users,
		})

	}
}

func (s *Server) handleUserLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientType, ok := clientTypeFromContext(r.Context())
		if !ok {
			s.respondWithError(w, 401, "no client type found", nil)
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
		}

		s.respondWithJSON(w, 200, apiResponse{Ok: true, Message: "refresh token revoked"})
	}
}

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
		}
		user, err := s.services.GetUserByID(paramUUID)
		if err != nil {
			s.respondWithError(w, 404, "user not found", err)
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

		dbUser, err := s.services.UpdateUser(database.User{
			ID:       userID,
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
		})
		if err != nil {
			s.respondWithError(w, 400, "could not update user", err)
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

func (s *Server) handleMakeUserAdmin(give bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := r.PathValue("userID")
		paramUUID, err := uuid.Parse(param)
		if err != nil {
			s.respondWithError(w, 404, "user not found", nil)
			return
		}
		_, err = s.services.UpdateUser(database.User{
			ID:      paramUUID,
			IsAdmin: give,
		})
		if err != nil {
			s.respondWithError(w, 400, "could not make user admin", err)
			return
		}
		s.respondWithJSON(w, 200, apiResponse{Ok: true, Message: "user with id: '" + param + "' created"})
	}
}
