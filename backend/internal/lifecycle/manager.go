package lifecycle

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sdraeger/DDALAB-docker-ext/internal/api"
)

// Manager handles DDALAB lifecycle operations
type Manager struct {
	installationPath string
}

// NewManager creates a new lifecycle manager
func NewManager(installationPath string) *Manager {
	return &Manager{
		installationPath: installationPath,
	}
}

// GetStatus returns the current status of DDALAB
func (m *Manager) GetStatus(ctx context.Context) (*api.StatusResponse, error) {
	if m.installationPath == "" {
		return nil, fmt.Errorf("no installation path configured")
	}

	// Check if services are running
	running, err := m.isRunning(ctx)
	if err != nil {
		return &api.StatusResponse{
			Running: false,
			State:   api.StateError,
			Services: []api.ServiceStatus{},
			Installation: api.InstallationInfo{
				Path:  m.installationPath,
				Valid: false,
			},
		}, err
	}

	// Get detailed service status
	services, err := m.getServiceStatus(ctx)
	if err != nil {
		// Return basic status even if detailed check fails
		services = []api.ServiceStatus{}
	}

	// Determine overall state
	state := api.StateDown
	if running {
		if m.allServicesHealthy(services) {
			state = api.StateUp
		} else {
			state = api.StateStarting
		}
	}

	return &api.StatusResponse{
		Running:  running,
		State:    state,
		Services: services,
		Installation: api.InstallationInfo{
			Path:        m.installationPath,
			Version:     "latest", // TODO: Get actual version
			LastUpdated: time.Now(), // TODO: Get actual last update time
			Valid:       true,
		},
	}, nil
}

// Start starts the DDALAB services
func (m *Manager) Start(ctx context.Context) error {
	if m.installationPath == "" {
		return fmt.Errorf("no installation path configured")
	}

	script := m.getScriptName()
	scriptPath := filepath.Join(m.installationPath, script)

	cmd := m.createCommandWithContext(ctx, scriptPath, "start")
	cmd.Dir = m.installationPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("failed to start DDALAB: %s\nOutput: %s", err, string(output))
	}

	return nil
}

// Stop stops the DDALAB services
func (m *Manager) Stop(ctx context.Context) error {
	if m.installationPath == "" {
		return fmt.Errorf("no installation path configured")
	}

	script := m.getScriptName()
	scriptPath := filepath.Join(m.installationPath, script)

	cmd := m.createCommandWithContext(ctx, scriptPath, "stop")
	cmd.Dir = m.installationPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("failed to stop DDALAB: %s\nOutput: %s", err, string(output))
	}

	return nil
}

// Restart restarts the DDALAB services
func (m *Manager) Restart(ctx context.Context) error {
	if err := m.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop services: %w", err)
	}

	// Wait a moment between stop and start
	time.Sleep(2 * time.Second)

	if err := m.Start(ctx); err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}

	return nil
}

// Update updates DDALAB to the latest version
func (m *Manager) Update(ctx context.Context) error {
	if m.installationPath == "" {
		return fmt.Errorf("no installation path configured")
	}

	// Stop services first
	if err := m.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop services before update: %w", err)
	}

	// Pull latest images
	cmd := exec.CommandContext(ctx, "docker-compose", "pull")
	cmd.Dir = m.installationPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("failed to pull latest images: %s\nOutput: %s", err, string(output))
	}

	// Start services with new images
	if err := m.Start(ctx); err != nil {
		return fmt.Errorf("failed to start services after update: %w", err)
	}

	return nil
}

// GetLogs retrieves DDALAB service logs
func (m *Manager) GetLogs(ctx context.Context, service string, lines int) (*api.LogsResponse, error) {
	if m.installationPath == "" {
		return nil, fmt.Errorf("no installation path configured")
	}

	args := []string{"logs"}
	if lines > 0 {
		args = append(args, "--tail", fmt.Sprintf("%d", lines))
	}
	if service != "" {
		args = append(args, service)
	}

	cmd := exec.CommandContext(ctx, "docker-compose", args...)
	cmd.Dir = m.installationPath

	output, err := cmd.Output()
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}

	return &api.LogsResponse{
		Logs:    string(output),
		Service: service,
		Since:   time.Now().Add(-1 * time.Hour), // Default to last hour
		Lines:   lines,
	}, nil
}

// isRunning checks if DDALAB services are currently running
func (m *Manager) isRunning(ctx context.Context) (bool, error) {
	cmd := exec.CommandContext(ctx, "docker-compose", "ps", "-q")
	cmd.Dir = m.installationPath

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check service status: %w", err)
	}

	// If there are running containers, output will not be empty
	return len(strings.TrimSpace(string(output))) > 0, nil
}

// getServiceStatus returns detailed status of all services
func (m *Manager) getServiceStatus(ctx context.Context) ([]api.ServiceStatus, error) {
	cmd := exec.CommandContext(ctx, "docker-compose", "ps")
	cmd.Dir = m.installationPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get service status: %w", err)
	}

	var services []api.ServiceStatus
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		// Skip header and empty lines
		if strings.Contains(line, "NAME") || strings.Contains(line, "---") || strings.TrimSpace(line) == "" {
			continue
		}

		// Look for lines that contain ddalab services
		if strings.Contains(line, "ddalab") {
			fields := strings.Fields(line)
			if len(fields) >= 6 {
				serviceName := fields[0]
				
				// Find the "Up" status by looking for it in the fields
				var status, health string
				for i, field := range fields {
					if field == "Up" || field == "Down" || field == "Exit" || field == "Restarting" {
						status = strings.ToLower(field)
						if status == "up" {
							status = api.StatusRunning
							health = api.HealthHealthy
						} else if status == "down" {
							status = api.StatusStopped
							health = api.HealthUnhealthy
						} else {
							status = api.StatusStarting
							health = api.HealthStarting
						}
						
						// Try to get uptime info
						uptime := ""
						if i+2 < len(fields) && (fields[i+2] == "hours" || fields[i+2] == "minutes" || fields[i+2] == "seconds") {
							uptime = fields[i+1] + " " + fields[i+2]
						}
						
						services = append(services, api.ServiceStatus{
							Name:   serviceName,
							Status: status,
							Health: health,
							Uptime: uptime,
						})
						break
					}
				}
			}
		}
	}

	return services, nil
}

// allServicesHealthy checks if all services are in a healthy state
func (m *Manager) allServicesHealthy(services []api.ServiceStatus) bool {
	if len(services) == 0 {
		return false
	}

	for _, service := range services {
		if service.Health != api.HealthHealthy {
			return false
		}
	}
	return true
}

// getScriptName returns the appropriate script name for the current OS
func (m *Manager) getScriptName() string {
	switch runtime.GOOS {
	case "windows":
		return "ddalab.ps1"
	default:
		return "ddalab.sh"
	}
}

// createCommandWithContext creates an appropriate command with context for the current OS
func (m *Manager) createCommandWithContext(ctx context.Context, scriptPath, action string) *exec.Cmd {
	switch runtime.GOOS {
	case "windows":
		if strings.HasSuffix(scriptPath, ".ps1") {
			return exec.CommandContext(ctx, "powershell", "-ExecutionPolicy", "Bypass", "-File", scriptPath, action)
		}
		return exec.CommandContext(ctx, scriptPath, action)
	default:
		return exec.CommandContext(ctx, "bash", scriptPath, action)
	}
}