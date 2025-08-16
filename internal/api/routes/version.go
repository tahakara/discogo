package routes

import (
	"discogo/internal/api/responses"
	rm "discogo/internal/api/responses/responseMessages"
	cfg "discogo/internal/config"
	lg "discogo/internal/logger"
	"encoding/json"
	"net/http"
)

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	version, verErr := cfg.GetDiscoGoVersion()
	appName, appErr := cfg.GetDiscoGoName()
	versionName, verNameErr := cfg.GetDiscoGoVersionName()

	switch {
	case verErr != nil:
		lg.Error(verErr.Error())
		resp := responses.NewErrorResponse(
			rm.MessageR(rm.FailedToGetTemplate, "version"),
			"")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	case appErr != nil:
		lg.Error(appErr.Error())
		resp := responses.NewErrorResponse(
			rm.MessageR(rm.FailedToGetTemplate, "app name"),
			"")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	case verNameErr != nil:
		lg.Error(verNameErr.Error())
		resp := responses.NewErrorResponse(
			rm.MessageR(rm.FailedToGetTemplate, "version name"),
			"")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := responses.NewVersionResponse(appName, version, versionName)
	json.NewEncoder(w).Encode(resp)
}
