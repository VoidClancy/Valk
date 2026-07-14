package cli

import (
	"fmt"
	"os"
)

func handleInit() {
	err := os.WriteFile("valk.json", []byte(configFileContent), 0644)
	if err != nil {
		fmt.Printf("Error writing valk.json: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("creating valk.json ....")
}

var configFileContent string = `
{
  "database": {
    "url_env": "DATABASE_URL",
    "direct_url_env": "DATABASE_DIRECT_URL"
  },

  "schema": "./schema.prisma",

  "output": {
    "client": "./valk",
    "migrations": "./valk/migrations"
  }
}
	`
