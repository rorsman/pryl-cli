package clipboard

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// Copy writes value to the system clipboard using a native command available
// on macOS or Linux.
func Copy(value string) error {
	command, args, err := commandForPlatform(runtime.GOOS, exec.LookPath)
	if err != nil {
		return err
	}

	cmd := exec.Command(command, args...)
	cmd.Stdin = strings.NewReader(value)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("copy to clipboard using %q: %w", command, err)
	}
	return nil
}

func commandForPlatform(goos string, lookPath func(string) (string, error)) (string, []string, error) {
	switch goos {
	case "darwin":
		return "pbcopy", nil, nil
	case "linux":
		for _, candidate := range []string{"wl-copy", "xclip", "xsel"} {
			if _, err := lookPath(candidate); err == nil {
				switch candidate {
				case "xclip":
					return candidate, []string{"-selection", "clipboard"}, nil
				case "xsel":
					return candidate, []string{"--clipboard", "--input"}, nil
				default:
					return candidate, nil, nil
				}
			}
		}
		return "", nil, fmt.Errorf("no Linux clipboard tool found; install wl-clipboard, xclip, or xsel")
	default:
		return "", nil, fmt.Errorf("clipboard is not supported on %s", goos)
	}
}
