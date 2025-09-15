package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

type Service struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Status struct {
	Running  bool      `json:"running"`
	Services []Service `json:"services"`
	Version  string    `json:"version"`
	Path     string    `json:"path"`
}

type Manager struct {
	dockerClient *client.Client
	setupPath    string
}

func NewManager() (*Manager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	// Find DDALAB setup path
	setupPath := findDDALABSetup()

	return &Manager{
		dockerClient: cli,
		setupPath:    setupPath,
	}, nil
}

func findDDALABSetup() string {
	// Check common locations
	paths := []string{
		"/DDALAB-setup",
		"../DDALAB-setup",
		"../../DDALAB-setup",
		os.Getenv("HOME") + "/DDALAB-setup",
		os.Getenv("HOME") + "/Desktop/DDALAB-setup",
	}

	for _, path := range paths {
		if _, err := os.Stat(filepath.Join(path, "docker-compose.yml")); err == nil {
			absPath, _ := filepath.Abs(path)
			return absPath
		}
	}

	return ""
}

func (m *Manager) getStatus() (*Status, error) {
	status := &Status{
		Running:  false,
		Services: []Service{},
		Version:  "Unknown",
		Path:     m.setupPath,
	}

	if m.setupPath == "" {
		status.Path = "Not found"
		return status, nil
	}

	// Get DDALAB services
	serviceNames := []string{"ddalab", "ddalab-postgres", "ddalab-redis", "ddalab-minio", "ddalab-traefik"}
	
	ctx := context.Background()
	containers, err := m.dockerClient.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	runningCount := 0
	for _, serviceName := range serviceNames {
		service := Service{
			Name:   serviceName,
			Status: "stopped",
		}

		for _, container := range containers {
			for _, name := range container.Names {
				if strings.Contains(name, serviceName) {
					if container.State == "running" {
						service.Status = "running"
						runningCount++
					}
					break
				}
			}
		}

		status.Services = append(status.Services, service)
	}

	status.Running = runningCount == len(serviceNames)

	// Try to get version from main container
	for _, container := range containers {
		for _, name := range container.Names {
			if strings.Contains(name, "ddalab") && !strings.Contains(name, "-") {
				if container.Image != "" {
					parts := strings.Split(container.Image, ":")
					if len(parts) > 1 {
						status.Version = parts[1]
					}
				}
				break
			}
		}
	}

	return status, nil
}

func (m *Manager) executeDockerCompose(args ...string) error {
	if m.setupPath == "" {
		return fmt.Errorf("DDALAB setup not found")
	}

	cmd := exec.Command("docker-compose", args...)
	cmd.Dir = m.setupPath
	cmd.Env = append(os.Environ(), fmt.Sprintf("COMPOSE_PROJECT_NAME=ddalab"))
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command failed: %v, output: %s", err, string(output))
	}

	return nil
}

func (m *Manager) handleStatus(w http.ResponseWriter, r *http.Request) {
	status, err := m.getStatus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (m *Manager) handleServiceAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["service"]
	action := vars["action"]

	var err error
	switch action {
	case "start":
		err = m.executeDockerCompose("start", serviceName)
	case "stop":
		err = m.executeDockerCompose("stop", serviceName)
	case "restart":
		err = m.executeDockerCompose("restart", serviceName)
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

func (m *Manager) handleStackAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["action"]

	var err error
	switch action {
	case "start":
		err = m.executeDockerCompose("up", "-d")
	case "stop":
		err = m.executeDockerCompose("down")
	case "restart":
		err = m.executeDockerCompose("restart")
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

func (m *Manager) handleLogs(w http.ResponseWriter, r *http.Request) {
	if m.setupPath == "" {
		http.Error(w, "DDALAB setup not found", http.StatusNotFound)
		return
	}

	cmd := exec.Command("docker-compose", "logs", "--tail=100")
	cmd.Dir = m.setupPath
	
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"logs": string(output)})
}

func (m *Manager) handleBackup(w http.ResponseWriter, r *http.Request) {
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

func enableCORS(next http.Handler) http.Handler {
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
}

func main() {
	manager, err := NewManager()
	if err != nil {
		log.Fatal("Failed to initialize manager:", err)
	}

	router := mux.NewRouter()
	
	// API routes
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/status", manager.handleStatus).Methods("GET")
	api.HandleFunc("/services/{service}/{action}", manager.handleServiceAction).Methods("POST")
	api.HandleFunc("/stack/{action}", manager.handleStackAction).Methods("POST")
	api.HandleFunc("/logs", manager.handleLogs).Methods("GET")
	api.HandleFunc("/backup", manager.handleBackup).Methods("POST")

	// Apply CORS middleware
	handler := enableCORS(router)

	// Check if running in Docker Desktop extension mode
	socketPath := os.Getenv("SOCKET_PATH")
	if socketPath != "" {
		// Unix socket mode for Docker Desktop
		os.Remove(socketPath)
		
		listener, err := net.Listen("unix", socketPath)
		if err != nil {
			log.Fatal("Failed to create unix socket:", err)
		}
		defer listener.Close()
		
		log.Printf("Server listening on unix socket: %s", socketPath)
		log.Fatal(http.Serve(listener, handler))
	} else {
		// TCP mode for development
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		
		log.Printf("Server listening on port %s", port)
		log.Fatal(http.ListenAndServe(":"+port, handler))
	}
}