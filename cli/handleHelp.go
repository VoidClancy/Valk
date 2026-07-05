package cli

import (
	"fmt"
	"strings"
)

func PrintHelp() {

	fmt.Println("Usage: valkyrie <command/flag> [arguments]")
	fmt.Println("\nAvailable commands:")

	for _, cmd := range Commands {
		var aliasStr string
		if len(cmd.Aliases) > 0 {
			aliasStr = fmt.Sprintf(" (%s)", join(cmd.Aliases, ", "))
		}
		fmt.Printf("  %-12s %s\n", cmd.Name+aliasStr, cmd.Description)
	}
}

func join(slice []string, sep string) string {
	if len(slice) == 0 {
		return ""
	}
	var res strings.Builder
	res.WriteString(slice[0])
	for _, item := range slice[1:] {
		res.WriteString(sep)
		res.WriteString(item)
	}
	return res.String()
}
