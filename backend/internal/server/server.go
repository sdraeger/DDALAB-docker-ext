package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sdraeger/DDALAB-docker-ext/internal/handlers"
)

// Server represents the HTTP server
type Server struct {
	handlerManager    *handlers.Manager
	lifecycleHandler  *handlers.LifecycleHandler
	envConfigHandler  *handlers.EnvConfigHandler
	router            *mux.Router
}

// NewServer creates a new HTTP server
func NewServer(handlerManager *handlers.Manager) *Server {
	// Get installation path from handler manager
	installationPath := handlerManager.GetSetupPath()
	
	return &Server{
		handlerManager:   handlerManager,
		lifecycleHandler: handlers.NewLifecycleHandler(installationPath),
		envConfigHandler: handlers.NewEnvConfigHandler(installationPath),
		router:           mux.NewRouter(),
	}
}

// SetupRoutes configures all the HTTP routes
func (s *Server) SetupRoutes() {
	// API version info
	s.router.HandleFunc("/api/version", s.handleVersion).Methods("GET")
	
	// V1 API routes
	api := s.router.PathPrefix("/api/v1").Subrouter()
	
	// Also support unversioned routes for backward compatibility (defaults to v1)
	apiCompat := s.router.PathPrefix("/api").Subrouter()
	
	// New standardized lifecycle endpoints
	api.HandleFunc("/status", s.lifecycleHandler.HandleStatusV1).Methods("GET")
	api.HandleFunc("/lifecycle/start", s.lifecycleHandler.HandleLifecycleStart).Methods("POST")
	api.HandleFunc("/lifecycle/stop", s.lifecycleHandler.HandleLifecycleStop).Methods("POST")
	api.HandleFunc("/lifecycle/restart", s.lifecycleHandler.HandleLifecycleRestart).Methods("POST")
	api.HandleFunc("/lifecycle/update", s.lifecycleHandler.HandleLifecycleUpdate).Methods("POST")
	api.HandleFunc("/lifecycle", s.lifecycleHandler.HandleGenericLifecycle).Methods("POST")
	api.HandleFunc("/logs", s.lifecycleHandler.HandleLogsV1).Methods("GET")
	
	// Legacy endpoints (still supported)
	api.HandleFunc("/services/{service}/{action}", s.handlerManager.HandleServiceAction).Methods("POST")
	api.HandleFunc("/stack/{action}", s.handlerManager.HandleStackAction).Methods("POST")
	api.HandleFunc("/backup", s.handlerManager.HandleBackup).Methods("POST")
	
	// Health and monitoring
	api.HandleFunc("/health", s.handlerManager.HandleHealthCheck).Methods("GET")
	api.HandleFunc("/metrics", s.handlerManager.HandleMetrics).Methods("GET")
	
	// Path management endpoints
	api.HandleFunc("/paths", s.handlerManager.HandleGetPaths).Methods("GET")
	api.HandleFunc("/paths/validate", s.handlerManager.HandleValidatePath).Methods("POST")
	api.HandleFunc("/paths/select", s.handlerManager.HandleSelectPath).Methods("POST")
	api.HandleFunc("/paths/discover", s.handlerManager.HandleDiscoverPaths).Methods("GET")
	
	// Environment configuration endpoints (new standardized API)
	api.HandleFunc("/config/env", s.envConfigHandler.HandleGetEnvConfig).Methods("GET")
	api.HandleFunc("/config/env", s.envConfigHandler.HandleUpdateEnvConfig).Methods("PUT")
	api.HandleFunc("/config/env/file", s.envConfigHandler.HandleGetEnvFile).Methods("GET")
	api.HandleFunc("/config/env/file", s.envConfigHandler.HandleUpdateEnvFile).Methods("PUT")
	api.HandleFunc("/config/env/validate", s.envConfigHandler.HandleValidateEnvConfig).Methods("POST")
	api.HandleFunc("/config/env/backup", s.envConfigHandler.HandleCreateEnvBackup).Methods("POST")
	api.HandleFunc("/config/env/backups", s.envConfigHandler.HandleListEnvBackups).Methods("GET")
	api.HandleFunc("/config/env/restore", s.envConfigHandler.HandleRestoreEnvBackup).Methods("POST")
	
	// Update endpoint
	api.HandleFunc("/update", s.handlerManager.HandleUpdateDDALAB).Methods("POST")

	// Add a simple test endpoint
	api.HandleFunc("/test", s.handleTest).Methods("GET")
	
	// Backward compatibility routes (unversioned - defaults to v1)
	s.setupBackwardCompatibilityRoutes(apiCompat)
}

