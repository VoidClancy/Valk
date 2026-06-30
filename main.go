package main

import (
	"fmt"
	"os"
	"slices"

	"valkyrie/cli"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	commands := cli.Commands

	if len(os.Args) < 2 {
		cli.PrintHelp()
		return
	}

	subcommand := os.Args[1]

	for _, cmd := range commands {
		if cmd.Name == subcommand || slices.Contains(cmd.Aliases, subcommand) {

			cmd.Callback(os.Args[2:])
			return
		}
	}

	fmt.Printf("Unknown command/flag: %s\n\n", subcommand)
	cli.PrintHelp()
}
