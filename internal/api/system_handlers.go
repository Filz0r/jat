package api

import "net/http"

type healthResponse struct {
}

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
		s.respondWithJSON(w, 200, systemResponse{
			Status: "server is up and running",
		})
	})
}

func (s *Server) InitializedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		initialized := s.cfg.IsInitialized()
		s.respondWithJSON(w, 200, systemResponse{
			Initialized: initialized,
		})
	})
}
