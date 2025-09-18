package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/sdraeger/DDALAB-docker-ext/internal/api"
	"github.com/sdraeger/DDALAB-docker-ext/internal/lifecycle"
)

// LifecycleHandler handles lifecycle operations
type LifecycleHandler struct {
	lifecycleManager *lifecycle.Manager
	responseHelper   *api.ResponseHelper
}

// NewLifecycleHandler creates a new lifecycle handler
func NewLifecycleHandler(installationPath string) *LifecycleHandler {
	return &LifecycleHandler{
		lifecycleManager: lifecycle.NewManager(installationPath),
		responseHelper:   api.NewResponseHelper(),
	}
}

// HandleStatusV1 handles GET /api/v1/status
func (h *LifecycleHandler) HandleStatusV1(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	status, err := h.lifecycleManager.GetStatus(ctx)
	if err != nil {
		h.responseHelper.ServiceError(w, "Failed to get status", err.Error())
		return
	}

	h.responseHelper.Success(w, status)
}

// HandleLifecycleStart handles POST /api/v1/lifecycle/start
func (h *LifecycleHandler) HandleLifecycleStart(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()

	if err := h.lifecycleManager.Start(ctx); err != nil {
		h.responseHelper.ServiceError(w, "Failed to start DDALAB", err.Error())
		return
	}

	// Get updated status
	status, err := h.lifecycleManager.GetStatus(ctx)
	if err != nil {
		// Start succeeded but status check failed
		h.responseHelper.Success(w, map[string]string{
			"operation": "start",
			"result":    "success",
			"message":   "DDALAB started successfully",
		})
		return
	}

	h.responseHelper.Success(w, map[string]interface{}{
		"operation": "start",
		"result":    "success", 
		"message":   "DDALAB started successfully",
		"status":    status,
	})
}

// HandleLifecycleStop handles POST /api/v1/lifecycle/stop
func (h *LifecycleHandler) HandleLifecycleStop(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	if err := h.lifecycleManager.Stop(ctx); err != nil {
		h.responseHelper.ServiceError(w, "Failed to stop DDALAB", err.Error())
		return
	}

	h.responseHelper.Success(w, map[string]string{
		"operation": "stop",
		"result":    "success",
		"message":   "DDALAB stopped successfully",
	})
}

// HandleLifecycleRestart handles POST /api/v1/lifecycle/restart  
func (h *LifecycleHandler) HandleLifecycleRestart(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()

	if err := h.lifecycleManager.Restart(ctx); err != nil {
		h.responseHelper.ServiceError(w, "Failed to restart DDALAB", err.Error())
		return
	}

	// Get updated status
	status, err := h.lifecycleManager.GetStatus(ctx)
	if err != nil {
		// Restart succeeded but status check failed
		h.responseHelper.Success(w, map[string]string{
			"operation": "restart",
			"result":    "success",
			"message":   "DDALAB restarted successfully",
		})
		return
	}

	h.responseHelper.Success(w, map[string]interface{}{
		"operation": "restart",
		"result":    "success",
		"message":   "DDALAB restarted successfully", 
		"status":    status,
	})
}

// HandleLifecycleUpdate handles POST /api/v1/lifecycle/update
func (h *LifecycleHandler) HandleLifecycleUpdate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
	defer cancel()

	if err := h.lifecycleManager.Update(ctx); err != nil {
		h.responseHelper.ServiceError(w, "Failed to update DDALAB", err.Error())
		return
	}

	h.responseHelper.Success(w, map[string]string{
		"operation": "update",
		"result":    "success",
		"message":   "DDALAB updated successfully",
	})
}

// HandleLogsV1 handles GET /api/v1/logs
func (h *LifecycleHandler) HandleLogsV1(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Parse query parameters
	service := r.URL.Query().Get("service")
	linesStr := r.URL.Query().Get("lines")
	
	lines := 100 // default
	if linesStr != "" {
		if parsedLines, err := strconv.Atoi(linesStr); err == nil && parsedLines > 0 {
			lines = parsedLines
		}
	}

	logs, err := h.lifecycleManager.GetLogs(ctx, service, lines)
	if err != nil {
		h.responseHelper.ServiceError(w, "Failed to get logs", err.Error())
		return
	}

	h.responseHelper.Success(w, logs)
}

// HandleGenericLifecycle handles POST /api/v1/lifecycle with operation in body
func (h *LifecycleHandler) HandleGenericLifecycle(w http.ResponseWriter, r *http.Request) {
	var request api.LifecycleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.responseHelper.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	switch request.Operation {
	case "start":
		h.HandleLifecycleStart(w, r)
	case "stop":
		h.HandleLifecycleStop(w, r) 
	case "restart":
		h.HandleLifecycleRestart(w, r)
	case "update":
		h.HandleLifecycleUpdate(w, r)
	default:
		h.responseHelper.BadRequest(w, "Invalid operation", "Supported operations: start, stop, restart, update")
	}
}

// SetInstallationPath updates the installation path
func (h *LifecycleHandler) SetInstallationPath(path string) {
	h.lifecycleManager = lifecycle.NewManager(path)
}