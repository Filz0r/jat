package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/filz0r/jat/internal/database"
)

type companyBodyRequest struct {
	Name    string `json:"name" validate:"required"`
	Website string `json:"website,omitempty"`
}

type companyResponse struct {
	ID            uint                `json:"id,omitempty" validate:"required"`
	Name          string              `json:"name" validate:"required"`
	Website       string              `json:"website,omitempty"`
	CreatedAt     time.Time           `json:"created_at" validate:"required"`
	UpdatedAt     time.Time           `json:"updated_at" validate:"required"`
	UserCount     int64               `json:"user_count,omitempty"`
	TotalCount    int64               `json:"total_count,omitempty"`
	CreatedByUser *userCreateResponse `json:"created_by_user,omitempty"`
	UpdatedByUser *userCreateResponse `json:"updated_by_user,omitempty"`
}

type companyListQuery struct {
	UserCount    bool `query:"user_count"`
	TotalCount   bool `query:"total_count"`
	PreloadUsers bool `query:"preload_users"`
}

type countCompanyQuery struct {
	TotalCount bool `query:"total_count"`
}

type companyChangeHistoryResponse struct {
	ID         uint               `json:"id" validate:"required"`
	CompanyID  uint               `json:"company_id" validate:"required"`
	NewName    string             `json:"new_name,omitempty"`
	OldWebsite string             `json:"old_website,omitempty"`
	NewWebsite string             `json:"new_website,omitempty"`
	OldName    string             `json:"old_name,omitempty"`
	CreatedAt  time.Time          `json:"created_at" validate:"required"`
	ChangedBy  userCreateResponse `json:"changed_by" validate:"required"`
	Reverted   bool               `json:"reverted" validate:"required"`
}

//TODO: Fix 500 response codes when possible

func createCompanyResponse(data database.Company, totalCount, userCount int64, includeCounts bool, cbUser, ebUser *database.User) companyResponse {
	result := companyResponse{
		ID:        data.ID,
		Name:      data.Name,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
	if data.Website != nil {
		result.Website = *data.Website
	}
	if includeCounts {
		result.TotalCount = totalCount
		result.UserCount = userCount
	}
	if cbUser != nil {
		userResponse := createUserResponse(*cbUser)
		result.CreatedByUser = &userResponse
	}
	if ebUser != nil {
		userResponse := createUserResponse(*ebUser)
		result.UpdatedByUser = &userResponse
	}
	return result
}

func createCompanyChangeHistoryResponse(data database.CompanyChangeHistory) companyChangeHistoryResponse {
	result := companyChangeHistoryResponse{
		ID:        data.ID,
		CompanyID: data.CompanyID,
		CreatedAt: data.CreatedAt,
		Reverted:  data.Reverted,
		ChangedBy: createUserResponse(data.ChangedByUser),
	}
	if data.OldNameValue != nil {
		result.OldName = *data.OldNameValue
	}
	if data.NewNameValue != nil {
		result.NewName = *data.NewNameValue
	}
	if data.OldWebsiteValue != nil {
		result.OldWebsite = *data.OldWebsiteValue
	}
	if data.NewWebsiteValue != nil {
		result.NewWebsite = *data.NewWebsiteValue
	}
	return result
}

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
		userID, _ := userIDFromContext(r.Context())
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

		response := createCompanyResponse(company, 0, 0, false, nil, nil)
		s.respondWithJSON(w, 201, apiResponse{
			Ok:      true,
			Data:    response,
			Message: "Company created",
		})
	}
}

