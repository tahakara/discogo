package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tahakara/discogo/internal/logger"
	redisclient "github.com/tahakara/discogo/internal/redis"
	redisHelper "github.com/tahakara/discogo/internal/redis/helper"
	"github.com/tahakara/discogo/internal/utils"
)

type HeartbeatResponse struct {
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

// HeartbeatHandler godoc
// @Summary      Heartbeat endpoint
// @Description  Checks the health of a service by UUID and updates its status in Redis.
// @Tags         DiscoGo
// @Accept       json
// @Produce      json
// @Param        uuid  query     string  true  "Service UUID"
// @Success      200   {object}  HeartbeatResponse
// @Failure      400   {object}  HeartbeatResponse
// @Failure      500   {object}  HeartbeatResponse
// @Router       /disco/heartbeat [post]
func HeartbeatHandler(w http.ResponseWriter, r *http.Request, rclient redisclient.Client, uuid string) {
	startTime := time.Now()
	if uuid == "" {
		utils.WriteJSONResponse(w, http.StatusBadRequest, HeartbeatResponse{
			Status: "error",
			Reason: "Missing uuid",
		})
		return
	}

	matched, err := utils.ValidateUUID(uuid)
	if err != nil || !matched {
		utils.WriteJSONResponse(w, http.StatusBadRequest, HeartbeatResponse{
			Status: "error",
			Reason: "Invalid uuid format",
		})
		return
	}

	updated, err := redisHelper.UpdateServiceEntry(rclient, uuid)
	if !updated {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, HeartbeatResponse{
			Status: "error",
			Reason: err.Error(),
		})
		return
	}

	logger.HeartBeat(fmt.Sprintf("%s healthy", uuid), time.Since(startTime))
	utils.WriteJSONResponse(w, http.StatusOK, HeartbeatResponse{
		Status: "ok",
	})

}
