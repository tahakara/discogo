package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tahakara/discogo/internal/api"
	env "github.com/tahakara/discogo/internal/config"
	"github.com/tahakara/discogo/internal/logger"
	redisclient "github.com/tahakara/discogo/internal/redis"
)

func StartHTTPServer(rclient redisclient.Client) {
	startTime := time.Now()
	addr := env.GetDiscoGoHTTPAddr()
	router := api.NewRouter(rclient)
	// Define your HTTP routes here
	logger.Info(fmt.Sprintf("HTTP server listening on %s", addr), time.Since(startTime))
	http.ListenAndServe(addr, router)
}

func StartRedisService() redisclient.Client {
	startTime := time.Now()
	addr := env.GetRedisServerAddr()
	password := env.GetRedisPassword() // Şifre yoksa "" döndürsün
	db := env.GetRedisDB()             // Örn: 0

	var rclient redisclient.Client = redisclient.New(addr, password, db)
	// Set dummy data for testing
	dummyData := []byte(`{"dummy":"data"}`)
	rclient.Set("asdasdas", dummyData, 1000000)

	err := rclient.Ping()
	if err != nil {
		logger.Error(fmt.Sprintf("Redis connection failed: %v", err), time.Since(startTime))
		return nil
	}
	logger.Info(fmt.Sprintf("Redis connected to %s", addr), time.Since(startTime))
	return rclient
}
