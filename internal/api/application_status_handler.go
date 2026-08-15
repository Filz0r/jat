package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/filz0r/jat/internal/database"
	"github.com/google/uuid"
)

type applicationStatusRequest struct {
	ID        uint      `json:"id" validate:"required"`
	Status    string    `json:"status" validate:"required"`
	Kind      string    `json:"kind" validate:"required"`
	UserID    uuid.UUID `json:"user_id,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// @Summary List current user's application statuses
// @Description Returns all application statuses belonging to the authenticated user.
// @Tags application_statuses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} apiResponse{data=[]applicationStatusRequest}
// @Failure 400 {object} apiResponse
// @Failure 401 {object} apiResponse
// @Router /application_statuses [get]
func (s *Server) handleGetUserApplicationStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromContext(r.Context())
		if !ok {
			s.respondWithError(w, 401, "no user id found", nil)
			return
		}
		data, err := s.services.GetAllApplicationStatus(userID)
		if err != nil {
			s.respondWithError(w, 400, "error getting application status", err)
			return
		}
		converted := make([]applicationStatusRequest, 0, len(data))
		for _, d := range data {
			temp := applicationStatusRequest{
				ID:        d.ID,
				Status:    d.Status,
				Kind:      d.Kind.String(),
				UpdatedAt: d.UpdatedAt,
				CreatedAt: d.CreatedAt,
				UserID:    d.UserID,
			}
			converted = append(converted, temp)
		}

		s.respondWithJSON(w, 200, apiResponse{Ok: true, Data: converted})
	}
}

func (s *Server) handleGetAllApplicationStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := s.services.GetAllApplicationStatusAdmin()
		if err != nil {
			s.respondWithError(w, 400, "error getting application status", err)
			return
		}
		// TODO: Might need to change this to user applicationStatusRequest instead of the db type
		s.respondWithJSON(w, 200, apiResponse{Ok: true, Data: data})
	}
}

// @Summary Update application status
// @Description Updates an existing application status. Only the owner or an admin can update it.
// @Tags application_statuses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param statusID path int true "Application status ID"
// @Param request body applicationStatusRequest true "Updated application status"
// @Success 200 {object} apiResponse{data=applicationStatusRequest}
// @Failure 400 {object} apiResponse
// @Failure 403 {object} apiResponse
// @Router /application_statuses/{statusID} [put]
func (s *Server) handleUpdateApplicationStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := r.PathValue("statusID")
		conv, err := strconv.ParseUint(param, 10, 32)
		if err != nil {
			s.respondWithError(w, 400, "invalid id format", err)
			return
		}
		statusID := uint(conv)
		userID, _ := userIDFromContext(r.Context())

		decoder := json.NewDecoder(r.Body)
		body := applicationStatusRequest{}
		err = decoder.Decode(&body)
		if err != nil {
			s.respondWithError(w, 400, "error decoding body", err)
			return
		}

		saved, err := s.services.UpdateApplicationStatus(
			userID,
			statusID,
			body.Status,
			body.Kind,
		)
		if err != nil {
			s.respondWithError(w, 400, "error updating application status", err)
			return
		}
		response := apiResponse{Ok: true, Data: applicationStatusRequest{
			ID:        saved.ID,
			Kind:      saved.Kind.String(),
			Status:    saved.Status,
			UpdatedAt: saved.UpdatedAt,
			CreatedAt: saved.CreatedAt,
			UserID:    saved.UserID,
		}}
		s.respondWithJSON(w, 200, response)
	}
}

// @Summary Get application status
// @Description Returns a single application status. Users can read their own; admins can read any.
// @Tags application_statuses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param statusID path int true "Application status ID"
// @Success 200 {object} apiResponse{data=applicationStatusRequest}
// @Failure 400 {object} apiResponse
// @Failure 403 {object} apiResponse
// @Router /application_statuses/{statusID} [get]
func (s *Server) handleGetAnApplicationStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := r.PathValue("statusID")
		conv, err := strconv.Atoi(param)
		if err != nil {
			s.respondWithError(w, 400, "invalid id format", err)
			return
		}
		userID, _ := userIDFromContext(r.Context())

		record, err := s.services.GetApplicationStatus(uint(conv))
		if err != nil {
			s.respondWithError(w, 400, "error getting application status", err)
			return
		}
		isAdmin := s.services.IsUserAdmin(userID)
		if !isAdmin && userID != record.UserID {
			s.respondWithError(
				w,
				403,
				"permission denied",
				fmt.Errorf(
					"%s tried to view an application status with id of %d from %s",
					userID,
					conv,
					record.UserID.String(),
				),
			)
			return
		}
		response := apiResponse{Ok: true, Data: applicationStatusRequest{
			ID:        record.ID,
			Kind:      record.Kind.String(),
			Status:    record.Status,
			UpdatedAt: record.UpdatedAt,
			CreatedAt: record.CreatedAt,
			UserID:    record.UserID,
		}}
		s.respondWithJSON(w, 200, response)
	}
}

// @Summary Create application status
// @Description Creates a new application status for the authenticated user.
// @Tags application_statuses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body applicationStatusRequest true "Application status payload"
// @Success 201 {object} apiResponse{data=applicationStatusRequest}
// @Failure 400 {object} apiResponse
// @Router /application_statuses [post]
func (s *Server) handleCreateApplicationStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := userIDFromContext(r.Context())
		decoder := json.NewDecoder(r.Body)

		body := applicationStatusRequest{}
		err := decoder.Decode(&body)
		if err != nil {
			s.respondWithError(w, 400, "error decoding body", err)
			return
		}

		convKind := database.ApplicationStatusKind(body.Kind)
		if !convKind.Valid() {
			s.respondWithError(w, 400, "invalid kind format", nil)
			return
		}

		saved, err := s.services.CreateApplicationStatus(
			userID,
			body.Status,
			convKind,
		)
		if err != nil {
			s.respondWithError(w, 400, "error creating application status", err)
			return
		}

		response := apiResponse{
			Ok: true,
			Data: applicationStatusRequest{
				ID:        saved.ID,
				Kind:      saved.Kind.String(),
				Status:    saved.Status,
				UpdatedAt: saved.UpdatedAt,
				CreatedAt: saved.CreatedAt,
				UserID:    saved.UserID,
			},
			Message: "Created new application status",
		}
		s.respondWithJSON(w, 201, response)
	}
}

// @Summary Delete application status
// @Description Soft-deletes an application status. Only the owner or an admin can delete it.
// @Tags application_statuses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param statusID path int true "Application status ID"
// @Success 200 {object} apiResponse
// @Failure 400 {object} apiResponse
// @Failure 403 {object} apiResponse
// @Router /application_statuses/{statusID} [delete]
func (s *Server) handleDeleteApplicationStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := r.PathValue("statusID")
		userID, _ := userIDFromContext(r.Context())
		conv, err := strconv.Atoi(param)
		if err != nil {
			s.respondWithError(w, 400, "invalid id format", err)
			return
		}
		record, err := s.services.GetApplicationStatus(uint(conv))
		if err != nil {
			s.respondWithError(w, 400, "error getting application status", err)
			return
		}
		isAdmin := s.services.IsUserAdmin(userID)
		if !isAdmin && userID != record.UserID {
			s.respondWithError(w, 403, "permission denied", nil)
			return
		}
		err = s.services.DeleteApplicationStatus(userID, record.ID)
		if err != nil {
			s.respondWithError(w, 400, "error deleting application status", err)
			return
		}
		response := apiResponse{Ok: true, Message: "application status deleted"}
		s.respondWithJSON(w, 200, response)
	}
}
