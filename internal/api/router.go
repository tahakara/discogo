package api

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/tahakara/discogo/internal/api/routes"
	redisclient "github.com/tahakara/discogo/internal/redis"
)

func NewRouter(rclient redisclient.Client) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	router.HandleFunc("/disco/version", routes.VersionHandler).Methods("GET", "POST", "PUT")
	router.HandleFunc("/disco/health", routes.HealthCheckHandler).Methods("GET")

	router.HandleFunc("/disco/register", func(w http.ResponseWriter, r *http.Request) {
		routes.RegisterHandler(w, r, rclient)
	}).Methods("POST")

	router.HandleFunc("/disco/heartbeat/{uuid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		routes.HeartbeatHandler(w, r, rclient, vars["uuid"])
	}).Methods("POST")

	router.HandleFunc("/disco/discover", func(w http.ResponseWriter, r *http.Request) {
		routes.DiscoverHandler(w, r, rclient)
	}).Methods("GET")

	router.HandleFunc("/deregister", func(w http.ResponseWriter, r *http.Request) {
		routes.DeregisterHandler(w, r, rclient)
	}).Methods("POST")

	// TODO: will implement report handler
	// It handles reporting of service status from service clients
	// when a certain event occurs service health status changes etc.
	// router.HandleFunc("/report", routes.ReportHandler(memcachedService)).Methods("POST")

	// router.HandleFunc("/error", routes.ErrorHandler).Methods("GET")

	return router
}
