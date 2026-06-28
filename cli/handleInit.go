package cli

import (
	"fmt"
	"os"
)

func handleInit() {
	os.WriteFile("valkyrie.json", []byte(configFileContent), 0644)
	fmt.Println("creating valkyrie.json ....")
}

var configFileContent string = `
{
  "database": {
    "url_env": "DATABASE_URL",
    "direct_url_env": "DATABASE_DIRECT_URL"
  },

  "schema": "./schema.prisma",

  "output": {
    "client": "./valkyrie",
    "migrations": "./valkyrie/migrations"
  }
}
	`
