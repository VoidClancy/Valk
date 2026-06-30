package cli

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"valkyrie/migration"
	"valkyrie/schema"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"golang.org/x/term"
	_ "modernc.org/sqlite"
)

type MigrationOptions struct {
	Verbose       bool
	DBUrl         string
	MigrationName string
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

	goose.SetDialect(string(schemaDef.Datasource.Provider))
	if !opts.Verbose {
		goose.SetLogger(goose.NopLogger())
	}

	if isReset {
		if isTTY() {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("[WARNING]: Resetting migrations will wipe ALL data in the database. Are you sure you want to proceed? [y/N]: ")
			response, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading response, aborting reset.")
				os.Exit(1)
			}
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Reset aborted.")
				os.Exit(0)
			}
		} else {
			fmt.Println("[WARNING]: Resetting migrations in non-interactive mode. ALL data in the database will be wiped.")
		}

		fmt.Println("Resetting database migrations...")
		err = goose.Reset(db, cfg.Output.Migrations)
		if err != nil && err != goose.ErrNoMigrationFiles {
			fmt.Printf("Error resetting migrations: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("Running pending migrations...")
	versionBefore, _ := goose.GetDBVersion(db)
	err = goose.Up(db, cfg.Output.Migrations)
	if err != nil && err != goose.ErrNoMigrationFiles {
		fmt.Printf("Migration failed: %v\n", err)
		os.Exit(1)
	}
	versionAfter, _ := goose.GetDBVersion(db)
	if versionAfter > versionBefore {
		fmt.Printf("Successfully applied pending migration(s) (database version is now %d).\n", versionAfter)
	}

	// diff current db schema against target schema.prisma
	upSql, downSql, err := migration.DiffAndPlan(db, schemaDef.Datasource.Provider, schemaDef, isTTY())
	if err != nil {
		fmt.Printf("Error generating migration plan: %v\n", err)
		os.Exit(1)
	}

	// if there is a diff, write a new migration file and apply it
	if len(strings.TrimSpace(upSql)) > 0 {
		migrationName := opts.MigrationName
		if migrationName == "" {
			if isTTY() {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter migration name: ")
				nameInput, err := reader.ReadString('\n')
				if err == nil {
					migrationName = strings.TrimSpace(nameInput)
				}
			}
			if migrationName == "" {
				migrationName = "migration_" + time.Now().Format("2006-01-02_150405")
			}
		}
		migrationName = strings.ReplaceAll(migrationName, " ", "_")

		err = writeMigrationFile(cfg, migrationName, upSql, downSql)
		if err != nil {
			fmt.Printf("Failed to write migration file: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Applying new migration...")
		err = goose.Up(db, cfg.Output.Migrations)
		if err != nil {
			fmt.Printf("Applying new migration failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Database migrated successfully.")
	} else {
		fmt.Println("No schema changes detected. Database is up-to-date.")
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

	migrationName := ""
	posArgs := fs.Args()
	if len(posArgs) > 0 {
		migrationName = posArgs[0]
	}

	return &MigrationOptions{
		Verbose:       *verbose,
		DBUrl:         *dbUrl,
		MigrationName: migrationName,
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

func writeMigrationFile(cfg *Config, migrationName, upSql, downSql string) error {
	err := os.MkdirAll(cfg.Output.Migrations, 0755)
	if err != nil {
		return fmt.Errorf("failed to create migrations directory '%s': %w", cfg.Output.Migrations, err)
	}

	filename := getNextMigrationFilename(cfg.Output.Migrations, migrationName)
	filePath := filepath.Join(cfg.Output.Migrations, filename)

	var sb strings.Builder
	sb.WriteString("-- +goose Up\n")
	sb.WriteString(upSql)
	sb.WriteString("\n\n-- +goose Down\n")
	sb.WriteString(downSql)
	sb.WriteString("\n")

	err = os.WriteFile(filePath, []byte(sb.String()), 0666)
	if err != nil {
		return fmt.Errorf("failed to write file '%s': %w", filePath, err)
	}

	fmt.Printf("Goose migration created successfully at: %s\n", filePath)
	return nil
}

func getNextMigrationFilename(migrationDir, migrationName string) string {
	files, err := os.ReadDir(migrationDir)
	if err != nil || len(files) == 0 {
		return fmt.Sprintf("00001_%s.sql", migrationName)
	}

	maxVersion := 0
	seqDigitLen := 5 //goose padding

	for _, entry := range files {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".sql") {
			continue
		}

		before, _, ok := strings.Cut(name, "_")
		if !ok {
			continue
		}

		prefix := before
		var val int
		_, err := fmt.Sscanf(prefix, "%d", &val)
		if err != nil {
			continue
		}

		if len(prefix) < 9 { // just a guard if i switch to timestamps
			if val > maxVersion {
				maxVersion = val
			}
			if len(prefix) > seqDigitLen {
				seqDigitLen = len(prefix)
			}
		}
	}

	nextVersion := maxVersion + 1
	formatStr := fmt.Sprintf("%%0%dd_%%s.sql", seqDigitLen)
	return fmt.Sprintf(formatStr, nextVersion, migrationName)
}

func isTTY() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}
