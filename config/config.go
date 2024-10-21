package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port int `json:"port"`
	GithubToken string `json:"github_token"`
}

func ReadConfig(useDefault bool) (*Config, error) {
	// Default config
	defaultConfig := Config{
		Port: 5000,
		GithubToken: "",
	}
	if useDefault {
		return &defaultConfig, nil
	}
	absConfigPath, _ := filepath.Abs("config.json")
	file, err := os.Open(absConfigPath)
	if err != nil {
		log.Errorf("Config file not found. Using default config.")
		return &defaultConfig, nil
	}
	defer file.Close()

	jsonParser := json.NewDecoder(file)

	var cfg Config
	if err = jsonParser.Decode(&cfg); err != nil {
		log.Errorf("Error parsing config file. Error: %v", err)
		return &defaultConfig, err
	}
	return &cfg, nil
}
