package api

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	APIVersion    = "v1"
	ServerVersion = "1.0.0"
)

// Error codes
const (
	ErrCodeDDALABNotFound     = "DDALAB_NOT_FOUND"
	ErrCodeInvalidPath        = "INVALID_PATH"
	ErrCodeServiceError       = "SERVICE_ERROR"
	ErrCodeConfigError        = "CONFIG_ERROR"
	ErrCodeOperationFailed    = "OPERATION_FAILED"
	ErrCodeValidationFailed   = "VALIDATION_FAILED"
	ErrCodeInternalError      = "INTERNAL_ERROR"
)

// State constants
const (
	StateUp       = "up"
	StateDown     = "down"
	StateStarting = "starting"
	StateStopping = "stopping"
	StateError    = "error"
)

// Health constants
const (
	HealthHealthy   = "healthy"
	HealthUnhealthy = "unhealthy"
	HealthStarting  = "starting"
	HealthUnknown   = "unknown"
)

// Service status constants
const (
	StatusRunning  = "running"
	StatusStopped  = "stopped"
	StatusStarting = "starting"
	StatusError    = "error"
)

// ResponseHelper helps create standardized API responses
type ResponseHelper struct{}

// NewResponseHelper creates a new response helper
func NewResponseHelper() *ResponseHelper {
	return &ResponseHelper{}
}

// Success creates a successful response
func (rh *ResponseHelper) Success(w http.ResponseWriter, data interface{}) {
	response := StandardResponse{
		Success: true,
		Data:    data,
		Metadata: &Metadata{
			Timestamp:     time.Now(),
			APIVersion:    APIVersion,
			ServerVersion: ServerVersion,
		},
	}
	rh.writeJSON(w, http.StatusOK, response)
}

// Error creates an error response
func (rh *ResponseHelper) Error(w http.ResponseWriter, code string, message string, details string, statusCode int) {
	response := StandardResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
		Metadata: &Metadata{
			Timestamp:     time.Now(),
			APIVersion:    APIVersion,
			ServerVersion: ServerVersion,
		},
	}
	rh.writeJSON(w, statusCode, response)
}

// NotFound creates a not found error response
func (rh *ResponseHelper) NotFound(w http.ResponseWriter, message string, details string) {
	rh.Error(w, ErrCodeDDALABNotFound, message, details, http.StatusNotFound)
}

// BadRequest creates a bad request error response
func (rh *ResponseHelper) BadRequest(w http.ResponseWriter, message string, details string) {
	rh.Error(w, ErrCodeValidationFailed, message, details, http.StatusBadRequest)
}

// InternalError creates an internal server error response
func (rh *ResponseHelper) InternalError(w http.ResponseWriter, message string, details string) {
	rh.Error(w, ErrCodeInternalError, message, details, http.StatusInternalServerError)
}

// ServiceError creates a service error response
func (rh *ResponseHelper) ServiceError(w http.ResponseWriter, message string, details string) {
	rh.Error(w, ErrCodeServiceError, message, details, http.StatusServiceUnavailable)
}

// writeJSON writes a JSON response
func (rh *ResponseHelper) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If we can't encode the response, send a basic error
		http.Error(w, `{"success":false,"error":{"code":"ENCODING_ERROR","message":"Failed to encode response"}}`, http.StatusInternalServerError)
	}
}