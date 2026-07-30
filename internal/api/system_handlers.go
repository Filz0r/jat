package api

import "net/http"

type systemResponse struct {
	Status      string `json:"status,omitempty"`
	Initialized bool   `json:"initialized,omitempty"`
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
		initialized := s.cfg.IsInitialized()
		s.respondWithJSON(w, 200, apiResponse{
			Data: systemResponse{Initialized: initialized},
			Ok:   true,
		})
	})
}
