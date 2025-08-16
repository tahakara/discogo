package routes

import (
	"discogo/internal/logger"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	logger.Error("An error occurred while processing the request")
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
