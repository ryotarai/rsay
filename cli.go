package main

import (
	"flag"
	"fmt"
	"io"
	"log"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		queueName string
		voice     string
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.StringVar(&voice, "voice", "", "Voice name")
	flags.StringVar(&voice, "v", "", "Voice name(Short)")
	flags.StringVar(&queueName, "queue", "", "Queue name")
	flags.StringVar(&queueName, "q", "", "Queue name(Short)")

	flVersion := flags.Bool("version", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if *flVersion {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	if queueName == "" {
		log.Fatal("--queue is required")
	}

	repeater := &Repeater{
		queueName: queueName,
		voice:     voice,
	}
	repeater.StartLoop()

	return ExitCodeOK
}
