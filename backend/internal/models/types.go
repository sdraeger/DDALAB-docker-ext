package models

import "time"

// Service represents a Docker service status
type Service struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// Status represents the overall DDALAB system status
type Status struct {
	Running  bool      `json:"running"`
	Services []Service `json:"services"`
	Version  string    `json:"version"`
	Path     string    `json:"path"`
}

// LauncherConfig represents the DDALAB launcher configuration
type LauncherConfig struct {
	DDALABPath      string    `json:"ddalab_path"`
	FirstRun        bool      `json:"first_run"`
	LastOperation   string    `json:"last_operation"`
	Version         string    `json:"version"`
	AutoUpdateCheck bool      `json:"auto_update_check"`
	LastUpdateCheck time.Time `json:"last_update_check"`
}

// ExtensionConfig represents the Docker extension configuration
type ExtensionConfig struct {
	SelectedPath string   `json:"selected_path"`
	KnownPaths   []string `json:"known_paths"`
}

// PathValidationResult represents the result of path validation
type PathValidationResult struct {
	Valid           bool   `json:"valid"`
	Path            string `json:"path"`
	Message         string `json:"message"`
	HasCompose      bool   `json:"has_compose"`
	HasDDALABScript bool   `json:"has_ddalab_script"`
}

// PathSelectionRequest represents a request to select a path
type PathSelectionRequest struct {
	Path string `json:"path"`
}

// EnvConfig represents environment configuration from .env files
type EnvConfig struct {
	URL     string `json:"url"`
	Host    string `json:"host"`
	Port    string `json:"port"`
	Scheme  string `json:"scheme"`
	Domain  string `json:"domain"`
}

// HealthCheck represents a service health check result
type HealthCheck struct {
	Service string            `json:"service"`
	Healthy bool              `json:"healthy"`
	Status  string            `json:"status"`
	Message string            `json:"message,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

// SystemHealth represents the overall system health
type SystemHealth struct {
	Overall    bool          `json:"overall"`
	Services   []HealthCheck `json:"services"`
	Timestamp  time.Time     `json:"timestamp"`
	ConfigPath string        `json:"config_path"`
}