package api

import (
	"discogo/internal/api/routes"
	"discogo/internal/memcachedservice"

	"github.com/gorilla/mux"
)

func NewRouter(memcachedService *memcachedservice.MemcachedService) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/health", routes.HealthCheckHandler).Methods("GET")
	router.HandleFunc("/version", routes.VersionHandler).Methods("GET")
	router.HandleFunc("/heartbeat", routes.HeartbeatHandler).Methods("GET")

	router.HandleFunc("/discover", routes.DiscoverHandler(memcachedService)).Methods("GET")
	router.HandleFunc("/register", routes.RegisterHandler(memcachedService)).Methods("POST")
	router.HandleFunc("/deregister", routes.DeregisterHandler(memcachedService)).Methods("POST")

	router.HandleFunc("/error", routes.ErrorHandler).Methods("GET")

	return router
}
