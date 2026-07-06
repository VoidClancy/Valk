package cli

import (
	"encoding/json"
	"log"
	"os"
	"slices"
)

var LogLevels = []string{
	"query",
	"info",
	"warn",
	"error",
	"all",
	"none",
}

type Config struct {
	Database DatabaseConfig `json:"database"`
	Schema   string         `json:"schema"`
	Output   OutputConfig   `json:"output"`
	Log      []string       `json:"log"`
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
	configFile, err := os.ReadFile("valk.json")
	if err != nil {
		log.Fatal("valk.json not found")
		return nil
	}

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	// hasAll := false
	for _, l := range config.Log {
		if l == "all" {
			// hasAll = true
		}
		if !slices.Contains(LogLevels, l) && l != "all" {
			log.Fatalf("invalid log level in valk.json: %q (must be one of: query, info, warn, error, all)", l)
			return nil
		}
	}
	// if hasAll && len(config.Log) > 1 {
	// 	log.Fatal("invalid log configuration: 'all' must be the only log level specified")
	// 	return nil
	// }

	return &config
}
