package main

import (
	"os"
	"fmt"
	"github.com/simantovyousoufov/rd-client/pkg/gobrid"
)

func main() {
	token := getTokenFromEnv()
	commands := initFlags()

	requireSubcommand(commands)

	subCommand := os.Args[1]

	switch subCommand {
	case "help":
		printDefaults(commands)
	case "magnet":
		magnetCmd := commands["magnet"]
		magnetCmdFile := magnetCmd.flags["f"]

		magnetCmd.cmd.Parse(os.Args[2:])

		m := NewMagnetCommand(MagnetCommandConfig{
			client: gobrid.NewClient(token),
		})

		if *magnetCmdFile != "" {
			err := m.AddLinksFromFile(*magnetCmdFile)
			HandleError(err)
		} else {
			m.AddLink(os.Args[2])
		}

		m.Download()
	default:
		fmt.Println("Unknown command.")

		printDefaults(commands)
		os.Exit(1)
	}
}

func HandleError(err error) {
	if err != nil {
		fmt.Printf("Encountered an error: %v", err)
	}
}
