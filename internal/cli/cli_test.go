package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"--version"}, strings.NewReader(""), &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(stdout.String(), "-dev\n") {
		t.Fatalf("version output = %q; want a development version", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("unexpected stderr output: %q", stderr.String())
	}
}

func TestRunEpochMilliseconds(t *testing.T) {
	var stdout, stderr bytes.Buffer
	args := []string{"time", "epoch", "--unit", "milliseconds", "1712345678000"}
	if err := Run(args, strings.NewReader(""), &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if got, want := stdout.String(), "2024-04-05T19:34:38Z\n"; got != want {
		t.Fatalf("output = %q; want %q", got, want)
	}
}

func TestRunRejectsVersionArguments(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"version", "extra"}, strings.NewReader(""), &stdout, &stderr); err == nil {
		t.Fatal("expected version arguments to be rejected")
	}
}

func TestRunRejectsInvalidArgumentsWithUsage(t *testing.T) {
	tests := []struct {
		name, want string
		args       []string
	}{
		{name: "unknown command", args: []string{"unknown"}, want: `unknown command "unknown"`},
		{name: "missing encoding", args: []string{"encode"}, want: "encode requires an encoding"},
		{name: "invalid length", args: []string{"secret", "generate", "--length", "nope"}, want: `"nope" is not an integer`},
		{name: "invalid epoch", args: []string{"time", "epoch", "nope"}, want: `invalid epoch value "nope"`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			err := Run(test.args, strings.NewReader(""), &stdout, &stderr)
			if err == nil {
				t.Fatal("expected invalid arguments to return an error")
			}
			if !strings.Contains(err.Error(), test.want) || !strings.Contains(err.Error(), "Usage:") {
				t.Fatalf("error = %q; want %q and usage", err, test.want)
			}
		})
	}
}

func TestRunEncodeBase64AcceptsOptionAfterValue(t *testing.T) {
	var stdout, stderr bytes.Buffer
	args := []string{"encode", "base64", "hello", "--encoding", "raw"}
	if err := Run(args, strings.NewReader(""), &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if got, want := stdout.String(), "aGVsbG8"; got != want {
		t.Fatalf("output = %q; want %q", got, want)
	}
}
