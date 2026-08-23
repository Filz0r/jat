package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"

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

type apiResponse struct {
	Data    any    `json:"data,omitempty"`
	Ok      bool   `json:"ok" validate:"required"`
	Message string `json:"message" validate:"required"`
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
	s.respondWithJSON(w, code, apiResponse{
		Message: msg,
		Ok:      false,
	})
}

func (s *Server) respondWithJSON(w http.ResponseWriter, code int, payload apiResponse) {
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

func BindQuery[T any](r *http.Request) (T, error) {
	var out T
	vals := r.URL.Query()
	v := reflect.ValueOf(&out).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Tag.Get("query")
		if name == "" {
			continue
		}
		raw := vals.Get(name)
		if raw == "" {
			continue
		}
		f := v.Field(i)
		switch f.Kind() {
		case reflect.Bool:
			b, err := strconv.ParseBool(raw)
			if err != nil {
				return out, fmt.Errorf("%s must be a boolean (got %q)", name, raw)
			}
			f.SetBool(b)
		case reflect.Int, reflect.Int64:
			n, err := strconv.Atoi(raw)
			if err != nil {
				return out, fmt.Errorf("invalid integer value: %s", raw)
			}
			f.SetInt(int64(n))

		case reflect.String:
			f.SetString(raw)
		default:
			return out, fmt.Errorf("unsupported query field type %s for %s", f.Kind(), name)
		}
	}
	return out, nil
}
