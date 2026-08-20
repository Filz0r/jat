package api

import (
	"net/http"
)

type systemResponse struct {
	Initialized bool `json:"initialized"`
}

// @Summary Health check
// @Description Returns a success message if the database is reachable.
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} apiResponse{data=string}
// @Failure 500 {object} apiResponse
// @Router /health [get]
func (s *Server) healthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.db == nil {
			s.respondWithError(w, 500, "db not initialized", nil)
			return
		}
		s.respondWithJSON(w, 200, apiResponse{
			Message: "Server is up and running",
			Ok:      true,
		})
	})
}

// @Summary Check initialization state
// @Description Returns whether the server has been initialized (first admin created). This is used by the web UI to decide between the create-account and login flows.
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} apiResponse{data=systemResponse}
// @Router /initialized [get]
func (s *Server) InitializedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		initialized := s.services.IsInitialized()
		response := apiResponse{
			Data: systemResponse{
				Initialized: initialized,
			},
			Message: "System is initialized",
			Ok:      true,
		}
		s.respondWithJSON(w, 200, response)
	})
}

// @Summary Mark service as initialized
// @Description Sets the server initialized flag. Only the first admin can call this.
// @Tags system
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} apiResponse
// @Failure 400 {object} apiResponse
// @Failure 403 {object} apiResponse
// @Router /initialized/set [get]
func (s *Server) handleSetInitialized() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.services.IsInitialized() {
			s.respondWithError(w, 400, "Service is already initialized", nil)
			return
		}

		isAdmin, _ := userAdminFromContext(r.Context())
		if !isAdmin {
			s.respondWithError(w, 403, "Forbidden Action", nil)
			return
		}
		userID, _ := userIDFromContext(r.Context())
		firstAdmin := s.services.GetFirstAdmin()
		if userID != firstAdmin {
			s.respondWithError(w, 403, "Forbidden Action", nil)
			return
		}
		res := s.services.SetInitialized()
		if !res {
			s.respondWithError(w, 400, "Error setting service as initialized", nil)
			return
		}
		response := apiResponse{
			Ok:      true,
			Message: "Service is now initialized, good luck hunting for jobs!",
		}
		s.respondWithJSON(w, 201, response)
	})
}
