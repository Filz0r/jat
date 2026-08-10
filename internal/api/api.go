package api

import (
	"errors"
	"net/http"
	"os"

	"github.com/filz0r/jat/internal/config"
	"github.com/filz0r/jat/internal/services"
	"github.com/filz0r/jat/internal/version"
)

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.middlewareLogRequest(s.middlewareClientType(s.mux)),
	}
	s.loadRoutes()
	initialConfigsExist, err := s.services.InitialConfigsExist()
	if err != nil {
		s.logger.Fatal(err)
	}
	if !initialConfigsExist {
		s.logger.Println("initial configs not exist, creating them...")
		err := s.services.CreateInitialConfigs()
		if err != nil {
			s.logger.Fatal(err)
		}
	}
	s.services.UpdateCurrentVersion()
	if s.isDebug {
		s.logger.Println(version.Info())
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
	s.mux.Handle("GET /api/initialized/set", s.middlewareAdminUser(s.handleSetInitialized()))

	//auth handlers
	s.mux.Handle("POST /api/auth/login", s.handleUserLogin())
	s.mux.Handle("POST /api/auth/logout",
		s.middlewareAuth(s.middlewareRefreshToken(s.handleUserLogout())))
	s.mux.Handle("GET /api/auth/refresh_token",
		s.middlewareRefreshToken(s.handleUserTokenRefresh()))
	s.mux.Handle("GET /api/auth/revoke_token",
		s.middlewareRefreshToken(s.handleUserRevokeToken()))
	// TODO: add a method to fetch the current sessions

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

	// Application Statuses Handlers
	s.mux.Handle("GET /api/application_statuses", s.middlewareAuth(s.handleGetUserApplicationStatus()))
	s.mux.Handle("GET /api/application_statuses/{statusID}", s.middlewareAuth(s.handleGetAnApplicationStatus()))
	s.mux.Handle("PUT /api/application_statuses/{statusID}", s.middlewareAuth(s.handleUpdateApplicationStatus()))
	s.mux.Handle("DELETE /api/application_statuses/{statusID}", s.middlewareAuth(s.handleDeleteApplicationStatus()))
	s.mux.Handle("POST /api/application_statuses", s.middlewareAuth(s.handleCreateApplicationStatus()))

	// Job applications Handlers
	s.mux.Handle("GET /api/jobs",
		s.middlewareAuth(s.handleGetUserJobApplications()))
	s.mux.Handle("POST /api/jobs",
		s.middlewareAuth(s.handleCreateJobApplication()))
	s.mux.Handle("GET /api/jobs/{jobID}",
		s.middlewareAuth(s.handleGetJobApplication()))
	s.mux.Handle("PUT /api/jobs/{jobID}/status/{statusID}",
		s.middlewareAuth(s.handleUpdateJobApplicationStatus()))
	s.mux.Handle("DELETE /api/jobs/{jobID}",
		s.middlewareAuth(s.handleDeleteJobApplication()))

	// Job Application Notes Handlers
	// Same base path as jobs because all notes belong to a single job
	s.mux.Handle("GET /api/jobs/{jobID}/notes", s.middlewareAuth(s.handleGetJobNotes()))
	s.mux.Handle("POST /api/jobs/{jobID}/notes", s.middlewareAuth(s.handleCreateJobNote()))
	s.mux.Handle("GET /api/jobs/{jobID}/notes/{noteID}", s.middlewareAuth(s.handleGetJobNote()))
	s.mux.Handle("PUT /api/jobs/{jobID}/notes/{noteID}", s.middlewareAuth(s.handleUpdateJobNote()))
	s.mux.Handle("DELETE /api/jobs/{jobID}/notes/{noteID}", s.middlewareAuth(s.handleDeleteJobNote()))

	// admin handlers
	s.mux.Handle("GET /api/admin/users", s.middlewareAdminUser(s.handleGetAllUsers()))
	s.mux.Handle("GET /api/admin/users/{userID}", s.middlewareAdminUser(s.handleMakeUserAdmin(true)))
	s.mux.Handle("DELETE /api/admin/users/{userID}", s.middlewareAdminUser(s.handleMakeUserAdmin(false)))
	// TODO: add a restore company change endpoint for admins
	// TODO: add admin endpoints to get soft deleted application status and a way to restore them
	// TODO: add admin endpoints to get soft deleted Job application and a way to restore them
	// TODO: add admin endpoints to get soft deleted Job application Notes and a way to restore them
	// TODO: add admin endpoint to get server configs (VIEW ONLY)
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

	debugEnvVar := os.Getenv("JAT_DEV")
	var isDebug bool
	if debugEnvVar != "" {
		isDebug = true
	} else {
		isDebug = false
	}
	server := &Server{
		db:        cfg.GetDB(),
		mux:       http.NewServeMux(),
		port:      servePort,
		services:  services.NewServiceManager(cfg.GetDB()),
		logger:    newServerLogger(),
		jwtSecret: *cfg.SecretJWT,
		isDebug:   isDebug,
	}
	return server, nil
}
