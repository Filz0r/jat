package api

import (
	"errors"
	"net/http"

	"github.com/filz0r/jat/internal/config"
	"github.com/filz0r/jat/internal/services"
)

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.middlewareLogRequest(s.middlewareClientType(s.mux)),
	}
	s.loadRoutes()
	s.cfg = config.New()
	err := s.cfg.LoadFromEnv()
	if err != nil {
		return err
	}
	s.logger.Printf("API server started on port %s", s.port)
	err = s.server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) loadRoutes() {
	// system handlers
	s.mux.Handle("GET /api/health", s.healthHandler())
	s.mux.Handle("GET /api/initialized", s.InitializedHandler())

	//auth handlers
	s.mux.Handle("POST /api/auth/login", s.handleUserLogin())
	s.mux.Handle("POST /api/auth/logout",
		s.middlewareAuth(s.middlewareRefreshToken(s.handleUserLogout())))
	s.mux.Handle("GET /api/auth/refresh_token",
		s.middlewareRefreshToken(s.handleUserTokenRefresh()))
	s.mux.Handle("GET /api/auth/revoke_token",
		s.middlewareRefreshToken(s.handleUserRevokeToken()))

	// user handlers
	s.mux.Handle("POST /api/users", s.handleUserCreate())
	s.mux.Handle("GET /api/users/{userID}",
		s.middlewareAuth(s.handleGetSingleUser()))
	s.mux.Handle("PUT /api/users", s.middlewareAuth(s.handleUserUpdate()))

	// Company handlers
	s.mux.Handle("POST /api/company", s.middlewareAuth(s.handleCreateCompany()))
	s.mux.Handle("GET /api/company", s.middlewareAuth(s.handleGetAllCompanies()))
	s.mux.Handle("GET /api/company/{companyID}", s.middlewareAuth(s.handleGetACompany()))
	s.mux.Handle("PUT /api/company/{companyID}", s.middlewareAuth(s.handleUpdateACompany()))
	s.mux.Handle("DELETE /api/company/{companyID}", s.middlewareAdminUser(s.handleDeleteACompany()))

	// Application Statuses
	s.mux.Handle("GET /api/application_statuses", s.middlewareAuth(s.handleGetUserApplicationStatus()))
	s.mux.Handle("GET /api/application_statuses/{statusID}", s.middlewareAuth(s.handleGetAnApplicationStatus()))
	s.mux.Handle("PUT /api/application_statuses/{statusID}", s.middlewareAuth(s.handleUpdateApplicationStatus()))
	s.mux.Handle("DELETE /api/application_statuses", s.middlewareAuth(s.handleDeleteApplicationStatus()))
	s.mux.Handle("POST /api/application_statuses", s.middlewareAuth(s.handleCreateApplicationStatus()))

	// admin handlers
	s.mux.Handle("GET /api/admin/users", s.middlewareAdminUser(s.handleGetAllUsers()))
	s.mux.Handle("GET /api/admin/users/{userID}", s.middlewareAdminUser(s.handleMakeUserAdmin(true)))
	s.mux.Handle("DELETE /api/admin/users/{userID}", s.middlewareAdminUser(s.handleMakeUserAdmin(false)))
	// TODO: add a restore company change endpoint for admins
	// TODO: add admin endpoints to get soft deleted application status and a way to restore them
}

func New(cfg *config.ConfigFile) (*Server, error) {
	var servePort string
	if cfg == nil {
		return nil, errors.New("config file isn't initialized")
	}
	servePort, err := cfg.GetPort()
	if err != nil {
		return nil, err
	}

	server := &Server{
		db:       cfg.GetDB(),
		mux:      http.NewServeMux(),
		port:     servePort,
		services: services.NewServiceManager(cfg.GetDB()),
		logger:   newServerLogger(),
	}
	return server, nil
}
