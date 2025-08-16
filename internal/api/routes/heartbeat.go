package routes

import (
	"discogo/internal/logger"
	"net/http"
)

func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received heartbeat request")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Heartbeat is alive"))
}
