package docker

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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"github.com/sdraeger/DDALAB-docker-ext/internal/models"
)

// Service implements the DockerService interface
type Service struct {
	client *client.Client
}

// NewService creates a new Docker service
func NewService() (*Service, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	return &Service{client: cli}, nil
}

// GetStatus retrieves the status of DDALAB services
func (s *Service) GetStatus(setupPath string) (*models.Status, error) {
	status := &models.Status{
		Running:  false,
		Services: []models.Service{},
		Version:  "Unknown",
		Path:     setupPath,
	}

	if setupPath == "" {
		status.Path = "Not found"
		return status, nil
	}

	// Get DDALAB services - look for common service names with various project prefixes
	baseServiceNames := []string{"ddalab", "postgres", "redis", "minio", "traefik"}
	knownDDALABContainers := []string{"ddalab", "ddalab-postgres", "ddalab-redis", "ddalab-minio", "ddalab-traefik"}
	
	// Also check for project-prefixed containers but be more specific
	projectPrefixes := []string{"ddalab-deploy", "ddalab", "ddalabsetup"}
	
	// Use context with timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	containers, err := s.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	// Create a map to track found services
	foundServices := make(map[string]models.Service)
	runningCount := 0
	
	// First, look for known DDALAB container names (without project prefix)
	for _, cont := range containers {
		for _, name := range cont.Names {
			// Remove leading slash from container name
			cleanName := strings.TrimPrefix(name, "/")
			
			// Check against known DDALAB container names
			for _, knownName := range knownDDALABContainers {
				if cleanName == knownName {
					// Extract service name (remove "ddalab-" prefix if present)
					serviceName := strings.TrimPrefix(cleanName, "ddalab-")
					
					service := models.Service{
						Name:   serviceName,
						Status: strings.ToLower(cont.State),
					}
					
					foundServices[serviceName] = service
					if cont.State == "running" {
						runningCount++
						log.Printf("Found running DDALAB service: %s (container: %s)", serviceName, cleanName)
					}
					break
				}
			}
		}
	}
	
	// If we didn't find services with exact names, try to find with project prefixes
	if len(foundServices) == 0 {
		// Try to detect project-prefixed containers, but be more strict
		for _, cont := range containers {
			for _, name := range cont.Names {
				cleanName := strings.TrimPrefix(name, "/")
				
				// Check if this looks like a DDALAB container with project prefix
				isDDALABContainer := false
				for _, prefix := range projectPrefixes {
					if strings.HasPrefix(cleanName, prefix+"-") || strings.HasPrefix(cleanName, prefix+"_") {
						isDDALABContainer = true
						break
					}
				}
				
				if isDDALABContainer {
					for _, baseName := range baseServiceNames {
						if strings.Contains(cleanName, baseName) {
							// Extract the service name
							serviceName := baseName
							
							service := models.Service{
								Name:   serviceName,
								Status: strings.ToLower(cont.State),
							}
							
							foundServices[serviceName] = service
							if cont.State == "running" {
								runningCount++
								log.Printf("Found running DDALAB service: %s (container: %s)", serviceName, cleanName)
							}
							break
						}
					}
				}
			}
		}
	}
	
	// Convert map to slice
	for _, service := range foundServices {
		status.Services = append(status.Services, service)
	}

	// Log what we found
	log.Printf("Found %d services, %d running", len(status.Services), runningCount)
	
	// Consider the stack running only if we have running services
	status.Running = runningCount > 0

	// Try to get version from main container
	for _, cont := range containers {
		for _, name := range cont.Names {
			if strings.Contains(name, "ddalab") && !strings.Contains(name, "-") {
				if cont.Image != "" {
					parts := strings.Split(cont.Image, ":")
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

// ExecuteCompose executes docker-compose commands
func (s *Service) ExecuteCompose(setupPath string, args ...string) error {
	if setupPath == "" {
		return fmt.Errorf("DDALAB setup not found")
	}

	// Check if docker-compose.yml exists
	composeFile := filepath.Join(setupPath, "docker-compose.yml")
	if _, err := os.Stat(composeFile); os.IsNotExist(err) {
		return fmt.Errorf("docker-compose.yml not found in %s", setupPath)
	}
	
	// Set up environment variables
	var envVars []string
	envVars = append(envVars, os.Environ()...)
	
	// Use the directory name as the project name to ensure consistency
	// This prevents Docker Desktop from creating extension-dependent project names
	setupDirName := filepath.Base(setupPath)
	projectName := strings.ToLower(strings.ReplaceAll(setupDirName, " ", ""))
	
	// Check if there are existing DDALAB containers to detect their project name
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	containers, err := s.client.ContainerList(ctx, container.ListOptions{All: true})
	if err == nil {
		// Look for existing DDALAB containers and use their project name
		for _, cont := range containers {
			// Check if this is a DDALAB container
			for _, name := range cont.Names {
				cleanName := strings.TrimPrefix(name, "/")
				if strings.Contains(cleanName, "ddalab") && !strings.Contains(cleanName, "manager") {
					// Check the compose project label
					if projectLabel, ok := cont.Labels["com.docker.compose.project"]; ok && projectLabel != "" {
						// Only use existing project if it's not from the extension
						if !strings.Contains(projectLabel, "desktop-extension") && !strings.Contains(projectLabel, "manager") {
							projectName = projectLabel
							log.Printf("Using existing DDALAB project name: %s", projectName)
							break
						}
					}
				}
			}
		}
	}
	
	// Log what project we're using
	log.Printf("Using project name: '%s' for command: %v", projectName, args)
	
	// Always use explicit project name to ensure independence from extension
	dockerArgs := []string{"compose", "-p", projectName, "-f", composeFile}
	dockerArgs = append(dockerArgs, args...)
	
	cmd := exec.Command("docker", dockerArgs...)
	cmd.Dir = setupPath
	cmd.Env = envVars
	
	log.Printf("Executing docker command: %v in directory: %s", dockerArgs, setupPath)
	
	output, err := cmd.CombinedOutput()
	log.Printf("Docker compose output: %s", string(output))
	
	if err != nil {
		// If docker compose fails, try docker-compose
		composeArgs := []string{"-p", projectName, "-f", composeFile}
		composeArgs = append(composeArgs, args...)
		
		cmd = exec.Command("docker-compose", composeArgs...)
		cmd.Dir = setupPath
		cmd.Env = envVars
		
		log.Printf("Fallback to docker-compose command: %v", composeArgs)
		
		output2, err2 := cmd.CombinedOutput()
		log.Printf("docker-compose output: %s", string(output2))
		
		if err2 != nil {
			// Return the most informative error
			return fmt.Errorf("docker compose failed: %v, output: %s\ndocker-compose failed: %v, output: %s", 
				err, string(output), err2, string(output2))
		}
	}

	return nil
}

// CheckServiceHealth checks the health of a specific service
func (s *Service) CheckServiceHealth(ctx context.Context, serviceName string) models.HealthCheck {
	health := models.HealthCheck{
		Service: serviceName,
		Healthy: false,
		Status:  "unknown",
		Details: make(map[string]string),
	}

	// Find the container
	containers, err := s.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		health.Status = "error"
		health.Message = "Failed to list containers"
		return health
	}

	var cont *types.Container
	for _, c := range containers {
		for _, name := range c.Names {
			if strings.Contains(name, serviceName) {
				cont = &c
				break
			}
		}
		if cont != nil {
			break
		}
	}

	if cont == nil {
		health.Status = "not_found"
		health.Message = "Container not found"
		return health
	}

	// Basic status
	health.Status = cont.State
	health.Details["image"] = cont.Image
	health.Details["created"] = time.Unix(cont.Created, 0).Format(time.RFC3339)

	if cont.State == "running" {
		// Check if container has health check
		inspect, err := s.client.ContainerInspect(ctx, cont.ID)
		if err == nil && inspect.State.Health != nil {
			health.Details["health_status"] = inspect.State.Health.Status
			if inspect.State.Health.Status == "healthy" {
				health.Healthy = true
			}
		} else if cont.State == "running" {
			// No health check, but running
			health.Healthy = true
		}

		// Add uptime
		if startedAt, err := time.Parse(time.RFC3339Nano, inspect.State.StartedAt); err == nil {
			health.Details["uptime"] = time.Since(startedAt).Round(time.Second).String()
		}
	}

	return health
}

// CheckDDALABAPI checks if the DDALAB API is responding
func (s *Service) CheckDDALABAPI() models.HealthCheck {
	check := models.HealthCheck{
		Service: "ddalab-api",
		Healthy: false,
		Status:  "unknown",
		Details: make(map[string]string),
	}

	// Try common DDALAB ports and endpoints
	endpoints := []string{
		"http://localhost:8000/health",
		"http://localhost:8000/api/health", 
		"http://localhost:3000/api/health",
		"https://localhost:8000/health",
		"https://localhost:8000/api/health",
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	for _, endpoint := range endpoints {
		resp, err := client.Get(endpoint)
		if err == nil {
			defer resp.Body.Close()
			check.Status = fmt.Sprintf("http_%d", resp.StatusCode)
			check.Details["endpoint"] = endpoint
			check.Details["status_code"] = fmt.Sprintf("%d", resp.StatusCode)
			
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				check.Healthy = true
				check.Message = "API responding"
				
				// Try to read response body for additional info
				var body map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&body); err == nil {
					check.Details["response"] = "valid_json"
				}
			}
			break
		}
	}

	if !check.Healthy && check.Status == "unknown" {
		check.Status = "unreachable"
		check.Message = "API not responding on any endpoint"
	}

	return check
}

// GetMetrics retrieves container metrics
func (s *Service) GetMetrics(ctx context.Context) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})
	
	// Get container stats
	containers, err := s.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get container list: %w", err)
	}

	for _, cont := range containers {
		for _, name := range cont.Names {
			if strings.Contains(name, "ddalab") {
				// Get container stats
				stats, err := s.client.ContainerStats(ctx, cont.ID, false)
				if err == nil {
					defer stats.Body.Close()
					var v types.StatsJSON
					if err := json.NewDecoder(stats.Body).Decode(&v); err == nil {
						containerMetrics := map[string]interface{}{
							"cpu_percent": calculateCPUPercent(&v),
							"memory_usage": v.MemoryStats.Usage,
							"memory_limit": v.MemoryStats.Limit,
						}
						
						// Add network stats if available
						if v.Networks != nil {
							if eth0, ok := v.Networks["eth0"]; ok {
								containerMetrics["network_rx"] = eth0.RxBytes
								containerMetrics["network_tx"] = eth0.TxBytes
							}
						}
						metrics[strings.TrimLeft(name, "/")] = containerMetrics
					}
				}
			}
		}
	}

	return metrics, nil
}

// calculateCPUPercent calculates CPU usage percentage
func calculateCPUPercent(stats *types.StatsJSON) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)
	
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		return (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return 0.0
}