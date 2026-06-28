package cli

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"valkyrie/migration"
	"valkyrie/schema"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type MigrationOptions struct {
	Verbose bool
	DBUrl   string
}

func handleMigrate(args []string) {
	isReset := false
	if len(os.Args) > 1 && (os.Args[1] == "reset" || os.Args[1] == "-r") {
		isReset = true
	} else if len(args) > 0 && args[0] == "reset" {
		isReset = true
	}

	opts, err := parseMigrationFlags(args)
	if err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	cfg := GetConfig()
	if cfg == nil {
		fmt.Println("Error: valkyrie.json not found or invalid")
		os.Exit(1)
	}

	dbUrl, err := resolveDbUrl(opts.DBUrl, cfg)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	schemaDef, err := parseSchemaFile(cfg.Schema)
	if err != nil {
		fmt.Printf("Error parsing schema: %v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open(string(schemaDef.Datasource.Provider), dbUrl)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if isReset {
		err = handleResetMigration(cfg, schemaDef, db, opts.Verbose)
	} else {
		err = handleUpMigration(cfg, schemaDef, db, opts.Verbose)
	}

	if err != nil {
		fmt.Printf("Migration failed: %v\n", err)
		os.Exit(1)
	}
}

func parseMigrationFlags(args []string) (*MigrationOptions, error) {
	fs := flag.NewFlagSet("migrate", flag.ExitOnError)
	verbose := fs.Bool("v", false, "verbose logging")
	dbUrl := fs.String("u", "", "Database connection URL")

	cleanArgs := []string{}
	for _, arg := range args {
		if arg != "reset" {
			cleanArgs = append(cleanArgs, arg)
		}
	}

	err := fs.Parse(cleanArgs)
	if err != nil {
		return nil, err
	}

	return &MigrationOptions{
		Verbose: *verbose,
		DBUrl:   *dbUrl,
	}, nil
}

func resolveDbUrl(flagUrl string, cfg *Config) (string, error) {
	if flagUrl != "" {
		return flagUrl, nil
	}

	envVarName := "DATABASE_DIRECT_URL"
	if cfg != nil && cfg.Database.DirectURLEnv != "" {
		envVarName = cfg.Database.DirectURLEnv
	}

	dbUrl := os.Getenv(envVarName)
	if dbUrl == "" {
		return "", fmt.Errorf("database connection URL is required (provide -u or set %s environment variable)", envVarName)
	}

	return dbUrl, nil
}

func parseSchemaFile(schemaPath string) (*schema.Schema, error) {
	schemaFileRaw, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file '%s': %w", schemaPath, err)
	}

	rawString := string(schemaFileRaw)
	parsedSchema, errs := schema.ParseSchema(rawString)
	if len(errs) > 0 {
		var sb strings.Builder
		sb.WriteString("schema diagnostics errors:\n")
		for _, diag := range errs {
			fmt.Fprintf(&sb, "- %s\n", diag.Message)
		}
		return nil, fmt.Errorf("%s", sb.String())
	}

	return parsedSchema, nil
}

func handleUpMigration(cfg *Config, schemaDef *schema.Schema, db *sql.DB, verbose bool) error {
	err := writeMigrationFile(cfg, schemaDef)
	if err != nil {
		return fmt.Errorf("failed to write migration file: %w", err)
	}

	goose.SetDialect(string(schemaDef.Datasource.Provider))

	if !verbose {
		goose.SetLogger(goose.NopLogger())
	}

	return goose.Up(db, cfg.Output.Migrations)
}

func handleResetMigration(cfg *Config, schemaDef *schema.Schema, db *sql.DB, verbose bool) error {
	goose.SetDialect(string(schemaDef.Datasource.Provider))

	if !verbose {
		goose.SetLogger(goose.NopLogger())
	}

	err := goose.Reset(db, cfg.Output.Migrations)
	if err != nil {
		return fmt.Errorf("failed to reset migrations: %w", err)
	}

	return goose.Up(db, cfg.Output.Migrations)
}

func writeMigrationFile(cfg *Config, schemaDef *schema.Schema) error {
	mig, err := migration.GenerateMigration(schemaDef)
	if err != nil {
		return err
	}

	err = os.MkdirAll(cfg.Output.Migrations, 0755)
	if err != nil {
		return fmt.Errorf("failed to create migrations directory '%s': %w", cfg.Output.Migrations, err)
	}

	filePath := filepath.Join(cfg.Output.Migrations, "001_migration.sql")
	err = os.WriteFile(filePath, mig, 0666)
	if err != nil {
		return fmt.Errorf("failed to write file '%s': %w", filePath, err)
	}

	return nil
}
