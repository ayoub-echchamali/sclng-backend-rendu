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
	cfg, err := config.ReadConfig()
	if err != nil {
		log.WithError(err).Error("Fail to initialize configuration.")
		os.Exit(1)
	}
	server := api.NewApiServer(cfg.GithubToken)
	server.ServeAndListen(cfg.Port)
}


