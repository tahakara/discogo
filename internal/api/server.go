package api

import (
	"discogo/internal/config"
	"discogo/internal/logger"
	"discogo/internal/memcachedservice"
	"net/http"
	"os"
)

func StartHTTPServer(memcachedService *memcachedservice.MemcachedService) {
	// HTTP router başlat
	router := NewRouter(memcachedService)

	// HTTP sunucu adresini config'den al
	httpAddr, err := config.GetDiscoGoHTTPAddr()
	if err != nil {
		logger.Error("Failed to retrieve HTTP server address: " + err.Error())
		os.Exit(1)
	}

	logger.Info("Listening on " + httpAddr + "...")
	if err := http.ListenAndServe(httpAddr, router); err != nil {
		logger.Error("Failed to start HTTP server: " + err.Error())
		os.Exit(1)
	}
}
