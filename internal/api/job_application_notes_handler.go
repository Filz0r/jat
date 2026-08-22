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

type noteData struct {
	ID        uint                     `json:"id" validate:"required"`
	Body      string                   `json:"body" validate:"required"`
	Status    applicationStatusRequest `json:"status" validate:"required"`
	UserID    uuid.UUID                `json:"user_id" validate:"required"`
	JobID     uint                     `json:"job_id" validate:"required"`
	CreatedAt time.Time                `json:"created_at" validate:"required"`
	UpdatedAt time.Time                `json:"updated_at" validate:"required"`
}

type noteCreateRequest struct {
	Body string `json:"body" validate:"required"`
}

func extractRequiredNoteData(r *http.Request, extractNote bool) (uint, uint, uuid.UUID, error) {
	jobIDParam := r.PathValue("jobID")
	jobIDConverted, err := strconv.ParseUint(jobIDParam, 10, 32)
	if err != nil {
		return 0, 0, uuid.Nil, err
	}
	jobID := uint(jobIDConverted)
	noteIDParam := r.PathValue("noteID")
	noteIDConverted, err := strconv.ParseUint(noteIDParam, 10, 32)
	if err != nil && extractNote {
		return 0, 0, uuid.Nil, err
	}
	noteID := uint(noteIDConverted)
	userID, _ := userIDFromContext(r.Context())
	return jobID, noteID, userID, nil
}

func convertNoteData(data database.ApplicationNote) noteData {
	return noteData{
		ID:        data.ID,
		Body:      data.Body,
		Status:    createApplicationStatusRequest(data.Status),
		UserID:    data.UserID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		JobID:     data.ApplicationID,
	}
}

// @Summary List notes for a job application
// @Description Returns all notes belonging to a specific job application.
// @Tags job_application_notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param jobID path int true "Job application ID"
// @Success 200 {object} apiResponse{data=[]noteData}
// @Failure 400 {object} apiResponse
// @Failure 404 {object} apiResponse
// @Router /jobs/{jobID}/notes [get]
func (s *Server) handleGetJobNotes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID, _, userID, err := extractRequiredNoteData(r, false)
		if err != nil {
			s.respondWithError(w, 400, "Error parsing Job Application Request", err)
			return
		}
		dbNotes, err := s.services.GetJobApplicationNotesByID(jobID, userID)
		if err != nil {
			s.respondWithError(w, 404, "Note not found", err)
			return
		}
		res := make([]noteData, 0, len(dbNotes))
		for _, note := range dbNotes {
			temp := convertNoteData(note)
			res = append(res, temp)
		}
		s.respondWithJSON(w, 200, apiResponse{
			Data:    res,
			Ok:      true,
			Message: fmt.Sprintf("Found %d job application notes", len(dbNotes)),
		})
	}
}

// @Summary Create job application note
// @Description Adds a note to a job application.
// @Tags job_application_notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param jobID path int true "Job application ID"
// @Param request body noteCreateRequest true "Note body"
// @Success 201 {object} apiResponse{data=noteData}
// @Failure 400 {object} apiResponse
// @Router /jobs/{jobID}/notes [post]
func (s *Server) handleCreateJobNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobIDParam := r.PathValue("jobID")
		jobIDConverted, err := strconv.ParseUint(jobIDParam, 10, 32)
		if err != nil {
			s.respondWithError(w, 400, "Error parsing Job Application ID", err)
			return
		}
		userID, _ := userIDFromContext(r.Context())
		jobID := uint(jobIDConverted)
		decoder := json.NewDecoder(r.Body)
		body := &noteCreateRequest{}
		if err := decoder.Decode(body); err != nil {
			s.respondWithError(w, http.StatusBadRequest, "Error parsing body", err)
			return
		}
		job, err := s.services.GetJobApplicationById(jobID)
		if err != nil {
			s.respondWithError(w, 400, "error processing data", err)
			return
		}
		data, err := s.services.CreateJobNote(jobID, job.StatusID, body.Body, userID)
		if err != nil {
			s.respondWithError(w, 400, "Error creating new job application note", err)
			return
		}
		s.respondWithJSON(w, 201, apiResponse{
			Data:    convertNoteData(data),
			Ok:      true,
			Message: "Created new job application note",
		})
	}
}

// @Summary Update job application note
// @Description Updates an existing note on a job application.
// @Tags job_application_notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param jobID path int true "Job application ID"
// @Param noteID path int true "Note ID"
// @Param request body noteCreateRequest true "Note body"
// @Success 200 {object} apiResponse{data=noteData}
// @Failure 400 {object} apiResponse
// @Router /jobs/{jobID}/notes/{noteID} [put]
func (s *Server) handleUpdateJobNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID, noteID, userID, err := extractRequiredNoteData(r, true)
		if err != nil {
			s.respondWithError(w, 400, "Error parsing Job Application Request", err)
			return
		}
		decoder := json.NewDecoder(r.Body)
		body := &noteCreateRequest{}
		if err := decoder.Decode(body); err != nil {
			s.respondWithError(w, http.StatusBadRequest, "Error parsing body", err)
			return
		}
		if body.Body == "" {
			s.respondWithError(w, http.StatusBadRequest, "Error parsing body", err)
			return
		}
		updated, err := s.services.UpdateJobNoteByID(
			noteID,
			jobID,
			userID,
			body.Body,
		)
		if err != nil {
			s.respondWithError(w, 400, "Error updating job application note", err)
			return
		}
		s.respondWithJSON(w, 200, apiResponse{
			Data:    convertNoteData(updated),
			Ok:      true,
			Message: "Updated job application note",
		})
	}
}

// @Summary Delete job application note
// @Description Deletes a note from a job application.
// @Tags job_application_notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param jobID path int true "Job application ID"
// @Param noteID path int true "Note ID"
// @Success 200 {object} apiResponse
// @Failure 400 {object} apiResponse
// @Router /jobs/{jobID}/notes/{noteID} [delete]
func (s *Server) handleDeleteJobNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID, noteID, userID, err := extractRequiredNoteData(r, true)
		if err != nil {
			s.respondWithError(w, 400, "Error parsing Job Application Request", err)
			return
		}
		err = s.services.DeleteApplicationNoteByID(nil, noteID, jobID, userID)
		if err != nil {
			s.respondWithError(w, 400, "Error deleting note", err)
			return
		}
		s.respondWithJSON(w, 200, apiResponse{
			Ok:      true,
			Message: "Deleted job application note",
		})
	}
}

// @Summary Get job application note
// @Description Returns a single note from a job application.
// @Tags job_application_notes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param jobID path int true "Job application ID"
// @Param noteID path int true "Note ID"
// @Success 200 {object} apiResponse{data=noteData}
// @Failure 400 {object} apiResponse
// @Failure 404 {object} apiResponse
// @Router /jobs/{jobID}/notes/{noteID} [get]
func (s *Server) handleGetJobNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID, noteID, userID, err := extractRequiredNoteData(r, true)
		if err != nil {
			s.respondWithError(w, 400, "Error parsing Job Application Request", err)
			return
		}
		data, err := s.services.GetJobApplicationNoteByID(noteID, jobID, userID)
		if err != nil {
			s.respondWithError(w, 404, "Job Application note not found", err)
			return
		}
		s.respondWithJSON(w, 200, apiResponse{
			Data:    convertNoteData(data),
			Ok:      true,
			Message: "Application note found",
		})
	}
}
