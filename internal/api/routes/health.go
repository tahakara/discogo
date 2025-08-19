package routes

import (
	"fmt"
	"net/http"
	"time"

	env "github.com/tahakara/discogo/internal/config"
	"github.com/tahakara/discogo/internal/logger"
	redisclient "github.com/tahakara/discogo/internal/redis"
	"github.com/tahakara/discogo/internal/utils"
)

type HealthStatus string

const (
	StatusHealthy   HealthStatus = "healthy"
	StatusUnhealthy HealthStatus = "unhealthy"
)

type HealthCheckResponse struct {
	Status HealthStatus `json:"status"`
}

func _startRedisService() redisclient.Client {
	startTime := time.Now()
	addr := env.GetRedisServerAddr()
	password := env.GetRedisPassword() // Şifre yoksa "" döndürsün
	db := env.GetRedisDB()             // Örn: 0

	client := redisclient.New(addr, password, db)
	err := client.Ping()

	if err != nil {
		logger.Error(fmt.Sprintf("Redis connection failed: %v", err), time.Since(startTime))
		return nil
	}
	logger.Info(fmt.Sprintf("Redis connected to %s", addr), time.Since(startTime))
	return client
}

// HealthCheck godoc
// @Summary      Health check endpoint
// @Description  Returns the health status of the API and Redis connection
// @Tags         DiscoGo
// @Produce      json
// @Success      200  {object}  HealthCheckResponse
// @Failure      500  {object}  HealthCheckResponse
// @Router       /disco/health [get]
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	client := _startRedisService()
	if client == nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, HealthCheckResponse{Status: StatusUnhealthy})
		return
	}
	client.Close()
	logger.HealthCheck("ok", time.Since(startTime))

	utils.WriteJSONResponse(w, http.StatusOK, HealthCheckResponse{Status: StatusHealthy})
}
