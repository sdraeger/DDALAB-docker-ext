package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sdraeger/DDALAB-docker-ext/internal/api"
	"github.com/sdraeger/DDALAB-docker-ext/internal/envconfig"
)

// EnvConfigHandler handles environment configuration operations
type EnvConfigHandler struct {
	envService     *envconfig.Service
	responseHelper *api.ResponseHelper
}

// NewEnvConfigHandler creates a new environment configuration handler
func NewEnvConfigHandler(installationPath string) *EnvConfigHandler {
	return &EnvConfigHandler{
		envService:     envconfig.NewService(installationPath),
		responseHelper: api.NewResponseHelper(),
	}
}

// HandleGetEnvConfig handles GET /api/v1/config/env
func (h *EnvConfigHandler) HandleGetEnvConfig(w http.ResponseWriter, r *http.Request) {
	envConfigResp, err := h.envService.GetEnvConfig()
	if err != nil {
		h.responseHelper.ServiceError(w, "Failed to get environment configuration", err.Error())
		return
	}

	h.responseHelper.Success(w, envConfigResp)
}

// HandleGetEnvFile handles GET /api/v1/config/env/file
func (h *EnvConfigHandler) HandleGetEnvFile(w http.ResponseWriter, r *http.Request) {
	content, err := h.envService.GetEnvFileContent()
	if err != nil {
		h.responseHelper.ServiceError(w, "Failed to get env file content", err.Error())
		return
	}

	response := map[string]interface{}{
		"content":       content,
		"last_modified": time.Now(), // TODO: Get actual modification time
	}

	h.responseHelper.Success(w, response)
}

// HandleUpdateEnvFile handles PUT /api/v1/config/env/file
func (h *EnvConfigHandler) HandleUpdateEnvFile(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Content      string `json:"content"`
		CreateBackup bool   `json:"create_backup,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.responseHelper.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	if err := h.envService.SetEnvFileContent(request.Content, request.CreateBackup); err != nil {
		h.responseHelper.ServiceError(w, "Failed to update env file", err.Error())
		return
	}

	h.responseHelper.Success(w, map[string]string{
		"message": "Environment file updated successfully",
	})
}

// HandleUpdateEnvConfig handles PUT /api/v1/config/env
func (h *EnvConfigHandler) HandleUpdateEnvConfig(w http.ResponseWriter, r *http.Request) {
	var request envconfig.EnvUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.responseHelper.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	if err := h.envService.UpdateEnvConfig(&request); err != nil {
		h.responseHelper.ServiceError(w, "Failed to update environment configuration", err.Error())
		return
	}

	h.responseHelper.Success(w, map[string]string{
		"message": "Environment configuration updated successfully",
	})
}

// HandleValidateEnvConfig handles POST /api/v1/config/env/validate
func (h *EnvConfigHandler) HandleValidateEnvConfig(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Variables []envconfig.EnvVar `json:"variables"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.responseHelper.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	result, err := h.envService.ValidateEnvConfig(request.Variables)
	if err != nil {
		h.responseHelper.ServiceError(w, "Failed to validate environment configuration", err.Error())
		return
	}

	h.responseHelper.Success(w, result)
}

// HandleCreateEnvBackup handles POST /api/v1/config/env/backup
func (h *EnvConfigHandler) HandleCreateEnvBackup(w http.ResponseWriter, r *http.Request) {
	backup, err := h.envService.CreateBackup()
	if err != nil {
		h.responseHelper.ServiceError(w, "Failed to create backup", err.Error())
		return
	}

	h.responseHelper.Success(w, backup)
}

// HandleListEnvBackups handles GET /api/v1/config/env/backups
func (h *EnvConfigHandler) HandleListEnvBackups(w http.ResponseWriter, r *http.Request) {
	backups, err := h.envService.ListBackups()
	if err != nil {
		h.responseHelper.ServiceError(w, "Failed to list backups", err.Error())
		return
	}

	h.responseHelper.Success(w, map[string]interface{}{
		"backups": backups,
		"count":   len(backups),
	})
}

// HandleRestoreEnvBackup handles POST /api/v1/config/env/restore
func (h *EnvConfigHandler) HandleRestoreEnvBackup(w http.ResponseWriter, r *http.Request) {
	var request struct {
		BackupFilename string `json:"backup_filename"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.responseHelper.BadRequest(w, "Invalid request body", err.Error())
		return
	}

	if request.BackupFilename == "" {
		h.responseHelper.BadRequest(w, "Backup filename is required", "")
		return
	}

	if err := h.envService.RestoreBackup(request.BackupFilename); err != nil {
		h.responseHelper.ServiceError(w, "Failed to restore backup", err.Error())
		return
	}

	h.responseHelper.Success(w, map[string]string{
		"message": "Backup restored successfully",
	})
}

// SetInstallationPath updates the installation path
func (h *EnvConfigHandler) SetInstallationPath(path string) {
	h.envService.SetInstallationPath(path)
}

// Legacy endpoints for backward compatibility

// HandleGetEnvConfigLegacy handles GET /api/env (legacy)
func (h *EnvConfigHandler) HandleGetEnvConfigLegacy(w http.ResponseWriter, r *http.Request) {
	h.HandleGetEnvConfig(w, r)
}

// HandleUpdateEnvConfigLegacy handles PUT /api/env (legacy)  
func (h *EnvConfigHandler) HandleUpdateEnvConfigLegacy(w http.ResponseWriter, r *http.Request) {
	h.HandleUpdateEnvConfig(w, r)
}

// HandleValidateEnvConfigLegacy handles POST /api/env/validate (legacy)
func (h *EnvConfigHandler) HandleValidateEnvConfigLegacy(w http.ResponseWriter, r *http.Request) {
	h.HandleValidateEnvConfig(w, r)
}