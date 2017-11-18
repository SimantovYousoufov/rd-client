package main

import (
	"os"
	"fmt"
	"flag"
	"github.com/simantovyousoufov/rd-client/pkg/gobrid"
	"os/exec"
	"runtime"
)

// @todo clean up code

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

func main() {
	token := os.Getenv(TOKEN_ENV)

	if token == "" {
		fmt.Printf("Missing the %s env variable", TOKEN_ENV)
		os.Exit(1)
	}

	magnetCmd := flag.NewFlagSet("magnet", flag.ExitOnError)
	magnetCmdFile := magnetCmd.String("f", "", "File to load with newline separated links.")

	commands := map[string]*flag.FlagSet{
		"magnet": magnetCmd,
	}

	if len(os.Args) < 2 {
		fmt.Println("A subcommand is required.\n")
		fmt.Println("Usage:\n")
		fmt.Println("ufo command [arg1] [arg2] ...\n")

		for cmd, set := range commands {
			fmt.Printf("Usage for `%s`:\n", cmd)
			set.PrintDefaults()
			fmt.Println()
		}

		os.Exit(1)
	}

	subCommand := os.Args[1]

	switch subCommand {
	case "magnet":
		magnetCmd.Parse(os.Args[2:])

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
		fmt.Println("Not supported yet.")
		os.Exit(1)
	}
}

func HandleError(err error) {
	if err != nil {
		fmt.Printf("Encountered an error: %v", err)
	}
}
