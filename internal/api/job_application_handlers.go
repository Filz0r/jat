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

type applicationResponse struct {
	ID        uint                      `json:"id" validate:"required"`
	Title     string                    `json:"title" validate:"required"`
	URL       string                    `json:"url" validate:"required"`
	UserID    uuid.UUID                 `json:"user_id" validate:"required"`
	CreatedAt time.Time                 `json:"created_at" validate:"required"`
	UpdatedAt time.Time                 `json:"updated_at" validate:"required"`
	Status    applicationStatusResponse `json:"status" validate:"required"`
	Company   companyResponse           `json:"company" validate:"required"`
}

type applicationRequest struct {
	Title     string    `json:"title" validate:"required"`
	URL       string    `json:"url" validate:"required"`
	StatusID  int       `json:"status_id" validate:"required"`
	CompanyID int       `json:"company_id" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
}

type jobApplicationQueries struct {
	StatusID  []uint `query:"status_id"`
	CompanyID []uint `query:"company_id"`
}

func generateApplicationResponseFromRow(row database.JobApplication) applicationResponse {
	res := applicationResponse{
		ID:        row.ID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
		UserID:    row.UserID,
		Title:     row.Title,
		URL:       row.Url,
		Company: companyResponse{
			ID:        row.Company.ID,
			Name:      row.Company.Name,
			CreatedAt: row.Company.CreatedAt,
			UpdatedAt: row.Company.UpdatedAt,
			CreatedBy: row.Company.CreatedBy,
			UpdatedBy: row.Company.EditedBy,
		},
		Status: applicationStatusResponse{
			ID:        row.StatusID,
			CreatedAt: row.Status.CreatedAt,
			UpdatedAt: row.Status.UpdatedAt,
			Status:    row.Status.Status,
			Kind:      row.Status.Kind.String(),
			UserID:    row.Status.UserID,
		},
	}
	if row.Company.Website != nil {
		res.Company.Website = *row.Company.Website
	}
	return res
}

// @Summary List current user's job applications
// @Description Returns all job applications belonging to the authenticated user.
// @Tags job_applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param company_id query []integer false "Filter by Company IDs"
// @Param status_id  query []integer false "Filter by Status IDs"
// @Success 200 {object} apiResponse{data=[]applicationResponse}
// @Failure 400 {object} apiResponse
// @Failure 401 {object} apiResponse
// @Router /jobs [get]
func (s *Server) handleGetUserJobApplications() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := userIDFromContext(r.Context())
		params, err := BindQuery[jobApplicationQueries](r)
		if err != nil {
			s.respondWithError(w, 400, err.Error(), err)
			return
		}
		data, err := s.services.GetUserJobApplications(userID, params.CompanyID, params.StatusID)
		if err != nil {
			s.respondWithError(w, 400, "Could not get job applications", err)
			return
		}
		res := make([]applicationResponse, 0, len(data))
		for _, job := range data {
			temp := generateApplicationResponseFromRow(job)
			res = append(res, temp)
		}
		response := apiResponse{
			Data:    res,
			Ok:      true,
			Message: fmt.Sprintf("Fetched %d job applications", len(res)),
		}
		s.respondWithJSON(w, 200, response)
	}
}

// @Summary Create job application
// @Description Creates a new job application for the authenticated user.
// @Tags job_applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body applicationRequest true "Job application payload"
// @Success 201 {object} apiResponse{data=applicationResponse}
// @Failure 400 {object} apiResponse
// @Failure 401 {object} apiResponse
// @Router /jobs [post]
func (s *Server) handleCreateJobApplication() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := userIDFromContext(r.Context())
		decoder := json.NewDecoder(r.Body)
		body := applicationRequest{}
		err := decoder.Decode(&body)
		if err != nil {
			s.respondWithError(w, 400, "Malformed request body", err)
			return
		}
		if body.StatusID < 1 || body.CompanyID < 1 || body.Title == "" || body.CreatedAt.IsZero() {
			s.respondWithError(w, 400, "Invalid request body", nil)
			return
		}

		row, err := s.services.CreateUserJobApplication(
			userID,
			uint(body.StatusID),
			uint(body.CompanyID),
			body.Title,
			body.URL,
			body.CreatedAt,
		)
		if err != nil {
			s.respondWithError(w, 400, "Could not create job application", err)
			return
		}
		response := apiResponse{
			Data:    generateApplicationResponseFromRow(row),
			Ok:      true,
			Message: "Created job application",
		}
		s.respondWithJSON(w, 201, response)
	}
}

// @Summary Get job application
// @Description Returns a single job application. Users can read their own; admins can read any.
// @Tags job_applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param jobID path int true "Job application ID"
// @Success 200 {object} apiResponse{data=applicationResponse}
// @Failure 400 {object} apiResponse
// @Failure 403 {object} apiResponse
// @Failure 404 {object} apiResponse
// @Router /jobs/{jobID} [get]
func (s *Server) handleGetJobApplication() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := userIDFromContext(r.Context())
		rawJobID := r.PathValue("jobID")
		jobID64, err := strconv.ParseUint(rawJobID, 10, 32)
		if err != nil {
			s.respondWithError(w, 400, "Invalid job id", err)
			return
		}
		jobID := uint(jobID64)
		job, err := s.services.GetJobApplicationById(jobID)
		if err != nil {
			s.respondWithError(w, 404, "Job not found", err)
			return
		}
		if !s.services.IsUserAdmin(userID) && job.UserID != userID {
			s.respondWithError(w, 403, "You do not have permission to view this job", nil)
			return
		}
		response := apiResponse{
			Data:    generateApplicationResponseFromRow(job),
			Ok:      true,
			Message: "Found job application",
		}
		s.respondWithJSON(w, 200, response)
	}
}

// @Summary Update job application status
// @Description Changes the status of a job application.
// @Tags job_applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param jobID path int true "Job application ID"
// @Param statusID path int true "Application status ID"
// @Success 200 {object} apiResponse{data=applicationResponse}
// @Failure 400 {object} apiResponse
// @Failure 403 {object} apiResponse
// @Failure 404 {object} apiResponse
// @Router /jobs/{jobID}/status/{statusID} [put]
func (s *Server) handleUpdateJobApplicationStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := userIDFromContext(r.Context())
		rawJobId := r.PathValue("jobID")
		rawStatusId := r.PathValue("statusID")
		jobID64, err := strconv.ParseUint(rawJobId, 10, 32)
		if err != nil {
			s.respondWithError(w, 400, "Invalid job id", err)
			return
		}
		statusID, err := strconv.ParseUint(rawStatusId, 10, 32)
		if err != nil {
			s.respondWithError(w, 400, "Invalid status id", err)
			return
		}
		jobApplication, err := s.services.GetJobApplicationById(uint(jobID64))
		if err != nil {
			s.respondWithError(w, 404, "Job not found", err)
			return
		}
		res, err := s.services.UpdateJobApplicationStatus(
			userID,
			jobApplication.ID,
			uint(statusID),
		)
		if err != nil {
			s.respondWithError(w, 400, "Could not update job application", err)
			return
		}
		response := apiResponse{
			Data:    generateApplicationResponseFromRow(res),
			Ok:      true,
			Message: "Updated job application",
		}
		s.respondWithJSON(w, 200, response)
	}
}

// @Summary Delete job application
// @Description Soft-deletes a job application. Only the owner or an admin can delete it.
// @Tags job_applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param jobID path int true "Job application ID"
// @Success 200 {object} apiResponse
// @Failure 400 {object} apiResponse
// @Failure 403 {object} apiResponse
// @Router /jobs/{jobID} [delete]
func (s *Server) handleDeleteJobApplication() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := userIDFromContext(r.Context())
		rawJobID := r.PathValue("jobID")
		jobID64, err := strconv.ParseUint(rawJobID, 10, 32)
		if err != nil {
			s.respondWithError(w, 400, "Invalid job id", err)
			return
		}
		err = s.services.DeleteJobApplication(nil, uint(jobID64), userID)
		if err != nil {
			s.respondWithError(w, 400, err.Error(), err)
			return
		}
		response := apiResponse{
			Ok:      true,
			Message: "Deleted job application",
		}
		s.respondWithJSON(w, 200, response)
	}
}
