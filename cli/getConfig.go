package cli

import (
	"encoding/json"
	"log"
	"os"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
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
	ClientName      string         `json:"client_name" yaml:"client_name"`
	EmbedMigrations *bool          `json:"embed_migrations" yaml:"embed_migrations"`
	Database        DatabaseConfig `json:"database" yaml:"database"`
	Schema          string         `json:"schema" yaml:"schema"`
	Output          OutputConfig   `json:"output" yaml:"output"`
	Log             []string       `json:"log" yaml:"log"`
}

type DatabaseConfig struct {
	URLEnv       string `json:"url_env" yaml:"url_env"`
	DirectURLEnv string `json:"direct_url_env" yaml:"direct_url_env"`
}

type OutputConfig struct {
	Client     string `json:"client" yaml:"client"`
	Migrations string `json:"migrations" yaml:"migrations"`
}

func GetConfig() *Config {
	var foundFiles []string
	candidateFiles := []string{"phi.yml", "phi.yaml", "phi.json"}

	for _, f := range candidateFiles {
		if _, err := os.Stat(f); err == nil {
			foundFiles = append(foundFiles, f)
		}
	}

	if len(foundFiles) > 1 {
		log.Fatalf("multiple configuration files found (%s). Please keep only one configuration file in the project directory.", strings.Join(foundFiles, ", "))
		return nil
	}

	if len(foundFiles) == 0 {
		log.Fatal("configuration file not found (expected phi.yml, phi.yaml, or phi.json)")
		return nil
	}

	foundFile := foundFiles[0]
	configFile, err := os.ReadFile(foundFile)
	if err != nil {
		log.Fatalf("failed to read %s: %v", foundFile, err)
		return nil
	}

	var config Config
	if strings.HasSuffix(foundFile, ".json") {
		err = json.Unmarshal(configFile, &config)
	} else {
		err = yaml.Unmarshal(configFile, &config)
	}

	if err != nil {
		log.Fatalf("failed to parse %s: %v", foundFile, err)
		return nil
	}

	// Apply default values if omitted
	if config.ClientName == "" {
		config.ClientName = "phi"
	}
	if config.EmbedMigrations == nil {
		defaultEmbed := true
		config.EmbedMigrations = &defaultEmbed
	}

	for _, l := range config.Log {
		if !slices.Contains(LogLevels, l) && l != "all" {
			log.Fatalf("invalid log level in %s: %q (must be one of: query, info, warn, error, all)", foundFile, l)
			return nil
		}
	}

	return &config
}
