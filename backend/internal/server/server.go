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
	handlerManager *handlers.Manager
	router         *mux.Router
}

// NewServer creates a new HTTP server
func NewServer(handlerManager *handlers.Manager) *Server {
	return &Server{
		handlerManager: handlerManager,
		router:         mux.NewRouter(),
	}
}

// SetupRoutes configures all the HTTP routes
func (s *Server) SetupRoutes() {
	// API routes
	api := s.router.PathPrefix("/api").Subrouter()
	
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
	
	// Environment configuration endpoint
	api.HandleFunc("/env", s.handlerManager.HandleGetEnvConfig).Methods("GET")

	// Add a simple test endpoint
	api.HandleFunc("/test", s.handleTest).Methods("GET")
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
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
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
	
	return http.ListenAndServe(":"+port, s.router)
}

// GetRouter returns the router for testing
func (s *Server) GetRouter() *mux.Router {
	return s.router
}