package main

import (
	"flag"
	"fmt"
	"io"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
	ExitCodeWrongArguments
)

const DefaultIncludeSuffixes = "*"
const DefaultExcludeSuffixes = ".bin,.jpg,.jpeg,.png,.gif"
const DefaultKeySeparator = ","

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		params    string
		arguments []string
		version bool
        excludes string
        includes string
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.StringVar(&params, "params", "", "parameter for template files")
	flags.StringVar(&params, "p", "", "parameter for template files(Short)")

	flags.BoolVar(&version, "version", false, "Print version information and quit.")
	flags.BoolVar(&version, "v", false, "Print version information and quit(Short).")

    flags.StringVar(&includes, "includes", DefaultIncludeSuffixes, "Include filtering suffixes(e.g. .txt,.html)")
    flags.StringVar(&includes, "i", DefaultIncludeSuffixes, "Include filtering suffixes")

    flags.StringVar(&excludes, "excludes", DefaultExcludeSuffixes, "Exclude filtering suffixes")
    flags.StringVar(&excludes, "e", DefaultExcludeSuffixes, "Exclude filtering suffixes")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	for 0 < flags.NArg() {
		arguments = append(arguments, flags.Arg(0))
		flags.Parse(flags.Args()[1:])
	}

	if len(arguments) != 2 {
		fmt.Fprintln(cli.errStream, "gokeleton [-p params] <src-template> <dest-template>")
		return ExitCodeWrongArguments
	}

    startParams := StartParams{
        Keywords: params,
        Arguments: arguments,
        IncludeSuffixes: includes,
        ExcludeSuffixes: excludes,
        KeySeparator: DefaultKeySeparator}

	err := StartMain(startParams)
    if err != nil {
        return ExitCodeError
    }

	return ExitCodeOK
}
