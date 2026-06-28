package cli

import (
	"flag"
	"fmt"
	"os"
)

func handleMigrate(args []string) {

	fs := flag.NewFlagSet("migrate", flag.ExitOnError)

	verbose := fs.Bool("v", false, "verbose logging")
	dbUrl := fs.String("u", "", "Database connection URL (defaults to defined Env Variable in valkyrie.json)")

	fs.Parse(args)

	if *dbUrl == "" {
		cfg := GetConfig()
		envVarName := "DATABASE_URL"

		if cfg != nil && cfg.Database.DirectURLEnv != "" {
			envVarName = cfg.Database.DirectURLEnv
		}

		*dbUrl = os.Getenv(envVarName)

		if *dbUrl == "" {
			fmt.Printf("Database connection URL is required (provide -u or set %s environment variable)\n", envVarName)
			return
		}
	}

	if *verbose {
		fmt.Println("Migrating database with verbose logging...")
		//  goose verbose logging
	} else {
		fmt.Println("Migrating database...")
	}
}
