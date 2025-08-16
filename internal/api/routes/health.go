package routes

import (
	"discogo/internal/api/responses"
	rm "discogo/internal/api/responses/responseMessages"
	cfg "discogo/internal/config"
	lg "discogo/internal/logger"
	lgm "discogo/internal/logger/messages"
	"discogo/internal/memcached"
	"encoding/json"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	memcachedAddr, err := cfg.GetMemcachedServerAddr()
	if err != nil {
		lg.Error(lgm.ErrorMemcachedFailedToRetrieveServerAddress)

		w.WriteHeader(http.StatusInternalServerError)
		resp := responses.NewErrorResponse(
			lgm.ErrorMemcachedFailedToRetrieveServerAddress,
			"")

		json.NewEncoder(w).Encode(resp)
		return
	}

	newMemcachedClient := memcached.NewMemcachedClient(memcachedAddr)
	if err := newMemcachedClient.Ping(); err != nil {
		lg.Error(lgm.ErrorMemcachedPingFailed)
		w.WriteHeader(http.StatusInternalServerError)
		resp := responses.NewErrorResponse(
			rm.ServiceHealthCheckFailed,
			"",
		)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := responses.NewSuccessResponse(
		rm.ServiceHealthCheckSuccess,
	)
	json.NewEncoder(w).Encode(resp)
}
