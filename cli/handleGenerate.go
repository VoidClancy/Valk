package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"valkyrie/generator"
	"valkyrie/schema"
)

func handleGenerate() {
	config := GetConfig()

	if err := os.MkdirAll(config.Output.Client, 0755); err != nil {
		fmt.Println(err)
		return
	}

	// Parse schema
	schemaBytes, err := os.ReadFile(config.Schema)
	if err != nil {
		fmt.Printf("failed to read schema file: %v\n", err)
		return
	}

	schemaDef, errs := schema.ParseSchema(string(schemaBytes))
	if len(errs) > 0 {
		fmt.Println("Schema parsing errors:")
		for _, e := range errs {
			fmt.Println(e)
		}
		return
	}

	relDir, err := filepath.Rel(config.Output.Client, config.Output.Migrations)
	var embedRelDir string
	if err == nil && !strings.HasPrefix(relDir, "..") && !filepath.IsAbs(relDir) {
		embedRelDir = filepath.ToSlash(filepath.Join(relDir, "*.sql"))
	} else {
		fmt.Printf("[WARNING]: Migrations directory %q is not a subdirectory of client output directory %q. Go's //go:embed does not support parent directory paths ('..'). Embedded migrations will be disabled.\n",
			config.Output.Migrations, config.Output.Client)
	}

	pkgName := filepath.Base(config.Output.Client)
	if pkgName == "." || pkgName == "" {
		pkgName = "valkyrie"
	}

	outputs, err := generator.GenerateClient(*schemaDef, pkgName, embedRelDir, config.Output.Migrations)
	if err != nil {
		fmt.Printf("failed to generate client: %v\n", err)
		return
	}

	for filename, content := range outputs {
		if err := os.WriteFile(filepath.Join(config.Output.Client, filename), []byte(content), 0644); err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("Generating Client...")
	fmt.Println("client generated at:", config.Output.Client)
}
