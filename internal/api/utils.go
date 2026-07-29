package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type serverClientTypes string

const (
	tuiClient serverClientTypes = "tui-client"
	webClient serverClientTypes = "web-client"
)

func (ct serverClientTypes) Valid() bool {
	if ct != tuiClient && ct != webClient {
		return false
	}
	return true
}

type contextKey string

const contextKeyUserID contextKey = "userID"
const contextKeyRefreshToken contextKey = "refreshToken"
const contextUserAdmin contextKey = "userAdmin"
const contextGetClientType contextKey = "clientType"

type errorResponse struct {
	Error string `json:"error,omitempty"`
	Ok    bool   `json:"ok,omitempty"`
}

type successResponse struct {
	Data    any    `json:"data,omitempty"`
	Ok      bool   `json:"ok,omitempty"`
	Message string `json:"message"`
}

func userIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(contextKeyUserID).(uuid.UUID)
	return id, ok
}

func refreshTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(contextKeyRefreshToken).(string)
	return token, ok
}

func userAdminFromContext(ctx context.Context) (bool, bool) {
	isAdmin, ok := ctx.Value(contextUserAdmin).(bool)
	return isAdmin, ok
}

func clientTypeFromContext(ctx context.Context) (serverClientTypes, bool) {
	clientType, ok := ctx.Value(contextGetClientType).(serverClientTypes)
	return clientType, ok
}

func (s *Server) respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		s.logger.Println(err)
	}
	if code > 499 {
		s.logger.Printf("Responding with 5XX error: %s", msg)
	}

	s.respondWithJSON(w, code, errorResponse{
		Error: msg,
		Ok:    false,
	})
}

func (s *Server) respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		s.logger.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		_, err := w.Write([]byte("Internal Server Error"))
		if err != nil {
			log.Printf("Error writing response: %s", err)
		}
		return
	}
	w.WriteHeader(code)
	_, err = w.Write(dat)
	if err != nil {
		log.Printf("Error writing response: %s", err)
	}
}
