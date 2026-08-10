package api

import (
	"log"
	"net/http"

	"github.com/filz0r/jat/internal/config"
	"github.com/filz0r/jat/internal/services"
	"gorm.io/gorm"
)

type Server struct {
	mux       *http.ServeMux
	db        *gorm.DB
	services  *services.ServiceManager
	port      string
	logger    *log.Logger
	server    *http.Server
	cfg       *config.ConfigFile
	jwtSecret string
	isDebug   bool
}
