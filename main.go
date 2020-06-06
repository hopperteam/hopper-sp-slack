package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sp-slack/config"
	"sp-slack/handler"
	"sp-slack/logger"
	"time"
)

func main() {
	var err error

	err = config.Init()
	if err != nil {
		logger.Fatal(err)
	}

	handler.RegisterRoutes()

	server := createServer()
	// will terminate with interrupt signal
	startServer(server)
}

func createServer() *http.Server {
	var server *http.Server
	server = &http.Server{
		Addr: ":" + config.Port,
	}

	return server
}

// this function will not terminate until an interrupt signal
// is provided, thus this function should be called last
func startServer(server *http.Server) {
	var err error
	go func() {
		logger.Infof("startign server on port: %s", config.Port)
		if err = server.ListenAndServe(); err != nil {
			logger.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		logger.Error(err)
	}
}
