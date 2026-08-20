package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type companyBodyRequest struct {
	Name    string `json:"name" validate:"required"`
	Website string `json:"website,omitempty"`
}

type companyResponse struct {
	ID        uint      `json:"id,omitempty" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Website   string    `json:"website,omitempty"`
	CreatedAt time.Time `json:"createdAt" validate:"required"`
	UpdatedAt time.Time `json:"updatedAt" validate:"required"`
	CreatedBy uuid.UUID `json:"createdBy" validate:"required"`
	UpdatedBy uuid.UUID `json:"updatedBy" validate:"required"`
}

//TODO: Fix 500 response codes when possible

// @Summary Create company
// @Description Creates a new company.
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body companyBodyRequest true "Company payload"
// @Success 201 {object} apiResponse{data=companyResponse}
// @Failure 400 {object} apiResponse
// @Failure 403 {object} apiResponse
// @Failure 409 {object} apiResponse
// @Router /company [post]
func (s *Server) handleCreateCompany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := userIDFromContext(r.Context())
		if !ok {
			s.respondWithError(w, 403, "forbidden", nil)
			return
		}
		decoder := json.NewDecoder(r.Body)
		body := &companyBodyRequest{}
		if err := decoder.Decode(body); err != nil {
			s.respondWithError(w, 400, "error parsing body", err)
			return
		}

		_, err := s.services.FindCompanyByName(body.Name)
		if err == nil {
			s.respondWithError(w, 409, "this company already exists (capitalization issue)", err)
			return
		}

		company, err := s.services.CreateCompany(userID, body.Name, body.Website)
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
		s.respondWithJSON(w, 201, apiResponse{
			Ok:      true,
			Data:    response,
			Message: "company created",
		})
	}
}

// @Summary List all companies
// @Description Returns every company in the system.
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} apiResponse{data=[]companyResponse}
// @Failure 400 {object} apiResponse
// @Router /company [get]
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

		s.respondWithJSON(w, 200, apiResponse{
			Ok:   true,
			Data: response,
		})
	}
}

// @Summary Update company
// @Description Updates an existing company.
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param companyID path int true "Company ID"
// @Param request body companyBodyRequest true "Company payload"
// @Success 200 {object} apiResponse{data=companyResponse}
// @Failure 400 {object} apiResponse
// @Failure 403 {object} apiResponse
// @Router /company/{companyID} [put]
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
			return
		}
		companyID := uint(convertedID)
		decoder := json.NewDecoder(r.Body)
		body := &companyBodyRequest{}
		if err := decoder.Decode(body); err != nil {
			s.respondWithError(w, 400, "error parsing body", err)
			return
		}
		company, err := s.services.UpdateCompany(companyID, userID, body.Name, body.Website)
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
		s.respondWithJSON(w, 200, apiResponse{
			Ok:      true,
			Data:    response,
			Message: "company updated",
		})
	}
}

// @Summary Get company
// @Description Returns a single company by ID.
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param companyID path int true "Company ID"
// @Success 200 {object} apiResponse{data=companyResponse}
// @Failure 400 {object} apiResponse
// @Failure 404 {object} apiResponse
// @Router /company/{companyID} [get]
func (s *Server) handleGetACompany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("companyID")
		companyParsed, err := strconv.ParseUint(id, 10, 32)
		companyID := uint(companyParsed)
		if err != nil {
			s.respondWithError(w, 400, "impossible to convert param", err)
			return
		}
		data, err := s.services.GetCompanyByID(companyID)
		if err != nil {
			s.respondWithError(w, 404, "company not found", err)
			return
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
		s.respondWithJSON(w, 200, apiResponse{
			Ok:      true,
			Data:    response,
			Message: "company found",
		})
	}
}

// @Summary Delete company
// @Description Admin-only. Soft-deletes a company by ID.
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param companyID path int true "Company ID"
// @Success 200 {object} apiResponse
// @Failure 400 {object} apiResponse
// @Failure 403 {object} apiResponse
// @Failure 404 {object} apiResponse
// @Router /company/{companyID} [delete]
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
			return
		}
		err = s.services.DeleteCompanyByID(uint(companyID), userID)
		if err != nil {
			s.respondWithError(w, 404, "company not found", err)
			return
		}

		s.respondWithJSON(w, 200, apiResponse{
			Ok:      true,
			Message: "company deleted",
		})
	}
}
