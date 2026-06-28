package cli

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Database DatabaseConfig `json:"database"`
	Schema   string         `json:"schema"`
	Output   OutputConfig   `json:"output"`
}

type DatabaseConfig struct {
	URLEnv       string `json:"url_env"`
	DirectURLEnv string `json:"direct_url_env"`
}

type OutputConfig struct {
	Client     string `json:"client"`
	Migrations string `json:"migrations"`
}

func GetConfig() *Config {
	var config Config
	configFile, err := os.ReadFile("valkyrie.json")
	if err != nil {
		log.Fatal("valkyrie.json not found")
		return nil
	}

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &config
}
