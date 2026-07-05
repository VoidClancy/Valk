package cli

type Command struct {
	Name        string
	Aliases     []string
	Description string
	Callback    func(args []string)
}

var Commands []Command

func init() {
	Commands = []Command{
		{
			Name:        "help",
			Aliases:     []string{"-h", "--help"},
			Description: "Print help info",
			Callback: func(args []string) {
				PrintHelp()
			},
		},
		{
			Name:        "init",
			Aliases:     []string{"-i"},
			Description: "Creates the config file in the root directory",
			Callback: func(args []string) {
				handleInit()
			},
		},
		{
			Name:        "migrate",
			Aliases:     []string{"-m"},
			Description: "Migrate the database",
			Callback:    handleMigrate,
		},
		{
			Name:        "generate",
			Aliases:     []string{"-g", "generate"},
			Description: "Generate migration files",
			Callback: func(args []string) {
				handleGenerate()
			},
		},
		{
			Name:        "reset",
			Aliases:     []string{"-r"},
			Description: "Reset the database to a clean state",
			Callback:    handleMigrate,
		},
	}
}
