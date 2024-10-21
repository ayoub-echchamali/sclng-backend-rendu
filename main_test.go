package main_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Scalingo/sclng-backend-test-v1/api"
	"github.com/Scalingo/sclng-backend-test-v1/config"
	log "github.com/sirupsen/logrus"
)

var testServer *api.ApiServer

func TestMain(m *testing.M) {
	cfg, err := config.ReadConfig(false)
	if err != nil {
		log.WithError(err).Error("Fail to initialize configuration.")
		os.Exit(1)
	}
    testServer = api.NewApiServer(cfg)
	testServer.CreateRoutes()
	code := m.Run()
	os.Exit(code)
}

func TestPing(t *testing.T) {
	req, _ := http.NewRequest("GET", "/ping", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetRepo(t *testing.T) {
	req, _ := http.NewRequest("GET", "/publicGithubRepos", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	testServer.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d, got %d\n", expected, actual)
	}
}