package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunEncodeBase64WithValue(t *testing.T) {
	var stdout, stderr bytes.Buffer
	args := []string{"encode", "base64", "hello"}
	if err := Run(args, strings.NewReader(""), &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if got, want := stdout.String(), "aGVsbG8="; got != want {
		t.Fatalf("output = %q; want %q", got, want)
	}
}

func TestRunDecodeBase64WithStdin(t *testing.T) {
	var stdout, stderr bytes.Buffer
	args := []string{"decode", "base64"}
	if err := Run(args, strings.NewReader("aGVsbG8="), &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if got, want := stdout.String(), "hello"; got != want {
		t.Fatalf("output = %q; want %q", got, want)
	}
}

func TestRunRejectsUnimplementedEncoding(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"encode", "rot13", "hello"}, strings.NewReader(""), &stdout, &stderr); err == nil {
		t.Fatal("expected an unimplemented command to return an error")
	} else if !strings.Contains(err.Error(), `encoding "rot13" is not implemented`) {
		t.Fatalf("error = %q; want an unimplemented-command message", err)
	}
}

func TestRunEncodeHexAndURL(t *testing.T) {
	for _, test := range []struct {
		name, command, input, want string
	}{
		{name: "hex", command: "hex", input: "hello", want: "68656c6c6f"},
		{name: "url", command: "url", input: "hello world", want: "hello%20world"},
	} {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			args := []string{"encode", test.command, test.input}
			if err := Run(args, strings.NewReader(""), &stdout, &stderr); err != nil {
				t.Fatal(err)
			}
			if got := stdout.String(); got != test.want {
				t.Fatalf("output = %q; want %q", got, test.want)
			}
		})
	}
}
