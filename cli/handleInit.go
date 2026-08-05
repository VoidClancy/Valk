package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func handleInit(args []string) {
	format := "yml"
	targetDir := "."

	if len(args) == 1 {
		arg := strings.ToLower(args[0])
		if arg == "yml" || arg == "yaml" || arg == "json" {
			format = arg
		} else {
			targetDir = args[0]
		}
	} else if len(args) >= 2 {
		arg0 := strings.ToLower(args[0])
		if arg0 == "yml" || arg0 == "yaml" || arg0 == "json" {
			format = arg0
			targetDir = args[1]
		} else {
			targetDir = args[0]
		}
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", targetDir, err)
		os.Exit(1)
	}

	filename := "phi." + format
	outPath := filepath.Join(targetDir, filename)

	var content string
	if format == "json" {
		content = jsonConfigTemplate
	} else {
		content = ymlConfigTemplate
	}

	if err := os.WriteFile(outPath, []byte(strings.TrimSpace(content)+"\n"), 0644); err != nil {
		fmt.Printf("Error writing %s: %v\n", outPath, err)
		os.Exit(1)
	}

	fmt.Printf("Created configuration file: %s\n", outPath)
}

var ymlConfigTemplate string = `
# Name of the generated Go client package (default: "phi")
client_name: phi

# Enable Go 1.16+ embedded migrations (//go:embed) inside client
embed_migrations: true

database:
  # Environment variable for runtime client database connections
  url_env: DATABASE_URL
  # Environment variable for migration direct DDL connections
  direct_url_env: DATABASE_DIRECT_URL

# Path to your Prisma schema definition file
schema: ./schema.prisma

output:
  # Output directory for generated Go client code
  client: ./phi
  # Output directory for Goose-compatible .sql migration files
  migrations: ./phi/migrations

# Logging levels: none, query, warn, error, info, all
log:
  - none
`

var jsonConfigTemplate string = `
{
  "client_name": "phi",
  "embed_migrations": true,
  "database": {
    "url_env": "DATABASE_URL",
    "direct_url_env": "DATABASE_DIRECT_URL"
  },
  "schema": "./schema.prisma",
  "output": {
    "client": "./phi",
    "migrations": "./phi/migrations"
  },
  "log": [
    "none"
  ]
}
`
