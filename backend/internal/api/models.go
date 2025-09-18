package api

import "time"

// StandardResponse wraps all API responses
type StandardResponse struct {
	Success  bool        `json:"success"`
	Data     interface{} `json:"data,omitempty"`
	Error    *ErrorInfo  `json:"error,omitempty"`
	Metadata *Metadata   `json:"metadata"`
}

// ErrorInfo provides detailed error information
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Metadata provides response metadata
type Metadata struct {
	Timestamp     time.Time `json:"timestamp"`
	APIVersion    string    `json:"api_version"`
	ServerVersion string    `json:"server_version"`
}

// StatusResponse represents the detailed status
type StatusResponse struct {
	Running      bool               `json:"running"`
	State        string             `json:"state"` // up|down|starting|stopping|error
	Services     []ServiceStatus    `json:"services"`
	Installation InstallationInfo  `json:"installation"`
}

// ServiceStatus represents individual service status
type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"` // running|stopped|starting|error
	Health string `json:"health"` // healthy|unhealthy|starting|unknown
	Uptime string `json:"uptime,omitempty"`
}

// InstallationInfo represents installation details
type InstallationInfo struct {
	Path        string    `json:"path"`
	Version     string    `json:"version"`
	LastUpdated time.Time `json:"last_updated"`
	Valid       bool      `json:"valid"`
}

// ConfigurationResponse represents the system configuration
type ConfigurationResponse struct {
	InstallationPath string                 `json:"installation_path"`
	Environment      map[string]string      `json:"environment"`
	Features         map[string]bool        `json:"features"`
}

// LifecycleRequest represents lifecycle operation requests
type LifecycleRequest struct {
	Operation string            `json:"operation"` // start|stop|restart|update
	Options   map[string]string `json:"options,omitempty"`
}

// BackupResponse represents backup operation response
type BackupResponse struct {
	BackupID   string    `json:"backup_id"`
	Filename   string    `json:"filename"`
	CreatedAt  time.Time `json:"created_at"`
	Size       int64     `json:"size"`
	Path       string    `json:"path"`
}

// LogsResponse represents system logs
type LogsResponse struct {
	Logs      string    `json:"logs"`
	Service   string    `json:"service,omitempty"`
	Since     time.Time `json:"since,omitempty"`
	Lines     int       `json:"lines"`
}

// DetectionResult represents installation detection results
type DetectionResult struct {
	Found         bool     `json:"found"`
	Paths         []string `json:"paths"`
	Recommended   string   `json:"recommended,omitempty"`
	AutoSelected  bool     `json:"auto_selected"`
}