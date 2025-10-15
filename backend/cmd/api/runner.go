package main

import (
	"log"
	"os"
	"path/filepath"

	"task-management/internal/config"
	"task-management/internal/infra/db"
	"task-management/internal/infra/server"
)

func main() {
	configPath, err := filepath.Abs("config")

	if err != nil {
		println("Gagal mendapatkan path absolute:", err.Error())
		os.Exit(1)
	}

	// Load configuration
	if err := config.LoadConfig(configPath); err != nil {
		println(configPath)
		println("Failed to load configuration file:", err.Error())
		os.Exit(1)
	}

	// init database
	database, err := db.NewDatabase(config.Config.Database)

	if err != nil {
		log.Fatalf("database setup failed: %v", err)
	}
	defer database.Close()

	app := server.StartServer()
	server.WaitForShutdown(app, func() {
		_ = database.Close()
	})
}
