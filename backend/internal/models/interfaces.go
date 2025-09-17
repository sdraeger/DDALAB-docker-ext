package models

import "context"

// DockerService defines the interface for Docker operations
type DockerService interface {
	GetStatus(setupPath string) (*Status, error)
	ExecuteCompose(setupPath string, args ...string) error
	CheckServiceHealth(ctx context.Context, serviceName string) HealthCheck
	CheckDDALABAPI() HealthCheck
	GetMetrics(ctx context.Context) (map[string]interface{}, error)
}

// PathService defines the interface for path management operations
type PathService interface {
	LoadSelectedPath(configPath string) string
	SaveSelectedPath(configPath, selectedPath string) error
	ValidatePath(path string) PathValidationResult
	FindDDALABSetup() string
	DiscoverPaths() []string
}

// ConfigService defines the interface for configuration operations
type ConfigService interface {
	ParseEnvFile(envPath string) (*EnvConfig, error)
	GetEnvConfig(setupPath string) (*EnvConfig, error)
	LoadExtensionConfig(configPath string) (*ExtensionConfig, error)
	SaveExtensionConfig(configPath string, config *ExtensionConfig) error
}

// HealthService defines the interface for health monitoring
type HealthService interface {
	CheckSystemHealth(ctx context.Context, services []string, dockerSvc DockerService) *SystemHealth
	CheckServiceHealth(ctx context.Context, serviceName string, dockerSvc DockerService) HealthCheck
}