package api

import (
	"context"
	"net/http"
	"time"

	"github.com/filz0r/jat/internal/auth"
)

func (s *Server) middlewareLogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := newResponseRecorder(w)

		next.ServeHTTP(recorder, r)

		duration := time.Since(start)
		s.logger.Printf(
			"%s %s -> %s%d%s in %s",
			r.Method,
			r.URL.Path,
			statusColor(recorder.status),
			recorder.status,
			reset,
			duration,
		)
	})
}

func (s *Server) middlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r)
		if err != nil {
			s.respondWithError(w, 401, "invalid jwt token", err)
			return
		}
		userID, err := auth.ValidateJWT(token, *s.cfg.SecretJWT)
		if err != nil {
			s.respondWithError(w, 401, "invalid jwt token", err)
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyUserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) middlewareRefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetAPIKey(r)
		if err != nil {
			s.logger.Printf("Error parsing token: %s", err)
			s.respondWithError(
				w,
				401,
				"Invalid user token",
				err,
			)
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyRefreshToken, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) middlewareAdminUser(next http.Handler) http.Handler {
	return s.middlewareAuth(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := userIDFromContext(r.Context())
			if !ok {
				s.respondWithError(w, 401, "invalid jwt token", nil)
				return
			}
			isAdmin := s.services.IsUserAdmin(userID)
			ctx := context.WithValue(r.Context(), contextUserAdmin, isAdmin)
			next.ServeHTTP(w, r.WithContext(ctx))
		}),
	)
}

func (s *Server) middlewareClientType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientType := r.Header.Get("X-Jat-Client-Type")
		if clientType == "" {
			s.respondWithError(w, 401, "missing client type", nil)
			return

		}
		converted := serverClientTypes(clientType)
		if !converted.Valid() {
			s.respondWithError(w, 401, "invalid client type", nil)
			return
		}
		ctx := context.WithValue(r.Context(), contextGetClientType, converted)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
