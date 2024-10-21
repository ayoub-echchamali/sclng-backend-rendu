package main

import (
	"os"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/sclng-backend-test-v1/api"
	"github.com/Scalingo/sclng-backend-test-v1/config"
)

func main() {
	log := logger.Default()
	log.Info("Initializing app")
	cfg, err := config.ReadConfig(false)
	if err != nil {
		log.WithError(err).Error("Fail to initialize configuration.")
		os.Exit(1)
	}
	server := api.NewApiServer(cfg)
	server.ServeAndListen(cfg.Port)
}


