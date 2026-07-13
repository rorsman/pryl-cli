package clipboard

import (
	"errors"
	"reflect"
	"testing"
)

func TestCommandForPlatform(t *testing.T) {
	tests := []struct {
		name, goos, available, wantCommand string
		wantArgs                           []string
	}{
		{name: "macOS", goos: "darwin", wantCommand: "pbcopy"},
		{name: "Wayland", goos: "linux", available: "wl-copy", wantCommand: "wl-copy"},
		{name: "X11 xclip", goos: "linux", available: "xclip", wantCommand: "xclip", wantArgs: []string{"-selection", "clipboard"}},
		{name: "X11 xsel", goos: "linux", available: "xsel", wantCommand: "xsel", wantArgs: []string{"--clipboard", "--input"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lookup := func(command string) (string, error) {
				if command == test.available {
					return "/usr/bin/" + command, nil
				}
				return "", errors.New("not found")
			}

			command, args, err := commandForPlatform(test.goos, lookup)
			if err != nil {
				t.Fatal(err)
			}
			if command != test.wantCommand || !reflect.DeepEqual(args, test.wantArgs) {
				t.Fatalf("got %q %v; want %q %v", command, args, test.wantCommand, test.wantArgs)
			}
		})
	}
}

func TestCommandForPlatformErrors(t *testing.T) {
	lookup := func(string) (string, error) { return "", errors.New("not found") }

	if _, _, err := commandForPlatform("linux", lookup); err == nil {
		t.Fatal("expected missing Linux clipboard tools to return an error")
	}
	if _, _, err := commandForPlatform("windows", lookup); err == nil {
		t.Fatal("expected unsupported platform to return an error")
	}
}
