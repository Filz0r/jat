package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/filz0r/jat/internal/auth"
)

const jwtLifetime = time.Minute * 5              // 5 minutes
const refreshTokenLifetime = time.Hour * 24 * 60 // 60 days

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
				Ok:      true,
				Message: "User logged in",
			}
			// Still return minimal info so the UI can know who logged in
			s.respondWithJSON(w, 200, response)
			return
		}

		response := apiResponse{
			Data:    newLoginResponse(token, refreshToken.Token),
			Message: "User logged in",
			Ok:      true,
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
			Message: "Access token updated",
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
		clientType, _ := clientTypeFromContext(r.Context())

		refreshToken, _ := refreshTokenFromContext(r.Context())

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

		s.respondWithJSON(w, 200, apiResponse{
			Ok:      true,
			Message: "User logged out",
		})
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
		s.respondWithJSON(w, 200, apiResponse{
			Ok:      true,
			Message: "Refresh token revoked",
		})
	}
}
