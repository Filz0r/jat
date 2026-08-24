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
	Status string `json:"status" validate:"required"`
	Kind   string `json:"kind" validate:"required"`
}

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

type applicationStatusListQuery struct {
	IncludeArchived bool `query:"include_archived"`
}

func createApplicationHistoryRequest(d database.StatusHistory) applicationStatusHistoryResponse {
	response := applicationStatusHistoryResponse{
		ID:            d.ID,
		ApplicationID: d.ApplicationID,
		NewStatus:     createApplicationStatusRequest(d.NewStatus),
		CreatedAt:     d.CreatedAt,
	}

	if d.OldStatus != nil {
		old := createApplicationStatusRequest(*d.OldStatus)
		response.OldStatus = &old
	}

	return response
}

func createApplicationStatusRequest(d database.ApplicationStatus) applicationStatusResponse {
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

// @Summary List current user's application statuses
// @Description Returns all application statuses belonging to the authenticated user.
// @Tags application_statuses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param include_archived query bool false "Include archived statuses"
// @Success 200 {object} apiResponse{data=[]applicationStatusRequest}
// @Failure 400 {object} apiResponse
// @Failure 401 {object} apiResponse
// @Router /application_statuses [get]
func (s *Server) handleGetUserApplicationStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := userIDFromContext(r.Context())
		params, err := BindQuery[applicationStatusListQuery](r)
		if err != nil {
			s.respondWithError(w, 400, err.Error(), err)
			return
		}
		data, err := s.services.GetAllApplicationStatus(userID, params.IncludeArchived)
		if err != nil {
			s.respondWithError(w, 400, "error getting application status", err)
			return
		}
		converted := make([]applicationStatusResponse, 0, len(data))
		for _, d := range data {
			temp := createApplicationStatusRequest(d)
			converted = append(converted, temp)
		}

		s.respondWithJSON(w, 200, apiResponse{
			Ok:      true,
			Data:    converted,
			Message: fmt.Sprintf("Found %d Application Status", len(converted)),
		})
	}
}

// TODO: add this to the API when admin endpoints are implemented
func (s *Server) handleGetAllApplicationStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := BindQuery[applicationStatusListQuery](r)
		if err != nil {
			s.respondWithError(w, 400, err.Error(), err)
			return
		}
		data, err := s.services.GetAllApplicationStatusAdmin(params.IncludeArchived)
		if err != nil {
			s.respondWithError(w, 400, "error getting application status", err)
			return
		}
		// TODO: Might need to change this to user applicationStatusRequest instead of the db type
		s.respondWithJSON(w, 200, apiResponse{
			Ok:      true,
			Data:    data,
			Message: fmt.Sprintf("Found %d Application Status", len(data)),
		})
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
// @Success 200 {object} apiResponse{data=applicationStatusResponse}
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
		response := apiResponse{
			Ok: true,
			Data: applicationStatusResponse{
				ID:        saved.ID,
				Kind:      saved.Kind.String(),
				Status:    saved.Status,
				UpdatedAt: saved.UpdatedAt,
				CreatedAt: saved.CreatedAt,
				UserID:    saved.UserID,
			},
			Message: "Updated Application Status",
		}
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
// @Success 200 {object} apiResponse{data=applicationStatusResponse}
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
		response := apiResponse{
			Ok:      true,
			Data:    createApplicationStatusRequest(record),
			Message: fmt.Sprintf("Found Application Status with ID %d", record.ID),
		}
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
// @Success 201 {object} apiResponse{data=applicationStatusResponse}
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
			Ok:      true,
			Data:    createApplicationStatusRequest(saved),
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
		err = s.services.DeleteApplicationStatus(userID, record.ID, false)
		if err != nil {
			s.respondWithError(w, 400, "error deleting application status", err)
			return
		}
		response := apiResponse{
			Ok:      true,
			Message: "Application status deleted",
		}
		s.respondWithJSON(w, 200, response)
	}
}

// @Summary Get job application status history
// @Description Returns the status change history for a single job application. Users can read their own; admins can read any.
// @Tags job_applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param jobID path int true "Job application ID"
// @Success 200 {object} apiResponse{data=[]applicationStatusHistoryResponse}
// @Failure 400 {object} apiResponse
// @Failure 403 {object} apiResponse
// @Router /jobs/{jobID}/history [get]
func (s *Server) handleGetJobApplicationHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawJobId := r.PathValue("jobID")
		jobID64, err := strconv.ParseUint(rawJobId, 10, 32)
		if err != nil {
			s.respondWithError(w, 400, "Error parsing Job Application ID", err)
			return
		}
		jobID := uint(jobID64)
		userID, _ := userIDFromContext(r.Context())
		data, err := s.services.GetJobApplicationStatusChanges(userID, jobID)
		if err != nil {
			s.respondWithError(w, 400, "Could not find job application history data", err)
			return
		}
		result := make([]applicationStatusHistoryResponse, 0, len(data))
		for _, d := range data {
			temp := createApplicationHistoryRequest(d)
			result = append(result, temp)
		}
		s.respondWithJSON(w, 200, apiResponse{
			Ok:      true,
			Message: fmt.Sprintf("Found %d Job Application History changes", len(data)),
			Data:    result,
		})
	}
}
