package main

import (
	"os"

	"discogo/internal/api"
	cfg "discogo/internal/config"
	lg "discogo/internal/logger"
	lgm "discogo/internal/logger/messages"
	"discogo/internal/memcached"
	"discogo/internal/memcachedservice"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		lg.Error(lgm.ErrorLoadingEnvFile)
	}

	memcachedAddr, err := cfg.GetMemcachedServerAddr()
	if err != nil {
		lg.Error(lgm.ErrorMemcachedFailedToRetrieveServerAddress)
		os.Exit(1)
	}

	lg.Info(lgm.MessageR(lgm.InfoStartingServiceDiscovery, memcachedAddr))

	memcachedClient := memcached.NewMemcachedClient(memcachedAddr)
	if err := memcachedClient.Ping(); err != nil {
		lg.Error(lgm.ErrorMemcachedConnectionFailed)
		os.Exit(1)
	}

	memcachedService := memcachedservice.NewMemcachedService(memcachedAddr, memcachedClient)

	api.StartHTTPServer(memcachedService)
}
