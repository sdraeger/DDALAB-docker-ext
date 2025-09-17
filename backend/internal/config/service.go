package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sdraeger/DDALAB-docker-ext/internal/models"
)

// Service implements the ConfigService interface
type Service struct{}

// NewService creates a new config service
func NewService() *Service {
	return &Service{}
}

// ParseEnvFile parses a .env file and extracts configuration
func (s *Service) ParseEnvFile(envPath string) (*models.EnvConfig, error) {
	config := &models.EnvConfig{}
	
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return nil, fmt.Errorf(".env file not found at %s", envPath)
	}
	
	data, err := os.ReadFile(envPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read .env file: %w", err)
	}
	
	lines := strings.Split(string(data), "\n")
	envVars := make(map[string]string)
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			   (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}
			envVars[key] = value
		}
	}
	
	// Try to construct the URL from various possible environment variables
	// Check for direct PUBLIC_URL first
	if publicURL := envVars["PUBLIC_URL"]; publicURL != "" {
		config.URL = publicURL
		// Parse URL to extract components
		if strings.HasPrefix(publicURL, "https://") {
			config.Scheme = "https"
			config.Host = strings.TrimPrefix(publicURL, "https://")
		} else if strings.HasPrefix(publicURL, "http://") {
			config.Scheme = "http"
			config.Host = strings.TrimPrefix(publicURL, "http://")
		}
		// Extract port if present
		if strings.Contains(config.Host, ":") {
			parts := strings.Split(config.Host, ":")
			config.Host = parts[0]
			config.Port = parts[1]
		}
		config.Domain = config.Host
	} else if domain := envVars["DOMAIN"]; domain != "" {
		config.Domain = domain
		config.Host = domain
		config.Scheme = "https"
		config.URL = fmt.Sprintf("https://%s", domain)
	} else if domain := envVars["DDALAB_DOMAIN"]; domain != "" {
		config.Domain = domain
		config.Host = domain
		config.Scheme = "https"
		config.URL = fmt.Sprintf("https://%s", domain)
	} else if host := envVars["DDALAB_HOST"]; host != "" {
		config.Host = host
		scheme := "https"
		if s := envVars["DDALAB_SCHEME"]; s != "" {
			scheme = s
		}
		config.Scheme = scheme
		
		port := envVars["DDALAB_PORT"]
		if port != "" && port != "80" && port != "443" {
			config.Port = port
			config.URL = fmt.Sprintf("%s://%s:%s", scheme, host, port)
		} else {
			config.URL = fmt.Sprintf("%s://%s", scheme, host)
		}
	} else {
		// Default fallback
		config.URL = "https://localhost"
		config.Host = "localhost"
		config.Scheme = "https"
	}
	
	return config, nil
}

// GetEnvConfig gets environment configuration from setup path
func (s *Service) GetEnvConfig(setupPath string) (*models.EnvConfig, error) {
	if setupPath == "" {
		return nil, fmt.Errorf("setup path is empty")
	}

	// Try different possible .env file locations
	envPaths := []string{
		filepath.Join(setupPath, ".env"),
		filepath.Join(setupPath, ".env.master"),
		filepath.Join(setupPath, "ddalab-deploy", ".env"),
		filepath.Join(setupPath, "ddalab-deploy", ".env.master"),
	}

	var config *models.EnvConfig
	var err error
	
	for _, envPath := range envPaths {
		config, err = s.ParseEnvFile(envPath)
		if err == nil {
			return config, nil
		}
	}

	// Return default config if no .env found
	return &models.EnvConfig{
		URL:    "https://localhost",
		Host:   "localhost",
		Scheme: "https",
	}, nil
}

// LoadExtensionConfig loads extension configuration
func (s *Service) LoadExtensionConfig(configPath string) (*models.ExtensionConfig, error) {
	var config models.ExtensionConfig
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		// Return empty config if file doesn't exist
		return &config, nil
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// SaveExtensionConfig saves extension configuration
func (s *Service) SaveExtensionConfig(configPath string, config *models.ExtensionConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(configPath, data, 0644)
}