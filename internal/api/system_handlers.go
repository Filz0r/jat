package api

import (
	"net/http"
)

type systemResponse struct {
	Initialized bool `json:"initialized"`
}

func (s *Server) healthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.db == nil {
			s.respondWithError(w, 500, "db not initialized", nil)
			return
		}
		s.respondWithJSON(w, 200, apiResponse{
			Message: "server is up and running",
			Ok:      true,
		})
	})
}

func (s *Server) InitializedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		initialized := s.services.IsInitialized()
		response := apiResponse{
			Data: systemResponse{
				Initialized: initialized,
			},
			Ok: true,
		}
		s.respondWithJSON(w, 200, response)
	})
}

func (s *Server) handleSetInitialized() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.services.IsInitialized() {
			s.respondWithError(w, 400, "service is already initialized", nil)
			return
		}

		isAdmin, _ := userAdminFromContext(r.Context())
		if !isAdmin {
			s.respondWithError(w, 401, "unauthorized", nil)
			return
		}
		userID, _ := userIDFromContext(r.Context())
		firstAdmin := s.services.GetFirstAdmin()
		if userID != firstAdmin {
			s.respondWithError(w, 401, "unauthorized", nil)
			return
		}
		res := s.services.SetInitialized()
		if !res {
			s.respondWithError(w, 400, "error setting service as initialized", nil)
			return
		}
		response := apiResponse{
			Ok:      true,
			Message: "Service is now initialized, good luck hunting for jobs!",
		}
		s.respondWithJSON(w, 201, response)
	})
}
