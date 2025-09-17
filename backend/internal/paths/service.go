package paths

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sdraeger/DDALAB-docker-ext/internal/models"
)

// Service implements the PathService interface
type Service struct{}

// NewService creates a new path service
func NewService() *Service {
	return &Service{}
}

// LoadSelectedPath loads the selected path from config
func (s *Service) LoadSelectedPath(configPath string) string {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}

	var config models.ExtensionConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return ""
	}

	// Validate that the selected path still exists
	if s.ValidatePath(config.SelectedPath).Valid {
		return config.SelectedPath
	}

	return ""
}

// SaveSelectedPath saves the selected path to config
func (s *Service) SaveSelectedPath(configPath, selectedPath string) error {
	// Load existing config or create new one
	var config models.ExtensionConfig
	if data, err := os.ReadFile(configPath); err == nil {
		json.Unmarshal(data, &config)
	}

	config.SelectedPath = selectedPath

	// Add to known paths if not already present
	found := false
	for _, path := range config.KnownPaths {
		if path == selectedPath {
			found = true
			break
		}
	}
	if !found {
		config.KnownPaths = append(config.KnownPaths, selectedPath)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(configPath, data, 0644)
}

// ValidatePath validates a DDALAB installation path
func (s *Service) ValidatePath(path string) models.PathValidationResult {
	result := models.PathValidationResult{
		Valid: false,
		Path:  path,
	}

	if path == "" {
		result.Message = "Path cannot be empty"
		return result
	}

	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result.Message = "Path does not exist"
		return result
	}

	// Check for docker-compose.yml
	composePath := filepath.Join(path, "docker-compose.yml")
	if _, err := os.Stat(composePath); err == nil {
		result.HasCompose = true
	}

	// Check for DDALAB scripts
	scriptPaths := []string{
		filepath.Join(path, "ddalab.sh"),
		filepath.Join(path, "ddalab.ps1"),
		filepath.Join(path, "ddalab.bat"),
	}

	for _, scriptPath := range scriptPaths {
		if _, err := os.Stat(scriptPath); err == nil {
			result.HasDDALABScript = true
			break
		}
	}

	// Additional validation: check if docker-compose.yml contains DDALAB services
	if result.HasCompose {
		composeData, err := os.ReadFile(composePath)
		if err == nil {
			composeContent := string(composeData)
			hasServices := strings.Contains(composeContent, "ddalab") || 
						strings.Contains(composeContent, "postgres") ||
						strings.Contains(composeContent, "redis") ||
						strings.Contains(composeContent, "minio")
			
			if hasServices {
				result.Valid = true
				result.Message = "Valid DDALAB installation found"
			} else {
				result.Message = "docker-compose.yml found but doesn't appear to contain DDALAB services"
			}
		}
	} else {
		result.Message = "No docker-compose.yml found in the specified path"
	}

	return result
}

// FindDDALABSetup searches for DDALAB installation
func (s *Service) FindDDALABSetup() string {
	log.Printf("Starting DDALAB installation search...")
	
	// Try to read launcher config from host
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = "/root"
	}

	// Check launcher config
	launcherConfigPath := filepath.Join(homeDir, ".ddalab-launcher")
	if data, err := os.ReadFile(launcherConfigPath); err == nil {
		var config models.LauncherConfig
		if err := json.Unmarshal(data, &config); err == nil && config.DDALABPath != "" {
			if s.ValidatePath(config.DDALABPath).Valid {
				log.Printf("Found DDALAB path from launcher config: %s", config.DDALABPath)
				return config.DDALABPath
			}
		}
	}
	
	// Search common paths
	searchPaths := []string{
		"/Users/simon/Desktop/DDALAB-setup",
		homeDir + "/DDALAB-setup",
		homeDir + "/Desktop/DDALAB-setup",
		homeDir + "/Desktop/DDALAB",
		homeDir + "/Documents/DDALAB-setup",
		"/opt/DDALAB-setup",
		"/usr/local/DDALAB-setup",
	}

	for _, path := range searchPaths {
		if s.ValidatePath(path).Valid {
			log.Printf("Found DDALAB path: %s", path)
			return path
		}
	}

	log.Printf("No DDALAB installation found in standard locations")
	return ""
}

// DiscoverPaths discovers available DDALAB installation paths
func (s *Service) DiscoverPaths() []string {
	discoveredPaths := []string{}

	// Try to read launcher config
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = "/root"
	}

	launcherConfigPath := filepath.Join(homeDir, ".ddalab-launcher")
	if data, err := os.ReadFile(launcherConfigPath); err == nil {
		var config models.LauncherConfig
		if err := json.Unmarshal(data, &config); err == nil && config.DDALABPath != "" {
			if s.ValidatePath(config.DDALABPath).Valid {
				discoveredPaths = append(discoveredPaths, config.DDALABPath)
			}
		}
	}

	// Search common paths
	searchPaths := []string{
		"/Users/simon/Desktop/DDALAB-setup",
		"/Users/simon/Desktop/DDALAB",
		homeDir + "/DDALAB-setup",
		homeDir + "/Desktop/DDALAB-setup", 
		homeDir + "/Desktop/DDALAB",
		homeDir + "/Documents/DDALAB-setup",
		"/opt/DDALAB-setup",
		"/usr/local/DDALAB-setup",
	}

	for _, path := range searchPaths {
		if s.ValidatePath(path).Valid {
			// Avoid duplicates
			found := false
			for _, existing := range discoveredPaths {
				if existing == path {
					found = true
					break
				}
			}
			if !found {
				discoveredPaths = append(discoveredPaths, path)
			}
		}
	}

	return discoveredPaths
}