package envconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Service handles environment configuration operations
type Service struct {
	installationPath string
}

// NewService creates a new environment configuration service
func NewService(installationPath string) *Service {
	return &Service{
		installationPath: installationPath,
	}
}

// GetEnvConfig loads and returns the current environment configuration
func (s *Service) GetEnvConfig() (*EnvConfigResponse, error) {
	if s.installationPath == "" {
		return nil, fmt.Errorf("no installation path configured")
	}

	// Find the .env file
	envFilePath, err := GetEnvFilePath(s.installationPath)
	if err != nil {
		return &EnvConfigResponse{
			Config:     nil,
			FilePath:   "",
			FileExists: false,
			Sections:   make(map[string][]EnvVar),
			Summary:    &ConfigSummary{},
		}, err
	}

	// Check if file exists
	fileInfo, err := os.Stat(envFilePath)
	if err != nil {
		return &EnvConfigResponse{
			Config:     nil,
			FilePath:   envFilePath,
			FileExists: false,
			Sections:   make(map[string][]EnvVar),
			Summary:    &ConfigSummary{},
		}, err
	}

	// Load the configuration
	config, err := LoadEnvFile(envFilePath)
	if err != nil {
		return &EnvConfigResponse{
			Config:       nil,
			FilePath:     envFilePath,
			FileExists:   true,
			LastModified: fileInfo.ModTime(),
			Sections:     make(map[string][]EnvVar),
			Summary:      &ConfigSummary{},
		}, err
	}

	return &EnvConfigResponse{
		Config:       config,
		FilePath:     envFilePath,
		FileExists:   true,
		LastModified: fileInfo.ModTime(),
		Sections:     config.GetVariablesBySection(),
		Summary:      config.GetConfigSummary(),
	}, nil
}

// UpdateEnvConfig updates the environment configuration
func (s *Service) UpdateEnvConfig(request *EnvUpdateRequest) error {
	if s.installationPath == "" {
		return fmt.Errorf("no installation path configured")
	}

	// Find the .env file
	envFilePath, err := GetEnvFilePath(s.installationPath)
	if err != nil {
		return err
	}

	// Load current configuration
	config, err := LoadEnvFile(envFilePath)
	if err != nil {
		return fmt.Errorf("failed to load current config: %w", err)
	}

	// Create backup if requested
	if request.CreateBackup {
		backupPath := fmt.Sprintf("%s.backup.%s", envFilePath, time.Now().Format("20060102-150405"))
		if err := copyFile(envFilePath, backupPath); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Update variables
	config.Variables = request.Variables

	// Save the updated configuration
	if err := config.SaveEnvFile(); err != nil {
		return fmt.Errorf("failed to save env file: %w", err)
	}

	return nil
}

// ValidateEnvConfig validates the environment configuration
func (s *Service) ValidateEnvConfig(variables []EnvVar) (*ValidationResult, error) {
	config := &EnvConfig{
		Variables: variables,
	}
	
	return ValidateEnvConfig(config), nil
}

// GetEnvFileContent returns the raw content of the .env file
func (s *Service) GetEnvFileContent() (string, error) {
	if s.installationPath == "" {
		return "", fmt.Errorf("no installation path configured")
	}

	envFilePath, err := GetEnvFilePath(s.installationPath)
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(envFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read env file: %w", err)
	}

	return string(content), nil
}

// SetEnvFileContent sets the raw content of the .env file
func (s *Service) SetEnvFileContent(content string, createBackup bool) error {
	if s.installationPath == "" {
		return fmt.Errorf("no installation path configured")
	}

	envFilePath, err := GetEnvFilePath(s.installationPath)
	if err != nil {
		return err
	}

	// Create backup if requested
	if createBackup {
		backupPath := fmt.Sprintf("%s.backup.%s", envFilePath, time.Now().Format("20060102-150405"))
		if err := copyFile(envFilePath, backupPath); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Write new content
	if err := os.WriteFile(envFilePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write env file: %w", err)
	}

	return nil
}

// CreateBackup creates a backup of the current .env file
func (s *Service) CreateBackup() (*BackupInfo, error) {
	if s.installationPath == "" {
		return nil, fmt.Errorf("no installation path configured")
	}

	envFilePath, err := GetEnvFilePath(s.installationPath)
	if err != nil {
		return nil, err
	}

	// Check if file exists
	if _, err := os.Stat(envFilePath); err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Create backup filename
	timestamp := time.Now().Format("20060102-150405")
	backupFilename := fmt.Sprintf("%s.backup.%s", filepath.Base(envFilePath), timestamp)
	backupPath := filepath.Join(filepath.Dir(envFilePath), backupFilename)

	// Create backup
	if err := copyFile(envFilePath, backupPath); err != nil {
		return nil, fmt.Errorf("failed to create backup: %w", err)
	}

	// Get backup file info
	backupInfo, err := os.Stat(backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get backup file info: %w", err)
	}

	return &BackupInfo{
		Filename:     backupFilename,
		FilePath:     backupPath,
		CreatedAt:    backupInfo.ModTime(),
		Size:         backupInfo.Size(),
		OriginalFile: envFilePath,
	}, nil
}

// ListBackups lists all backup files
func (s *Service) ListBackups() ([]*BackupInfo, error) {
	if s.installationPath == "" {
		return nil, fmt.Errorf("no installation path configured")
	}

	envFilePath, err := GetEnvFilePath(s.installationPath)
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(envFilePath)
	baseName := filepath.Base(envFilePath)
	
	// Read directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var backups []*BackupInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		// Check if this is a backup file
		name := entry.Name()
		if !isBackupFile(name, baseName) {
			continue
		}

		fullPath := filepath.Join(dir, name)
		fileInfo, err := entry.Info()
		if err != nil {
			continue
		}

		backups = append(backups, &BackupInfo{
			Filename:     name,
			FilePath:     fullPath,
			CreatedAt:    fileInfo.ModTime(),
			Size:         fileInfo.Size(),
			OriginalFile: envFilePath,
		})
	}

	return backups, nil
}

// RestoreBackup restores a backup file
func (s *Service) RestoreBackup(backupFilename string) error {
	if s.installationPath == "" {
		return fmt.Errorf("no installation path configured")
	}

	envFilePath, err := GetEnvFilePath(s.installationPath)
	if err != nil {
		return err
	}

	backupPath := filepath.Join(filepath.Dir(envFilePath), backupFilename)
	
	// Check if backup exists
	if _, err := os.Stat(backupPath); err != nil {
		return fmt.Errorf("backup file not found: %w", err)
	}

	// Create backup of current file before restoring
	currentBackup := fmt.Sprintf("%s.before-restore.%s", envFilePath, time.Now().Format("20060102-150405"))
	if err := copyFile(envFilePath, currentBackup); err != nil {
		return fmt.Errorf("failed to backup current file: %w", err)
	}

	// Restore the backup
	if err := copyFile(backupPath, envFilePath); err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	return nil
}

// SetInstallationPath updates the installation path
func (s *Service) SetInstallationPath(path string) {
	s.installationPath = path
}

// GetInstallationPath returns the current installation path
func (s *Service) GetInstallationPath() string {
	return s.installationPath
}

// Helper functions

func isBackupFile(filename, originalBasename string) bool {
	// Check for .backup suffix or .backup.timestamp pattern
	return (filename != originalBasename) && 
		   (strings.Contains(filename, ".backup"))
}