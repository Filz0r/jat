package main

import (
	"log"
	"os"

	"github.com/filz0r/jat/internal/api"
	"github.com/filz0r/jat/internal/config"
	"github.com/joho/godotenv"
)

// @title Job Application Tracker (JAT) API
// @version 0.2.0
// @description REST API for the JAT job application tracker. The API requires the X-Jat-Client-Type header on all routes (values: web-client, tui-client). Web clients use HttpOnly Secure SameSite=Strict cookies; TUI clients receive tokens in the JSON response.
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	args := os.Args
	var cmd string
	if len(args) == 2 {
		cmd = args[1]
	} else if len(args) == 1 {
		log.Fatal("the TUI interface is disabled for now")
		// Config handle for TUI interface
		//cfg := config.New()
		//err = cfg.LoadData()
		//if err != nil {
		//	log.Fatal(err)
		//}
		//if cfg.IsInitialized() && (cfg.IsStandalone() || cfg.IsServer()) {
		//	var stdLogging bool
		//	if cfg.IsServer() {
		//		stdLogging = true
		//	} else {
		//		stdLogging = false
		//	}
		//
		//	db, err := database.ConnectDb(cfg.DbUri(), stdLogging)
		//	if err != nil {
		//		log.Fatal(err)
		//	}
		//	cfg.SetDB(db)
		//	cfg.Services = services.NewServiceManager(db)
		//}
		//
		//// TUI Handle
		//model := ui.NewModel(cfg)
		//p := tea.NewProgram(model)
		//if _, err := p.Run(); err != nil {
		//	log.Fatal(err)
		//}
		//return
	}
	if cmd == "" && len(args) == 2 {
		log.Fatal("nice try")
	}
	if cmd == "server" {
		cfg := config.New()
		err := cfg.LoadFromEnv()
		if err != nil {
			log.Fatal(err)
		}
		server, err := api.New(cfg)
		if err != nil {
			log.Fatal(err)
		}
		err = server.Start()
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Fatal("invalid command was given, use 'server' to run as server or no commands to run as a TUI interface")

}
