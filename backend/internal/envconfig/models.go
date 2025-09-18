package envconfig

import "time"

// EnvVar represents a single environment variable
type EnvVar struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	Comment    string `json:"comment,omitempty"`
	Section    string `json:"section,omitempty"`
	IsRequired bool   `json:"is_required"`
	IsSecret   bool   `json:"is_secret"`
	Example    string `json:"example,omitempty"`
}

// EnvConfig manages environment configuration
type EnvConfig struct {
	Variables []EnvVar `json:"variables"`
	FilePath  string   `json:"file_path"`
	Sections  []string `json:"sections"`
}

// EnvConfigResponse represents the API response for env config
type EnvConfigResponse struct {
	Config       *EnvConfig         `json:"config"`
	FilePath     string             `json:"file_path"`
	FileExists   bool               `json:"file_exists"`
	LastModified time.Time          `json:"last_modified,omitempty"`
	Sections     map[string][]EnvVar `json:"sections"`
	Summary      *ConfigSummary      `json:"summary"`
}

// ConfigSummary provides overview statistics
type ConfigSummary struct {
	TotalVariables    int `json:"total_variables"`
	RequiredVariables int `json:"required_variables"`
	SecretVariables   int `json:"secret_variables"`
	EmptyVariables    int `json:"empty_variables"`
	SectionCount      int `json:"section_count"`
}

// EnvUpdateRequest represents a request to update environment variables
type EnvUpdateRequest struct {
	Variables []EnvVar `json:"variables"`
	CreateBackup bool  `json:"create_backup,omitempty"`
}

// BackupInfo represents information about a backup file
type BackupInfo struct {
	Filename     string    `json:"filename"`
	FilePath     string    `json:"file_path"`
	CreatedAt    time.Time `json:"created_at"`
	Size         int64     `json:"size"`
	OriginalFile string    `json:"original_file"`
}

// ValidationResult represents validation results for env variables
type ValidationResult struct {
	Valid     bool      `json:"valid"`
	Errors    []string  `json:"errors,omitempty"`
	Warnings  []string  `json:"warnings,omitempty"`
	Variables []EnvVar  `json:"variables"`
}