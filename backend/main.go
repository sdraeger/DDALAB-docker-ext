package main

import (
	"log"

	"github.com/sdraeger/DDALAB-docker-ext/internal/config"
	"github.com/sdraeger/DDALAB-docker-ext/internal/docker"
	"github.com/sdraeger/DDALAB-docker-ext/internal/handlers"
	"github.com/sdraeger/DDALAB-docker-ext/internal/health"
	"github.com/sdraeger/DDALAB-docker-ext/internal/paths"
	"github.com/sdraeger/DDALAB-docker-ext/internal/server"
)

func main() {
	// Initialize services
	dockerSvc, err := docker.NewService()
	if err != nil {
		log.Fatal("Failed to initialize Docker service:", err)
	}

	pathSvc := paths.NewService()
	configSvc := config.NewService()
	healthSvc := health.NewService()

	// Set up configuration paths
	configPath := "/tmp/ddalab-manager-config.json"

	// Load or find DDALAB setup path
	setupPath := pathSvc.LoadSelectedPath(configPath)
	if setupPath == "" {
		setupPath = pathSvc.FindDDALABSetup()
		if setupPath != "" {
			pathSvc.SaveSelectedPath(configPath, setupPath)
		}
	}

	// Initialize handler manager
	handlerManager := handlers.NewManager(dockerSvc, pathSvc, configSvc, healthSvc, setupPath, configPath)

	// Initialize and configure server
	srv := server.NewServer(handlerManager)
	srv.SetupRoutes()
	srv.EnableCORS()

	// Start server
	port := "8080"
	if err := srv.Start(port); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}
