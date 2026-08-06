package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/filz0r/jat/internal/database"
	"github.com/google/uuid"
)

type noteData struct {
	ID        uint      `json:"id"`
	Body      string    `json:"body"`
	Status    string    `json:"status"`
	UserID    uuid.UUID `json:"user_id"`
	JobID     uint      `json:"job_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type noteCreateRequest struct {
	Body string `json:"body"`
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
		Status:    data.Status.Status,
		UserID:    data.UserID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

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
			Message: "Found job application notes",
		})
	}
}

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

func (s *Server) handleDeleteJobNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID, noteID, userID, err := extractRequiredNoteData(r, true)
		if err != nil {
			s.respondWithError(w, 400, "Error parsing Job Application Request", err)
			return
		}
		err = s.services.DeleteApplicationNoteByID(noteID, jobID, userID)
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
