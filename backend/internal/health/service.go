package health

import (
	"context"
	"time"

	"github.com/sdraeger/DDALAB-docker-ext/internal/models"
)

// Service implements the HealthService interface
type Service struct{}

// NewService creates a new health service
func NewService() *Service {
	return &Service{}
}

// CheckSystemHealth checks the overall system health
func (s *Service) CheckSystemHealth(ctx context.Context, services []string, dockerSvc models.DockerService) *models.SystemHealth {
	health := &models.SystemHealth{
		Overall:   true,
		Timestamp: time.Now(),
	}

	healthChecks := make([]models.HealthCheck, 0, len(services))

	for _, service := range services {
		check := s.CheckServiceHealth(ctx, service, dockerSvc)
		healthChecks = append(healthChecks, check)
		
		// Update overall health
		if !check.Healthy {
			health.Overall = false
		}
	}

	// Try to check DDALAB API directly if main service is running
	if len(healthChecks) > 0 && healthChecks[0].Healthy {
		apiHealth := dockerSvc.CheckDDALABAPI()
		if apiHealth.Service != "" {
			healthChecks = append(healthChecks, apiHealth)
			if !apiHealth.Healthy {
				health.Overall = false
			}
		}
	}

	health.Services = healthChecks
	return health
}

// CheckServiceHealth checks the health of a specific service
func (s *Service) CheckServiceHealth(ctx context.Context, serviceName string, dockerSvc models.DockerService) models.HealthCheck {
	return dockerSvc.CheckServiceHealth(ctx, serviceName)
}