// @Summary List all companies
// @Description Returns every company in the system. (returns 403 if a non admin user passes the preload_users query param)
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_count query bool false "Include a count of user applications for this company"
// @Param preload_users query bool false "Include user data for the creation/edits of companies (admin only)"
// @Param total_count query bool false "Include a count of total applications (global) for this company"
// @Success 200 {object} apiResponse{data=[]companyResponse}
// @Failure 403 {object} apiResponse
// @Failure 400 {object} apiResponse
// @Router /company [get]
func (s *Server) handleGetAllCompanies() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := userIDFromContext(r.Context())
		params, err := BindQuery[companyListQuery](r)
		if err != nil {
			s.respondWithError(w, 400, "invalid query", err)
			return
		}
		if params.PreloadUsers && !s.services.IsUserAdmin(userID) {
			s.respondWithError(w, 403, "forbidden", nil)
			return
		}
		companies, counts, err := s.services.GetAllCompanies(
			nil,
			userID,
			params.PreloadUsers,
			params.TotalCount,
			params.UserCount,
		)
		if err != nil {
			s.respondWithError(w, 400, "error getting all companies", err)
			return
		}
		response := make([]companyResponse, 0, len(companies))
		for _, company := range companies {
			includeCounts := params.UserCount || params.TotalCount
			var cbUser, ebUser *database.User
			if params.PreloadUsers {
				cbUser = &company.CreatedByUser
				ebUser = &company.CreatedByUser
			} else {
				cbUser = nil
				ebUser = nil
			}
			temp := createCompanyResponse(
				company,
				counts[company.ID].TotalCount,
				counts[company.ID].UserCount,
				includeCounts,
				cbUser,
				ebUser,
			)
			response = append(response, temp)
		}

		s.respondWithJSON(w, 200, apiResponse{
			Ok:      true,
			Data:    response,
			Message: fmt.Sprintf("Found %d companies", len(response)),
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
		userID, _ := userIDFromContext(r.Context())
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
		company, err := s.services.UpdateCompany(nil, companyID, userID, body.Name, body.Website)
		if err != nil {
			s.respondWithError(w, 400, "error updating company", err)
			return
		}
		response := createCompanyResponse(company, 0, 0, false, nil, nil)
		s.respondWithJSON(w, 200, apiResponse{
			Ok:      true,
			Data:    response,
			Message: "Company updated",
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
// @Param user_count query bool false "Include a count of user applications for this company"
// @Param total_count query bool false "Include a count of total applications (global) for this company"
// @Success 200 {object} apiResponse{data=companyResponse}
// @Failure 400 {object} apiResponse
// @Failure 404 {object} apiResponse
// @Router /company/{companyID} [get]
func (s *Server) handleGetACompany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("companyID")
		companyParsed, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			s.respondWithError(w, 400, "impossible to convert param", err)
			return
		}

		params, err := BindQuery[companyListQuery](r)
		if err != nil {
			s.respondWithError(w, 400, "invalid query", err)
			return
		}
		userID, _ := userIDFromContext(r.Context())
		data, counts, err := s.services.GetCompanyByID(nil, uint(companyParsed), userID, params.TotalCount, params.UserCount)
		if err != nil {
			s.respondWithError(w, 404, "company not found", err)
			return
		}
		includeCounts := params.UserCount || params.TotalCount
		response := createCompanyResponse(data, counts.TotalCount, counts.UserCount, includeCounts, nil, nil)
		if data.Website != nil {
			response.Website = *data.Website
		}
		s.respondWithJSON(w, 200, apiResponse{
			Ok:      true,
			Data:    response,
			Message: "Company found",
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
		userID, _ := userIDFromContext(r.Context())
		id := r.PathValue("companyID")
		companyID, err := strconv.Atoi(id)
		if err != nil {
			s.respondWithError(w, 400, "impossible to convert param", err)
			return
		}
		err = s.services.DeleteCompanyByID(nil, uint(companyID), userID)
		if err != nil {
			s.respondWithError(w, 404, "company not found", err)
			return
		}

		s.respondWithJSON(w, 200, apiResponse{
			Ok:      true,
			Message: "Company deleted",
		})
	}
}

// @Summary Count Company Job Applications
// @Description Counts the Job Applications a company has company ID. If the total_count query param is true it returns the total applications for a company, otherwise it returns the total applications for the user
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param companyID path int true "Company ID"
// @Param total_count query bool false "Include a count of total applications (global) for this company"
// @Success 200 {object} apiResponse{data=apiCountResult}
// @Failure 400 {object} apiResponse
// @Failure 404 {object} apiResponse
// @Router /company/{companyID}/count [get]
func (s *Server) handleCountCompanyApplications() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := userIDFromContext(r.Context())
		rawID := r.PathValue("companyID")
		companyID64, err := strconv.ParseUint(rawID, 10, 32)
		if err != nil {
			s.respondWithError(w, 400, "impossible to convert param", err)
			return
		}
		params, err := BindQuery[countCompanyQuery](r)
		if err != nil {
			s.respondWithError(w, 400, "invalid query", err)
			return
		}
		count, err := s.services.CountCompanyApplications(nil, uint(companyID64), userID, !params.TotalCount)
		if err != nil {
			s.respondWithError(w, 404, "company not found", err)
			return
		}
		response := apiResponse{
			Ok: true,
			Data: apiCountResult{
				Count: count,
			},
			Message: fmt.Sprintf("Found %d company applications", count),
		}
		s.respondWithJSON(w, 200, response)
	}
}

// @Summary Get history of changes to a company
// @Description Returns the history of changes to a company record (admin only)
// @Tags companies
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param companyID path int true "Company ID"
// @Success 200 {object} apiResponse{data=[]companyChangeHistoryResponse}
// @Failure 400 {object} apiResponse
// @Failure 404 {object} apiResponse
// @Router /company/{companyID}/history [get]
func (s *Server) handleGetCompanyHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := userIDFromContext(r.Context())
		id := r.PathValue("companyID")
		companyID64, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			s.respondWithError(w, 400, "impossible to convert param", err)
			return
		}
		history, err := s.services.GetCompanyHistory(nil, userID, uint(companyID64))
		if err != nil {
			s.respondWithError(w, 404, "company not found", err)
			return
		}
		responseData := make([]companyChangeHistoryResponse, 0, len(history))
		for _, record := range history {
			temp := createCompanyChangeHistoryResponse(record)
			responseData = append(responseData, temp)
		}
		s.respondWithJSON(w, 200, apiResponse{
			Ok:      true,
			Data:    responseData,
			Message: fmt.Sprintf("Found %d company history", len(history)),
		})
	}
}

// @Summary Reverts a company change record
// @Description Reverts a change that was made to a company record (admin only)
// @Tags companies
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param companyID path int true "Company ID"
// @Param changeID path int true "Change ID"
// @Success 200 {object} apiResponse
// @Failure 400 {object} apiResponse
// @Failure 404 {object} apiResponse
// @Router /company/{companyID}/history/{changeID} [put]
func (s *Server) handleRestoreACompanyChange() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := userIDFromContext(r.Context())
		id := r.PathValue("companyID")
		companyID64, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			s.respondWithError(w, 400, "impossible to convert param", err)
			return
		}
		change := r.PathValue("changeID")
		changeID64, err := strconv.ParseUint(change, 10, 32)
		if err != nil {
			s.respondWithError(w, 400, "impossible to convert param", err)
			return
		}
		err = s.services.RevertCompanyChange(nil, userID, uint(companyID64), uint(changeID64))
		if err != nil {
			s.respondWithError(w, 404, "company not found", err)
			return
		}
		s.respondWithJSON(w, 200, apiResponse{
			Ok:      true,
			Message: "Company change restored",
		})
	}
}
