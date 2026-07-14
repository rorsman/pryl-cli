package cli

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"pryl/internal/base64util"
	"pryl/internal/clipboard"
	"pryl/internal/hexutil"
	"pryl/internal/secret"
	"pryl/internal/timeutil"
	"pryl/internal/urlutil"
	"pryl/internal/version"
)

const (
	timeUsageText    = "Usage: pryl time epoch [--unit seconds|milliseconds] <value>\n"
	secretUsageText  = "Usage: pryl secret generate [--length N] [--alphabet alphanumeric|hex|base64url] [--print]\n"
	base64UsageText  = "Usage: pryl {encode|decode} base64 [--encoding standard|raw|url|rawurl] [value]\n"
	simpleCodecUsage = "Usage: pryl {encode|decode} {hex|url} [value]\n"

	usageText = `Usage:
  pryl time epoch [--unit seconds|milliseconds] <value>
  pryl secret generate [--length N] [--alphabet alphanumeric|hex|base64url] [--print]
  pryl encode base64 [--encoding standard|raw|url|rawurl] [value]
  pryl decode base64 [--encoding standard|raw|url|rawurl] [value]
  pryl encode hex [value]
  pryl decode hex [value]
  pryl encode url [value]
  pryl decode url [value]

Commands:
  time      Time and timestamp utilities
  secret    Secure secret generation
  version   Show the CLI version
  help      Show this help
`
)

// Run executes the CLI and returns an error suitable for display to users.
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		_, err := io.WriteString(stdout, usageText)
		return err
	}
	if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		if len(args) != 1 {
			return usageError("help does not accept arguments", usageText)
		}
		_, err := io.WriteString(stdout, usageText)
		return err
	}
	if args[0] == "--version" || args[0] == "-v" || args[0] == "version" {
		if len(args) != 1 {
			return usageError("version does not accept arguments", "Usage: pryl --version\n")
		}
		_, err := fmt.Fprintln(stdout, version.Value)
		return err
	}

	switch args[0] {
	case "time":
		return runTime(args[1:], stdout, stderr)
	case "secret":
		return runSecret(args[1:], stdout, stderr)
	case "encode":
		return runCodec(args[1:], true, stdin, stdout, stderr)
	case "decode":
		return runCodec(args[1:], false, stdin, stdout, stderr)
	default:
		return usageError(fmt.Sprintf("unknown command %q", args[0]), usageText)
	}
}

func usageError(message, usage string) error {
	return fmt.Errorf("%s\n%s", message, usage)
}

func optionValue(args []string, index *int, option string) (string, error) {
	argument := args[*index]
	if strings.HasPrefix(argument, option+"=") {
		value := strings.TrimPrefix(argument, option+"=")
		if value == "" {
			return "", fmt.Errorf("option %s requires a value", option)
		}
		return value, nil
	}
	if *index+1 >= len(args) {
		return "", fmt.Errorf("option %s requires a value", option)
	}
	*index++
	return args[*index], nil
}

func runCodec(args []string, encode bool, stdin io.Reader, stdout, _ io.Writer) error {
	operation := "decode"
	if encode {
		operation = "encode"
	}
	codecUsage := fmt.Sprintf("Usage: pryl %s <base64|hex|url> [options] [value]\n", operation)
	if len(args) == 0 {
		return usageError(fmt.Sprintf("%s requires an encoding", operation), codecUsage)
	}
	if len(args) == 1 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		_, err := io.WriteString(stdout, codecUsage)
		return err
	}
	if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		return usageError("help does not accept arguments", codecUsage)
	}

	switch args[0] {
	case "base64":
		return runBase64(args[1:], encode, stdin, stdout)
	case "hex", "url":
		return runSimpleCodec(args[0], args[1:], encode, stdin, stdout)
	default:
		return usageError(fmt.Sprintf("encoding %q is not implemented", args[0]), codecUsage)
	}
}

func runSimpleCodec(name string, args []string, encode bool, stdin io.Reader, stdout io.Writer) error {
	if len(args) == 1 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		_, err := io.WriteString(stdout, simpleCodecUsage)
		return err
	}
	if len(args) > 0 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		return usageError("help does not accept arguments", simpleCodecUsage)
	}
	if len(args) > 1 {
		return usageError(fmt.Sprintf("%s accepts at most one value", name), simpleCodecUsage)
	}
	if len(args) == 1 && strings.HasPrefix(args[0], "-") {
		return usageError(fmt.Sprintf("unknown option %q", args[0]), simpleCodecUsage)
	}

	input := stdin
	if len(args) == 1 {
		input = strings.NewReader(args[0])
	}

	switch name {
	case "hex":
		if encode {
			if err := hexutil.Encode(input, stdout); err != nil {
				return fmt.Errorf("encode hex: %w", err)
			}
			return nil
		}
		var decoded bytes.Buffer
		if err := hexutil.Decode(input, &decoded); err != nil {
			return usageError(err.Error(), simpleCodecUsage)
		}
		_, err := decoded.WriteTo(stdout)
		return err
	case "url":
		if encode {
			if err := urlutil.Encode(input, stdout); err != nil {
				return fmt.Errorf("encode URL: %w", err)
			}
			return nil
		}
		var decoded bytes.Buffer
		if err := urlutil.Decode(input, &decoded); err != nil {
			return usageError(err.Error(), simpleCodecUsage)
		}
		_, err := decoded.WriteTo(stdout)
		return err
	default:
		return fmt.Errorf("encoding %q is not implemented", name)
	}
}

