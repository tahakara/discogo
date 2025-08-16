package routes

import (
	"discogo/internal/logger"
	"discogo/internal/memcachedservice"
	"net/http"
)

func RegisterHandler(memcachedService *memcachedservice.MemcachedService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if memcachedService == nil {
			logger.Error("MemcachedService is nil")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Service registration successful"))
	}
}
