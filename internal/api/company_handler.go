package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/filz0r/jat/internal/database"
	"github.com/google/uuid"
)

type companyBodyRequest struct {
	Name    string `json:"name"`
	Website string `json:"website,omitempty"`
}

type companyResponse struct {
	ID        uint      `json:"id,omitempty"`
	Name      string    `json:"name"`
	Website   string    `json:"website,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedBy uuid.UUID `json:"createdBy"`
	UpdatedBy uuid.UUID `json:"updatedBy"`
}

//TODO: Fix 500 response codes when possible

func (s *Server) handleCreateCompany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromContext(r.Context())
		if !ok {
			s.respondWithError(w, 403, "forbidden", nil)
		}
		decoder := json.NewDecoder(r.Body)
		body := &companyBodyRequest{}
		if err := decoder.Decode(body); err != nil {
			s.respondWithError(w, 400, "error parsing body", err)
			return
		}
		companyRow := database.Company{
			Name: body.Name,
		}
		if body.Website != "" {
			companyRow.Website = &body.Website
		}
		company, err := s.services.CreateCompany(userID, companyRow)
		if err != nil {
			s.respondWithError(w, 400, "error creating company", err)
			return
		}
		response := companyResponse{
			ID:        company.ID,
			Name:      company.Name,
			CreatedAt: company.CreatedAt,
			UpdatedAt: company.UpdatedAt,
			CreatedBy: company.CreatedBy,
			UpdatedBy: company.EditedBy,
		}
		if company.Website != nil {
			response.Website = *company.Website
		}
		s.respondWithJSON(w, 201, successResponse{
			Ok:      true,
			Data:    response,
			Message: "company created",
		})
	}
}

func (s *Server) handleGetAllCompanies() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		companies, err := s.services.GetAllCompanies()
		if err != nil {
			s.respondWithError(w, 400, "error getting all companies", err)
			return
		}
		response := make([]companyResponse, 0, len(companies))
		for _, company := range companies {
			temp := companyResponse{
				Name:      company.Name,
				CreatedAt: company.CreatedAt,
				UpdatedAt: company.UpdatedAt,
				CreatedBy: company.CreatedBy,
				UpdatedBy: company.EditedBy,
			}
			if company.Website != nil {
				temp.Website = *company.Website
			}
			response = append(response, temp)
		}

		s.respondWithJSON(w, 200, successResponse{
			Ok:   true,
			Data: response,
		})
	}
}

func (s *Server) handleUpdateACompany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromContext(r.Context())
		if !ok {
			s.respondWithError(w, 403, "forbidden", nil)
			return
		}
		id := r.PathValue("companyID")
		convertedID, err := strconv.Atoi(id)
		if err != nil {
			s.respondWithError(w, 400, "impossible to convert param", err)
		}
		companyID := uint(convertedID)
		decoder := json.NewDecoder(r.Body)
		body := &companyBodyRequest{}
		if err := decoder.Decode(body); err != nil {
			s.respondWithError(w, 400, "error parsing body", err)
			return
		}
		companyToUpdate := database.Company{
			Name: body.Name,
		}
		companyToUpdate.ID = companyID
		companyToUpdate.UpdatedAt = time.Now()
		if body.Website != "" {
			companyToUpdate.Website = &body.Website
		}
		company, err := s.services.UpdateCompany(userID, companyToUpdate)
		if err != nil {
			s.respondWithError(w, 400, "error updating company", err)
			return
		}
		response := companyResponse{
			ID:        company.ID,
			Name:      company.Name,
			CreatedAt: company.CreatedAt,
			UpdatedAt: company.UpdatedAt,
			CreatedBy: company.CreatedBy,
			UpdatedBy: company.EditedBy,
		}
		if company.Website != nil {
			response.Website = *company.Website
		}
		s.respondWithJSON(w, 200, successResponse{
			Ok:      true,
			Data:    response,
			Message: "company updated",
		})
	}
}

func (s *Server) handleGetACompany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("companyID")
		companyID, err := strconv.Atoi(id)
		if err != nil {
			s.respondWithError(w, 400, "impossible to convert param", err)
		}
		data, err := s.services.GetCompanyByID(companyID)
		if err != nil {
			s.respondWithError(w, 404, "company not found", err)
		}
		response := companyResponse{
			ID:        data.ID,
			Name:      data.Name,
			CreatedAt: data.CreatedAt,
			UpdatedAt: data.UpdatedAt,
			CreatedBy: data.CreatedBy,
			UpdatedBy: data.EditedBy,
		}
		if data.Website != nil {
			response.Website = *data.Website
		}
		s.respondWithJSON(w, 200, successResponse{
			Ok:      true,
			Data:    response,
			Message: "company found",
		})
	}
}

func (s *Server) handleDeleteACompany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromContext(r.Context())
		if !ok {
			s.respondWithError(w, 403, "forbidden", nil)
			return
		}
		id := r.PathValue("companyID")
		companyID, err := strconv.Atoi(id)
		if err != nil {
			s.respondWithError(w, 400, "impossible to convert param", err)
		}
		err = s.services.DeleteCompanyByID(uint(companyID), userID)
		if err != nil {
			s.respondWithError(w, 404, "company not found", err)
		}

		s.respondWithJSON(w, 200, successResponse{
			Ok:      true,
			Message: "company deleted",
		})
	}
}