// handleTest is a simple test endpoint
func (s *Server) handleTest(w http.ResponseWriter, r *http.Request) {
	log.Printf("TEST endpoint hit")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Backend is working!","time":"` + time.Now().String() + `"}`))
}

// EnableCORS adds CORS middleware
func (s *Server) EnableCORS() {
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	})
}

// Start starts the HTTP server
func (s *Server) Start(port string) error {
	log.Printf("DDALAB Manager backend starting...")
	log.Printf("DDALAB installation path: %s", s.handlerManager.GetSetupPath())
	log.Printf("Backend listening on port: %s", port)
	
	// Docker extensions run in an isolated VM/container context
	// They need to bind to all interfaces (0.0.0.0) within their container
	// Docker Desktop provides network isolation at the extension level
	return http.ListenAndServe(":"+port, s.router)
}

// handleVersion returns API version information
func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{
		"version": "1.0.0",
		"api_version": "v1",
		"supported_versions": ["v1"],
		"deprecated_versions": [],
		"server": "DDALAB Manager Backend",
		"features": {
			"status_monitoring": true,
			"env_configuration": true,
			"docker_management": true,
			"backup_restore": true,
			"path_management": true
		}
	}`))
}

// setupBackwardCompatibilityRoutes sets up unversioned routes for backward compatibility
func (s *Server) setupBackwardCompatibilityRoutes(api *mux.Router) {
	// Add version check middleware
	api.Use(s.addVersionHeaders)
	
	// Status and service management
	api.HandleFunc("/status", s.handlerManager.HandleStatus).Methods("GET")
	api.HandleFunc("/services/{service}/{action}", s.handlerManager.HandleServiceAction).Methods("POST")
	api.HandleFunc("/stack/{action}", s.handlerManager.HandleStackAction).Methods("POST")
	api.HandleFunc("/logs", s.handlerManager.HandleLogs).Methods("GET")
	api.HandleFunc("/backup", s.handlerManager.HandleBackup).Methods("POST")
	
	// Health and monitoring
	api.HandleFunc("/health", s.handlerManager.HandleHealthCheck).Methods("GET")
	api.HandleFunc("/metrics", s.handlerManager.HandleMetrics).Methods("GET")
	
	// Path management endpoints
	api.HandleFunc("/paths", s.handlerManager.HandleGetPaths).Methods("GET")
	api.HandleFunc("/paths/validate", s.handlerManager.HandleValidatePath).Methods("POST")
	api.HandleFunc("/paths/select", s.handlerManager.HandleSelectPath).Methods("POST")
	api.HandleFunc("/paths/discover", s.handlerManager.HandleDiscoverPaths).Methods("GET")
	
	// Environment configuration endpoints (legacy compatibility)
	api.HandleFunc("/env", s.envConfigHandler.HandleGetEnvConfigLegacy).Methods("GET")
	api.HandleFunc("/env/file", s.handlerManager.HandleGetEnvFile).Methods("GET")
	api.HandleFunc("/env/file", s.handlerManager.HandleUpdateEnvFile).Methods("PUT")
	api.HandleFunc("/env/validate", s.envConfigHandler.HandleValidateEnvConfigLegacy).Methods("POST")
	api.HandleFunc("/env/backup", s.handlerManager.HandleBackupEnvFile).Methods("POST")
	api.HandleFunc("/env/backups", s.handlerManager.HandleListEnvBackups).Methods("GET")
	api.HandleFunc("/env/restore", s.handlerManager.HandleRestoreEnvFile).Methods("POST")
	
	// Update endpoint
	api.HandleFunc("/update", s.handlerManager.HandleUpdateDDALAB).Methods("POST")
	
	// Test endpoint
	api.HandleFunc("/test", s.handleTest).Methods("GET")
}

// addVersionHeaders adds version information to response headers
func (s *Server) addVersionHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-API-Version", "v1")
		w.Header().Set("X-Server-Version", "1.0.0")
		next.ServeHTTP(w, r)
	})
}

// GetRouter returns the router for testing
func (s *Server) GetRouter() *mux.Router {
	return s.router
}