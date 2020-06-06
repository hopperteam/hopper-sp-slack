package handler

import (
	"net/http"
	"sp-slack/logger"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	logger.Info("ping")
}
