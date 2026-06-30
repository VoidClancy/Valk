package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func handleGenerate() {
	config := GetConfig()

	if err := os.MkdirAll(config.Output.Client, 0755); err != nil {
		fmt.Println(err)
		return
	}

	relDir, err := filepath.Rel(config.Output.Client, config.Output.Migrations)
	if err != nil || strings.HasPrefix(relDir, "..") || filepath.IsAbs(relDir) {
		fmt.Printf("[WARNING]: Migrations directory %q is not a subdirectory of client output directory %q. Go's //go:embed does not support parent directory paths ('..').\n",
			config.Output.Migrations, config.Output.Client)
	}

	var content string

	content += fmt.Sprintf(`
	package valkyrie

	import "embed"

	//go:embed %s/*.sql
	var migrationsFS embed.FS
	
`, filepath.ToSlash(relDir))

	if err := os.WriteFile(filepath.Join(config.Output.Client, "client.go"), []byte(content), 0644); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Generating Client...")
	fmt.Println("client generated at:", config.Output.Client)
}
