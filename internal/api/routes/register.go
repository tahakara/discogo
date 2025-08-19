package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	requestDTOs "github.com/tahakara/discogo/internal/api/dtos/requestdto"
	validators "github.com/tahakara/discogo/internal/api/validators"
	env "github.com/tahakara/discogo/internal/config"
	logger "github.com/tahakara/discogo/internal/logger"
	redisclient "github.com/tahakara/discogo/internal/redis"
	redisHelper "github.com/tahakara/discogo/internal/redis/helper"
	utils "github.com/tahakara/discogo/internal/utils"
)

type RegisterResponse struct {
	Status           string   `json:"status"`                     // "success" or "error"
	ServiceUUID      string   `json:"serviceUUID,omitempty"`      // Only on success
	HealthCheckCycle int      `json:"healthCheckCycle,omitempty"` // Only on success
	Reason           []string `json:"reason,omitempty"`           // Only on error
}

// RegisterHandler handles service registration requests.
//
// @Summary      Register a new service
// @Description  Registers a new service instance with the discovery system.
// @Tags         DiscoGo
// @Accept       json
// @Produce      json
// @Param        request body requestdto.RegisterRequestDTO true "Service registration payload"
// @Success      200 {object} RegisterResponse
// @Failure      400 {object} RegisterResponse
// @Failure      409 {object} RegisterResponse
// @Router       /disco/register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request, rclient redisclient.Client) {
	startTime := time.Now()
	var req requestDTOs.RegisterRequestDTO
	if err := utils.DecodeJSONBody(w, r, &req); err != nil {
		// utils.DecodeJSONBody already writes error response
		utils.WriteJSONResponse(w, http.StatusBadRequest,
			RegisterResponse{
				Status: "error",
				Reason: []string{"Invalid request payload"},
			})
		return
	}

	if err := validators.ValidateRegisterRequest(&req); err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest,
			RegisterResponse{
				Status: "error",
				Reason: err,
			})
		return
	}

	mappedEntry := redisHelper.ServiceEntry{
		ServiceUUID:   uuid.New().String(),
		Name:          req.Name,
		Type:          req.Type,
		Status:        redisHelper.StatusRegistered,
		Version:       req.Version,
		Provider:      req.Provider,
		Region:        req.Region,
		Zone:          req.Zone,
		Cluster:       req.Cluster,
		InstanceID:    req.InstanceID,
		NetworkID:     req.NetworkID,
		SubnetID:      req.SubnetID,
		NetworkDomain: req.NetworkDomain,
		Tags:          req.Tags,
		Addr4:         req.Addr4,
		Addr6:         req.Addr6,
		Port4:         req.Port4,
		Port6:         req.Port6,
	}

	if exists, _ := redisHelper.IsServiceExists(rclient, mappedEntry); exists {
		utils.WriteJSONResponse(w, http.StatusConflict,
			RegisterResponse{
				Status: "error",
				Reason: []string{"Service already exists"},
			})
		return
	}

	redisHelper.RegisterNewService(rclient, mappedEntry)

	logger.Register(fmt.Sprintf("%s:%s", mappedEntry.Type, mappedEntry.ServiceUUID), time.Since(startTime))

	utils.WriteJSONResponse(w, http.StatusOK, RegisterResponse{
		Status:           "ok",
		ServiceUUID:      mappedEntry.ServiceUUID,
		HealthCheckCycle: env.GetHealthCheckInterval(), // Default health check cycle in seconds
	})
}
