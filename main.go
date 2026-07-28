package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/filz0r/jat/internal/config"
	"github.com/filz0r/jat/internal/database"
	"github.com/filz0r/jat/internal/services"
	"github.com/filz0r/jat/internal/ui"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Config handle
	cfg := config.New()
	err = cfg.LoadData()
	if err != nil {
		log.Fatal(err)
	}
	if cfg.IsInitialized() && (cfg.IsStandalone() || cfg.IsServer()) {
		var stdLogging bool
		if cfg.IsServer() {
			stdLogging = true
		} else {
			stdLogging = false
		}

		db, err := database.ConnectDb(cfg.DbUri(), stdLogging)
		if err != nil {
			log.Fatal(err)
		}
		cfg.SetDB(db)
		cfg.Services = services.NewServiceManager(db)
	}

	// TUI Handle
	model := ui.NewModel(cfg)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
	// Server Handle to be implemented here
	// Basically create a factory that checks the conf file,
	// if its standalone/client load normal TUI
	// if its server launch the http server
}
