package handler

import (
	"net/http"
	"sp-slack/logger"
)

func RegisterRoutes() {
    http.HandleFunc("/", func (w http.ResponseWriter, r * http.Request) {
        logger.Info("ping")
    })
}

