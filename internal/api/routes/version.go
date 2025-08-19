package routes

import (
	"net/http"

	env "github.com/tahakara/discogo/internal/config"
	"github.com/tahakara/discogo/internal/utils"
)

type VersionResponse struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	VersionName string `json:"versionName"`
}

// VersionHandler godoc
// @Summary      Get service version
// @Description  Retrieves the version information of the service
// @Tags         DiscoGo
// @Accept       json
// @Produce      json
// @Success      200 {object} VersionResponse
// @Failure      406 {object} utils.JSONResponse
// @Router       /disco/version [get]
func VersionHandler(w http.ResponseWriter, r *http.Request) {
	appName := env.GetDiscoGoName()
	appVersion := env.GetDiscoGoVersion()
	appVersionName := env.GetDiscoGoVersionName()

	if appName == "" || appVersion == "" || appVersionName == "" {
		utils.WriteJSONResponse(w, http.StatusNotAcceptable, nil)
		return
	}

	resp := VersionResponse{
		Name:        appName,
		Version:     appVersion,
		VersionName: appVersionName,
	}

	utils.WriteJSONResponse(w, http.StatusOK, resp)
}
