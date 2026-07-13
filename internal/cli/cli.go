package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"pryl/internal/clipboard"
	"pryl/internal/secret"
	"pryl/internal/timeutil"
	"pryl/internal/version"
)

const (
	timeUsageText   = "Usage: pryl time epoch [--unit seconds|milliseconds] <value>\n"
	secretUsageText = "Usage: pryl secret generate [--length N] [--alphabet alphanumeric|hex|base64url] [--print]\n"

	usageText = `Usage:
  pryl time epoch [--unit seconds|milliseconds] <value>
  pryl secret generate [--length N] [--alphabet alphanumeric|hex|base64url] [--print]

Commands:
  time      Time and timestamp utilities
  secret    Secure secret generation
  version   Show the CLI version
  help      Show this help
`
)

// Run executes the CLI and returns an error suitable for display to users.
func Run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		_, err := io.WriteString(stdout, usageText)
		return err
	}
	if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		if len(args) != 1 {
			return errors.New("help does not accept arguments")
		}
		_, err := io.WriteString(stdout, usageText)
		return err
	}
	if args[0] == "--version" || args[0] == "-v" || args[0] == "version" {
		if len(args) != 1 {
			return errors.New("version does not accept arguments")
		}
		_, err := fmt.Fprintln(stdout, version.Value)
		return err
	}

	switch args[0] {
	case "time":
		return runTime(args[1:], stdout, stderr)
	case "secret":
		return runSecret(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("unknown command %q (run 'pryl help' for usage)", args[0])
	}
}

func runTime(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		_, err := io.WriteString(stdout, timeUsageText)
		return err
	}
	if args[0] != "epoch" {
		return fmt.Errorf("unknown time command %q", args[0])
	}
	flags := flag.NewFlagSet("time epoch", flag.ContinueOnError)
	flags.SetOutput(stderr)
	unitName := flags.String("unit", "seconds", "seconds or milliseconds")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return errors.New("time epoch requires exactly one value")
	}

	value, err := strconv.ParseInt(flags.Arg(0), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid epoch value %q: expected an integer", flags.Arg(0))
	}

	var unit timeutil.Unit
	switch strings.ToLower(*unitName) {
	case "seconds", "second", "s":
		unit = timeutil.Seconds
	case "milliseconds", "millisecond", "ms":
		unit = timeutil.Milliseconds
	default:
		return fmt.Errorf("invalid time unit %q: expected seconds or milliseconds", *unitName)
	}

	result := timeutil.EpochToISO(value, unit)
	_, err = fmt.Fprintln(stdout, result)
	return err
}

func runSecret(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		_, err := io.WriteString(stdout, secretUsageText)
		return err
	}
	if args[0] != "generate" {
		return fmt.Errorf("unknown secret command %q", args[0])
	}

	flags := flag.NewFlagSet("secret generate", flag.ContinueOnError)
	flags.SetOutput(stderr)
	length := flags.Int("length", 32, "number of characters to generate")
	alphabet := flags.String("alphabet", "alphanumeric", "alphanumeric, hex, or base64url")
	printSecret := flags.Bool("print", false, "print the generated secret to the terminal instead of copying it")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected argument %q", flags.Arg(0))
	}
	if *length <= 0 {
		return errors.New("length must be greater than zero")
	}

	result, err := secret.Generate(*length, strings.ToLower(*alphabet))
	if err != nil {
		return err
	}
	if *printSecret {
		_, err = fmt.Fprintln(stdout, result)
		return err
	}
	if err := clipboard.Copy(result); err != nil {
		return err
	}
	_, err = fmt.Fprintln(stderr, "secret copied to clipboard")
	return err
}
