package main

import (
	"time"

	_ "github.com/tahakara/discogo/docs"
	env "github.com/tahakara/discogo/internal/config"
	serviceconfigloader "github.com/tahakara/discogo/internal/config/serviceconfiguration"
	"github.com/tahakara/discogo/internal/service"
)

func main() {
	env.LoadEnv()
	if serviceconfigloader.LoadAllConfigs() != nil {
		panic(true)
	}

	client := service.StartRedisService()
	client.Set("key", []byte("value"), 10*time.Minute)
	// client.Close() // Ensure the Redis client is closed when the application exits
	service.StartHTTPServer(client)
}
