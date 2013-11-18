package main

import (
	"fmt"
	"os"
	"github.com/jessevdk/go-flags"
	. "polydawn.net/docket/commands"
)

var parser = flags.NewNamedParser("docket", flags.Default | flags.HelpFlag)

const EXIT_BADARGS = 1
const EXIT_PANIC = 2

// print only the error message (don't dump stacks).
// unless any debug mode is on; then don't recover, because we want to dump stacks.
func panicHandler() {
	if len(os.Getenv("DEBUG")) == 0 {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			fmt.Println("\n" + "Docket crashed! This could be a problem with docker or git, or docket itself." + "\n" + "To see more about what went wrong, turn on stack traces by running:" + "\n\n" + "export DEBUG=1" + "\n\n" + "Feel free to contact the developers for help:" + "\n" + "https://github.com/polydawn/docket" + "\n")
			os.Exit(EXIT_PANIC)
		}
	}
}

func main() {
	defer panicHandler()

	//Go-flags is a little too clever with sub-commands.
	//To keep the help-command parity with git & docker / etc, check for 'help' manually before args parse
	if len(os.Args) < 2 || os.Args[1] == "help" {
		parser.WriteHelp(os.Stdout)
		os.Exit(0)
	}

	//Parse for command & flags, and exit with a relevant return code.
	_, err := parser.Parse()
	if err != nil {
		os.Exit(EXIT_BADARGS)
	} else {
		os.Exit(0)
	}
}

func init() {
	// parser.AddCommand(
	// 	"command",
	// 	"description",
	// 	"long description",
	// 	&whateverCmd{}
	// )
	parser.AddCommand(
		"run",
		"Run a container",
		"Run a container based on configuration in the current directory.",

		//Default settings
		&RunCmdOpts{
			Source:      "graph",
		},
	)
	parser.AddCommand(
		"build",
		"Transform a container",
		"Transform a container based on configuration in the current directory.",

		//Default settings
		&BuildCmdOpts{
			Source:      "graph",
			Destination: "graph",
		},
	)
	parser.AddCommand(
		"version",
		"Print docket version",
		"Print docket version",
		&VersionCmdOpts{},
	)
}
