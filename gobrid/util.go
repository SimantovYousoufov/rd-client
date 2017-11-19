package main

import (
	"os/exec"
	"os"
	"runtime"
	"fmt"
	"flag"
)

var clear map[string]func()

func init() {
	clear = make(map[string]func())

	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cls")
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
