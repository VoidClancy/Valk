package cli

import (
	"fmt"
	"os"
)

func handleInit() {
	os.WriteFile("valk.json", []byte(configFileContent), 0644)
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
