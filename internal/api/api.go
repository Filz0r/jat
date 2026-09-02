package api

import (
	"errors"
	"net/http"
	"os"

	"github.com/filz0r/jat/internal/config"
	"github.com/filz0r/jat/internal/services"
	"github.com/filz0r/jat/internal/utils"
	"github.com/filz0r/jat/internal/version"
)

func (s *Server) Start() error {
	s.loadRoutes()

	// API routes live under /api and require the X-Jat-Client-Type header.
	s.mux.Handle("/api/", http.StripPrefix("/api", s.middlewareClientType(s.apiMux)))

	// Static files and SPA fallback are served at the root. In development the
	// Vite dev server handles this, so skip registration to avoid requiring
	// the client-type header on non-API requests.
	if !utils.IsDevMode() {
		s.mux.Handle("/{path...}", s.staticHandler())
	}

	s.server = &http.Server{
		Addr:    ":" + s.port,
		Handler: s.middlewareLogRequest(s.mux),
	}

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
	s.apiMux.Handle("GET /health", s.healthHandler())
	s.apiMux.Handle("GET /initialized", s.initializedHandler())
	s.apiMux.Handle("GET /initialized/set", s.middlewareAdminUser(s.handleSetInitialized()))

	//auth handlers
	s.apiMux.Handle("POST /auth/login", s.handleUserLogin())
	s.apiMux.Handle("POST /auth/logout",
		s.middlewareAuth(s.middlewareRefreshToken(s.handleUserLogout())))
	s.apiMux.Handle("GET /auth/refresh_token",
		s.middlewareRefreshToken(s.handleUserTokenRefresh()))
	s.apiMux.Handle("GET /auth/revoke_token",
		s.middlewareRefreshToken(s.handleUserRevokeToken()))
	// TODO: add a method to fetch the current sessions

	// user handlers
	s.apiMux.Handle("GET /users", s.middlewareAdminUser(s.handleGetAllUsers()))
	s.apiMux.Handle("POST /users", s.handleUserCreate())
	s.apiMux.Handle("GET /users/me", s.middlewareAuth(s.handleGetCurrentUser()))
	s.apiMux.Handle("GET /users/{userID}",
		s.middlewareAdminUser(s.handleGetSingleUser()))
	s.apiMux.Handle("PUT /users", s.middlewareAuth(s.handleUserUpdate()))

	// Company handlers
	s.apiMux.Handle("POST /company", s.middlewareAuth(s.handleCreateCompany()))
	s.apiMux.Handle("GET /company", s.middlewareAuth(s.handleGetAllCompanies()))
	s.apiMux.Handle("GET /company/{companyID}", s.middlewareAuth(s.handleGetACompany()))
	s.apiMux.Handle("PUT /company/{companyID}", s.middlewareAuth(s.handleUpdateACompany()))
	s.apiMux.Handle("DELETE /company/{companyID}", s.middlewareAdminUser(s.handleDeleteACompany()))
	s.apiMux.Handle("GET /company/{companyID}/history", s.middlewareAdminUser(s.handleGetCompanyHistory()))
	s.apiMux.Handle("PUT /company/{companyID}/history/{changeID}", s.middlewareAdminUser(s.handleRestoreACompanyChange()))
	// kinda irrelevant since the get /api/company/{companyID} has a query params that already returns count
	s.apiMux.Handle("GET /company/{companyID}/count", s.middlewareAuth(s.handleCountCompanyApplications()))

	// Application Statuses Handlers
	s.apiMux.Handle("GET /application_statuses", s.middlewareAuth(s.handleGetUserApplicationStatus()))
	s.apiMux.Handle("GET /application_statuses/{statusID}", s.middlewareAuth(s.handleGetAnApplicationStatus()))
	s.apiMux.Handle("PUT /application_statuses/{statusID}", s.middlewareAuth(s.handleUpdateApplicationStatus()))
	s.apiMux.Handle("DELETE /application_statuses/{statusID}", s.middlewareAuth(s.handleDeleteApplicationStatus()))
	s.apiMux.Handle("POST /application_statuses", s.middlewareAuth(s.handleCreateApplicationStatus()))
	s.apiMux.Handle("PUT /application_statuses/{statusID}/unarchive", s.middlewareAuth(s.handleUnarchiveJobApplicationStatus()))

	// Job applications Handlers
	s.apiMux.Handle("GET /jobs",
		s.middlewareAuth(s.handleGetUserJobApplications()))
	s.apiMux.Handle("POST /jobs",
		s.middlewareAuth(s.handleCreateJobApplication()))
	s.apiMux.Handle("GET /jobs/{jobID}",
		s.middlewareAuth(s.handleGetJobApplication()))
	s.apiMux.Handle("PUT /jobs/{jobID}/status/{statusID}",
		s.middlewareAuth(s.handleUpdateJobApplicationStatus()))
	s.apiMux.Handle("DELETE /jobs/{jobID}",
		s.middlewareAuth(s.handleDeleteJobApplication()))
	s.apiMux.Handle("GET /jobs/{jobID}/history",
		s.middlewareAuth(s.handleGetJobApplicationHistory()))

	// Job Application Notes Handlers
	// Same base path as jobs because all notes belong to a single job
	s.apiMux.Handle("GET /jobs/{jobID}/notes", s.middlewareAuth(s.handleGetJobNotes()))
	s.apiMux.Handle("POST /jobs/{jobID}/notes", s.middlewareAuth(s.handleCreateJobNote()))
	s.apiMux.Handle("GET /jobs/{jobID}/notes/{noteID}", s.middlewareAuth(s.handleGetJobNote()))
	s.apiMux.Handle("PUT /jobs/{jobID}/notes/{noteID}", s.middlewareAuth(s.handleUpdateJobNote()))
	s.apiMux.Handle("DELETE /jobs/{jobID}/notes/{noteID}", s.middlewareAuth(s.handleDeleteJobNote()))

	// admin handlers
	s.apiMux.Handle("GET /admin/users/{userID}", s.middlewareAdminUser(s.handleMakeUserAdmin(true)))
	s.apiMux.Handle("DELETE /admin/users/{userID}", s.middlewareAdminUser(s.handleMakeUserAdmin(false)))
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
		apiMux:    http.NewServeMux(),
		port:      servePort,
		services:  services.NewServiceManager(cfg.GetDB()),
		logger:    newServerLogger(),
		jwtSecret: *cfg.SecretJWT,
		isDebug:   isDebug,
	}
	return server, nil
}
