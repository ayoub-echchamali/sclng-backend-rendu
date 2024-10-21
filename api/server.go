package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/sclng-backend-test-v1/config"
	"github.com/Scalingo/sclng-backend-test-v1/util"
	log "github.com/sirupsen/logrus"
)

type ApiServer struct {
	Router *handlers.Router
	Config *config.Config
}

func NewApiServer(cfg *config.Config) *ApiServer {
	log.Info("creating server instance...")
	log := logger.Default()
	log.Infof("Creating api server with github token %s", util.Substring(cfg.GithubToken, 0, 9))
	return &ApiServer{
		Router: handlers.NewRouter(log),
		Config: cfg,
	}
}

func (s *ApiServer) ServeAndListen(port int) {
	log.Infof("initializing routes and listening on %v", port) 
	s.CreateRoutes()
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), s.Router)
	if err != nil {
		log.WithError(err).Error("Fail to listen to the given port")
		os.Exit(2)
	}
}