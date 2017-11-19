package main

import (
	"os"
	"fmt"
	"flag"
	"github.com/simantovyousoufov/rd-client/pkg/gobrid"
	"os/exec"
	"runtime"
)

// @todo implement a `watch` that will watch for magnet links copied to clipboard and download them
// 		https://github.com/atotto/clipboard

var clear map[string]func()

func init() {
	clear = make(map[string]func())

	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cls") //Windows example it is untested, but I think its working
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func ClearScreen() {
	value, ok := clear[runtime.GOOS]

	if ! ok {
		panic("Platform is unsupported.")
	}

	value()
}

func requireSubcommand(commands map[string]CmdContainer) {
	if len(os.Args) < 2 {
		fmt.Println("A subcommand is required.\n")
		fmt.Println("Usage:\n")
		fmt.Println("ufo command [arg1] [arg2] ...\n")

		printDefaults(commands)

		os.Exit(1)
	}
}

func printDefaults(commands map[string]CmdContainer) {
	for key, val := range commands {
		fmt.Printf("Usage for `%s`:\n", key)
		val.cmd.PrintDefaults()
		fmt.Println()
	}
}

func getTokenFromEnv() string {
	token := os.Getenv(TOKEN_ENV)

	if token == "" {
		fmt.Printf("Missing the %s env variable", TOKEN_ENV)
		os.Exit(1)
	}

	return token
}

type CmdContainer struct {
	cmd   *flag.FlagSet
	flags map[string]*string // @todo this won't handle other flag types
}

func initFlags() map[string]CmdContainer {
	magnetCmd := flag.NewFlagSet("magnet", flag.ExitOnError)
	magnetCmdFile := magnetCmd.String("f", "", "File to load with newline separated links.")

	return map[string]CmdContainer{
		"magnet": {
			cmd: magnetCmd,
			flags: map[string]*string{
				"f": magnetCmdFile,
			},
		},
	}
}

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
