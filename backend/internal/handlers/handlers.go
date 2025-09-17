package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sdraeger/DDALAB-docker-ext/internal/models"
)

// Manager manages HTTP handlers and their dependencies
type Manager struct {
	dockerSvc models.DockerService
	pathSvc   models.PathService
	configSvc models.ConfigService
	healthSvc models.HealthService
	setupPath string
	configPath string
}

// NewManager creates a new handler manager
func NewManager(dockerSvc models.DockerService, pathSvc models.PathService, configSvc models.ConfigService, healthSvc models.HealthService, setupPath, configPath string) *Manager {
	return &Manager{
		dockerSvc:  dockerSvc,
		pathSvc:    pathSvc,
		configSvc:  configSvc,
		healthSvc:  healthSvc,
		setupPath:  setupPath,
		configPath: configPath,
	}
}

// UpdateSetupPath updates the setup path
func (m *Manager) UpdateSetupPath(path string) {
	m.setupPath = path
}

// GetSetupPath returns the current setup path
func (m *Manager) GetSetupPath() string {
	return m.setupPath
}

// HandleStatus handles status requests
func (m *Manager) HandleStatus(w http.ResponseWriter, r *http.Request) {
	// Log status requests to monitor polling frequency
	log.Printf("Status request from %s", r.RemoteAddr)
	
	status, err := m.dockerSvc.GetStatus(m.setupPath)
	if err != nil {
		log.Printf("ERROR: Failed to get status: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// HandleServiceAction handles individual service actions
func (m *Manager) HandleServiceAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["service"]
	action := vars["action"]

	var err error
	switch action {
	case "start":
		err = m.dockerSvc.ExecuteCompose(m.setupPath, "start", serviceName)
	case "stop":
		err = m.dockerSvc.ExecuteCompose(m.setupPath, "stop", serviceName)
	case "restart":
		err = m.dockerSvc.ExecuteCompose(m.setupPath, "restart", serviceName)
	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// HandleStackAction handles stack-wide actions
func (m *Manager) HandleStackAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["action"]
	
	log.Printf("=== STACK ACTION REQUESTED: %s ===", action)
	log.Printf("Setup path: %s", m.setupPath)

	var err error
	switch action {
	case "start":
		log.Printf("Starting DDALAB stack...")
		err = m.dockerSvc.ExecuteCompose(m.setupPath, "up", "-d")
	case "stop":
		log.Printf("Stopping DDALAB stack...")
		err = m.dockerSvc.ExecuteCompose(m.setupPath, "down")
	case "restart":
		log.Printf("Restarting DDALAB stack...")
		err = m.dockerSvc.ExecuteCompose(m.setupPath, "restart")
	default:
		log.Printf("Invalid action requested: %s", action)
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("ERROR executing %s: %v", action, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("=== STACK ACTION COMPLETED SUCCESSFULLY: %s ===", action)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// HandleLogs handles log retrieval
func (m *Manager) HandleLogs(w http.ResponseWriter, r *http.Request) {
	if m.setupPath == "" {
		http.Error(w, "DDALAB setup not found", http.StatusNotFound)
		return
	}

	// Determine project name
	projectName := "ddalab-deploy"
	if dirName := filepath.Base(m.setupPath); dirName != "" {
		projectName = strings.ReplaceAll(strings.ToLower(dirName), " ", "")
	}
	
	composeFile := filepath.Join(m.setupPath, "docker-compose.yml")

	// Try new docker compose syntax first
	cmd := exec.Command("docker", "compose", "-p", projectName, "-f", composeFile, "logs", "--tail=100")
	cmd.Dir = m.setupPath
	cmd.Env = os.Environ()
	
	output, err := cmd.Output()
	if err != nil {
		// Fallback to docker-compose
		cmd = exec.Command("docker-compose", "-p", projectName, "-f", composeFile, "logs", "--tail=100")
		cmd.Dir = m.setupPath
		cmd.Env = os.Environ()
		
		output, err = cmd.Output()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"logs": string(output)})
}

// HandleBackup handles backup creation
func (m *Manager) HandleBackup(w http.ResponseWriter, r *http.Request) {
	if m.setupPath == "" {
		http.Error(w, "DDALAB setup not found", http.StatusNotFound)
		return
	}

	backupScript := filepath.Join(m.setupPath, "scripts", "backup.sh")
	if _, err := os.Stat(backupScript); os.IsNotExist(err) {
		// Try the backup command directly
		cmd := exec.Command("docker-compose", "exec", "-T", "postgres", 
			"pg_dump", "-U", "ddalab", "-d", "ddalab")
		cmd.Dir = m.setupPath
		
		output, err := cmd.Output()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Save backup
		timestamp := time.Now().Format("20060102-150405")
		filename := fmt.Sprintf("ddalab-backup-%s.sql", timestamp)
		backupPath := filepath.Join(m.setupPath, "backups", filename)
		
		os.MkdirAll(filepath.Dir(backupPath), 0755)
		if err := os.WriteFile(backupPath, output, 0644); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"filename": filename})
	} else {
		// Use the backup script
		cmd := exec.Command("bash", backupScript)
		cmd.Dir = m.setupPath
		
		output, err := cmd.CombinedOutput()
		if err != nil {
			http.Error(w, fmt.Sprintf("Backup failed: %s", string(output)), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"output": string(output)})
	}
}

// HandleHealthCheck handles health check requests
func (m *Manager) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	services := []string{"ddalab", "ddalab-postgres", "ddalab-redis", "ddalab-minio", "ddalab-traefik"}
	health := m.healthSvc.CheckSystemHealth(ctx, services, m.dockerSvc)
	health.ConfigPath = m.setupPath

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// HandleMetrics handles metrics requests
func (m *Manager) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	metrics, err := m.dockerSvc.GetMetrics(ctx)
	if err != nil {
		http.Error(w, "Failed to get metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// HandleGetPaths handles get paths requests
func (m *Manager) HandleGetPaths(w http.ResponseWriter, r *http.Request) {
	config, err := m.configSvc.LoadExtensionConfig(m.configPath)
	if err != nil {
		log.Printf("Failed to load extension config: %v", err)
		config = &models.ExtensionConfig{}
	}

	if config.SelectedPath == "" {
		config.SelectedPath = m.setupPath
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// HandleValidatePath handles path validation requests
func (m *Manager) HandleValidatePath(w http.ResponseWriter, r *http.Request) {
	var req models.PathSelectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result := m.pathSvc.ValidatePath(req.Path)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// HandleSelectPath handles path selection requests
func (m *Manager) HandleSelectPath(w http.ResponseWriter, r *http.Request) {
	var req models.PathSelectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the path first
	validation := m.pathSvc.ValidatePath(req.Path)
	if !validation.Valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validation)
		return
	}

	// Save the selected path
	if err := m.pathSvc.SaveSelectedPath(m.configPath, req.Path); err != nil {
		http.Error(w, "Failed to save path", http.StatusInternalServerError)
		return
	}

	// Update manager's current path
	m.setupPath = req.Path
	log.Printf("Selected new DDALAB path: %s", req.Path)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(validation)
}

// HandleGetEnvConfig handles environment config requests
func (m *Manager) HandleGetEnvConfig(w http.ResponseWriter, r *http.Request) {
	if m.setupPath == "" {
		http.Error(w, "DDALAB setup not found", http.StatusNotFound)
		return
	}

	config, err := m.configSvc.GetEnvConfig(m.setupPath)
	if err != nil {
		log.Printf("Failed to get env config: %v", err)
		// Return default config on error
		config = &models.EnvConfig{
			URL:    "https://localhost",
			Host:   "localhost",
			Scheme: "https",
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// HandleDiscoverPaths handles path discovery requests
func (m *Manager) HandleDiscoverPaths(w http.ResponseWriter, r *http.Request) {
	discoveredPaths := m.pathSvc.DiscoverPaths()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{
		"discovered_paths": discoveredPaths,
	})
}