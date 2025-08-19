package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tahakara/discogo/internal/logger"
	redisclient "github.com/tahakara/discogo/internal/redis"
	redishelper "github.com/tahakara/discogo/internal/redis/helper"
	"github.com/tahakara/discogo/internal/utils"
)

type DeregisterRequestBody struct {
	ServiceUUID string `json:"serviceUUID"`
}

type DeregisterResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// @Summary Deregister a service
// @Description Deregisters a service from the registry using its UUID.
// @Tags DiscoGo
// @Accept json
// @Produce json
// @Param DeregisterRequestBody body DeregisterRequestBody true "Service UUID to deregister"
// @Success 200 {object} DeregisterResponse "Service deregistered successfully"
// @Failure 400 {object} DeregisterResponse "Invalid request body or serviceUUID"
// @Failure 404 {object} DeregisterResponse "Service not found"
// @Failure 500 {object} DeregisterResponse "Failed to deregister service"
// @Router /deregister [post]
func DeregisterHandler(w http.ResponseWriter, r *http.Request, rclient redisclient.Client) {
	startTime := time.Now()
	var body DeregisterRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, DeregisterResponse{
			Message: "Invalid request body",
			Status:  "error",
		})
		return
	}

	if body.ServiceUUID == "" {
		utils.WriteJSONResponse(w, http.StatusBadRequest, DeregisterResponse{
			Message: "serviceUUID is required",
			Status:  "error",
		})
		return
	}

	isValid, err := utils.ValidateUUID(body.ServiceUUID)
	if err != nil || !isValid {
		utils.WriteJSONResponse(w, http.StatusBadRequest, DeregisterResponse{
			Message: "Invalid serviceUUID format",
			Status:  "error",
		})
		return
	}

	result, err := redishelper.DeregisterServiceEntry(rclient, body.ServiceUUID)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, DeregisterResponse{
			Message: "Failed to deregister service",
			Status:  "error",
		})
		return
	}

	if !result {
		utils.WriteJSONResponse(w, http.StatusNotFound, DeregisterResponse{
			Message: "Service not found",
			Status:  "error",
		})
		return
	}

	logger.DeRegister(body.ServiceUUID, time.Since(startTime))
	utils.WriteJSONResponse(w, http.StatusOK, DeregisterResponse{
		Message: "Service deregistered successfully",
		Status:  "success",
	})
}