func runBase64(args []string, encode bool, stdin io.Reader, stdout io.Writer) error {
	if len(args) == 1 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		_, err := io.WriteString(stdout, base64UsageText)
		return err
	}
	if len(args) > 0 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		return usageError("help does not accept arguments", base64UsageText)
	}

	encodingName := "standard"
	values := make([]string, 0, 1)
	for index := 0; index < len(args); index++ {
		switch {
		case args[index] == "--help" || args[index] == "-h":
			return usageError("help does not accept arguments", base64UsageText)
		case args[index] == "--encoding" || strings.HasPrefix(args[index], "--encoding="):
			value, err := optionValue(args, &index, "--encoding")
			if err != nil {
				return usageError(err.Error(), base64UsageText)
			}
			encodingName = strings.ToLower(value)
		case strings.HasPrefix(args[index], "-"):
			return usageError(fmt.Sprintf("unknown option %q", args[index]), base64UsageText)
		default:
			values = append(values, args[index])
		}
	}
	if len(values) > 1 {
		return usageError("Base64 accepts at most one value", base64UsageText)
	}

	encoding := base64util.Encoding(encodingName)
	if encoding != base64util.Standard && encoding != base64util.Raw && encoding != base64util.URL && encoding != base64util.RawURL {
		return usageError(fmt.Sprintf("Base64 encoding option %q is not implemented; choose standard, raw, url, or rawurl", encodingName), base64UsageText)
	}
	input := stdin
	if len(values) == 1 {
		input = strings.NewReader(values[0])
	}
	if encode {
		if err := base64util.Encode(input, stdout, encoding); err != nil {
			return fmt.Errorf("encode Base64: %w", err)
		}
		return nil
	}
	var decoded bytes.Buffer
	if err := base64util.Decode(input, &decoded, encoding); err != nil {
		return usageError(err.Error(), base64UsageText)
	}
	_, err := decoded.WriteTo(stdout)
	return err
}

func runTime(args []string, stdout, _ io.Writer) error {
	if len(args) == 0 {
		_, err := io.WriteString(stdout, timeUsageText)
		return err
	}
	if len(args) == 1 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		_, err := io.WriteString(stdout, timeUsageText)
		return err
	}
	if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		return usageError("help does not accept arguments", timeUsageText)
	}
	if args[0] != "epoch" {
		return usageError(fmt.Sprintf("unknown time command %q", args[0]), timeUsageText)
	}

	unitName := "seconds"
	values := make([]string, 0, 1)
	for index := 1; index < len(args); index++ {
		switch {
		case args[index] == "--help" || args[index] == "-h":
			return usageError("help does not accept arguments", timeUsageText)
		case args[index] == "--unit" || strings.HasPrefix(args[index], "--unit="):
			value, err := optionValue(args, &index, "--unit")
			if err != nil {
				return usageError(err.Error(), timeUsageText)
			}
			unitName = strings.ToLower(value)
		case strings.HasPrefix(args[index], "-") && args[index] != "-":
			return usageError(fmt.Sprintf("unknown option %q", args[index]), timeUsageText)
		default:
			values = append(values, args[index])
		}
	}
	if len(values) != 1 {
		return usageError("time epoch requires exactly one value", timeUsageText)
	}

	value, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		return usageError(fmt.Sprintf("invalid epoch value %q: expected an integer", values[0]), timeUsageText)
	}

	var unit timeutil.Unit
	switch unitName {
	case "seconds", "second", "s":
		unit = timeutil.Seconds
	case "milliseconds", "millisecond", "ms":
		unit = timeutil.Milliseconds
	default:
		return usageError(fmt.Sprintf("invalid time unit %q: expected seconds or milliseconds", unitName), timeUsageText)
	}

	result := timeutil.EpochToISO(value, unit)
	_, err = fmt.Fprintln(stdout, result)
	return err
}

func runSecret(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		_, err := io.WriteString(stdout, secretUsageText)
		return err
	}
	if len(args) == 1 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		_, err := io.WriteString(stdout, secretUsageText)
		return err
	}
	if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		return usageError("help does not accept arguments", secretUsageText)
	}
	if args[0] != "generate" {
		return usageError(fmt.Sprintf("unknown secret command %q", args[0]), secretUsageText)
	}

	length := 32
	alphabet := "alphanumeric"
	printSecret := false
	for index := 1; index < len(args); index++ {
		switch {
		case args[index] == "--help" || args[index] == "-h":
			return usageError("help does not accept arguments", secretUsageText)
		case args[index] == "--print":
			printSecret = true
		case args[index] == "--length" || strings.HasPrefix(args[index], "--length="):
			value, err := optionValue(args, &index, "--length")
			if err != nil {
				return usageError(err.Error(), secretUsageText)
			}
			length, err = strconv.Atoi(value)
			if err != nil {
				return usageError(fmt.Sprintf("invalid value for --length: %q is not an integer", value), secretUsageText)
			}
		case args[index] == "--alphabet" || strings.HasPrefix(args[index], "--alphabet="):
			value, err := optionValue(args, &index, "--alphabet")
			if err != nil {
				return usageError(err.Error(), secretUsageText)
			}
			alphabet = strings.ToLower(value)
		case strings.HasPrefix(args[index], "-"):
			return usageError(fmt.Sprintf("unknown option %q", args[index]), secretUsageText)
		default:
			return usageError(fmt.Sprintf("unexpected argument %q", args[index]), secretUsageText)
		}
	}
	if length <= 0 {
		return usageError("length must be greater than zero", secretUsageText)
	}

	result, err := secret.Generate(length, alphabet)
	if err != nil {
		return usageError(err.Error(), secretUsageText)
	}
	if printSecret {
		_, err = fmt.Fprintln(stdout, result)
		return err
	}
	if err := clipboard.Copy(result); err != nil {
		return err
	}
	_, err = fmt.Fprintln(stderr, "secret copied to clipboard")
	return err
}
