package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"github.com/voidclancy/valk/generator"
	"github.com/voidclancy/valk/schema"
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
		pkgName = "valk"
	}

	parentImportPath, err := generator.ResolveImportPath(config.Output.Client)
	if err != nil {
		fmt.Printf("failed to resolve parent import path: %v\n", err)
		return
	}

	outputs, err := generator.GenerateClient(*schemaDef, pkgName, parentImportPath, embedRelDir, config.Output.Migrations, config.Log)
	if err != nil {
		fmt.Printf("failed to generate client: %v\n", err)
		return
	}

	for filename, content := range outputs {
		outPath := filepath.Join(config.Output.Client, filename)
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			fmt.Println(err)
			return
		}
		if err := os.WriteFile(outPath, []byte(content), 0644); err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("Generating Client...")
	fmt.Println("client generated at:", config.Output.Client)
}
