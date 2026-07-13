package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"--version"}, &stdout, &stderr); err != nil {
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
	if err := Run(args, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if got, want := stdout.String(), "2024-04-05T19:34:38Z\n"; got != want {
		t.Fatalf("output = %q; want %q", got, want)
	}
}

func TestRunRejectsVersionArguments(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"version", "extra"}, &stdout, &stderr); err == nil {
		t.Fatal("expected version arguments to be rejected")
	}
}
