package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

func handleGenerate() {
	config := GetConfig()
	err := os.MkdirAll(config.Output.Client, 0755)

	if err != nil {
		fmt.Println(err)
		return
	}
	clientFile := filepath.Join(config.Output.Client, "client.go")

	err = os.WriteFile(clientFile, []byte("package valkyrie"), 0644)
	fmt.Println("Generating Client...")

	fmt.Println("client generated at: ", config.Output.Client)
}
