# pryl

A small, modular CLI for everyday developer utilities.

## Usage

### Requirements

To run a built `pryl` binary:

- Go is not required.
- The binary must match the target operating system and CPU architecture.
- Supported operating systems are macOS and Linux.
- macOS provides clipboard support through the built-in `pbcopy` command.
- Linux requires one of `wl-copy`, `xclip`, or `xsel` for the default secret-copy behavior.

If no clipboard command is available, secrets can still be printed explicitly with `--print`.

### Examples

Convert a Unix epoch timestamp in seconds to ISO 8601:

```sh
pryl time epoch 1712345678
```

Convert a Unix epoch timestamp in milliseconds:

```sh
pryl time epoch --unit milliseconds 1712345678000
```

Generate a secure secret and copy it to the clipboard:

```sh
pryl secret generate --length 32
```

Generate a hexadecimal secret:

```sh
pryl secret generate --length 32 --alphabet hex
```

Generate a raw URL-safe Base64 secret:

```sh
pryl secret generate --length 32 --alphabet base64url
```

Explicitly print a secret to the terminal:

```sh
pryl secret generate --length 32 --print
```

Encode and decode Base64. Input can be supplied as an argument or through stdin:

```sh
pryl encode base64 hello
# aGVsbG8=

pryl decode base64 aGVsbG8=
# hello

printf 'hello' | pryl encode base64 --encoding rawurl
printf 'aGVsbG8' | pryl decode base64 --encoding rawurl

pryl encode hex hello
# 68656c6c6f

pryl decode hex 68656c6c6f
# hello

pryl encode url 'hello world/?'
# hello%20world%2F%3F

pryl decode url 'hello%20world%2F%3F'
# hello world/?
```

Base64 supports `standard`, `raw`, `url`, and `rawurl`. Hex and URL encoding do not require an additional encoding option. Options may be placed before or after the input value:

```sh
pryl encode base64 --encoding raw hello
pryl encode base64 hello --encoding raw
```

Output is written without an extra newline, making the commands suitable for pipelines and binary data. Invalid arguments and input return a non-zero exit code and include the relevant command usage. Unsupported encodings are reported as not implemented.

Check the CLI version:

```sh
pryl --version
```

Development builds report a version with the `-dev` suffix. Release builds can override it with linker flags:

```sh
go build -ldflags "-X pryl/internal/version.Value=0.0.1" -o pryl ./cmd/pryl
```

## Build

### Requirements

- Go 1.26.5 or newer
- macOS or Linux

The exact Go toolchain is specified in `go.mod`. Go is only required to build, test, or run the project from source.

### Build for the current machine

```sh
go build -o pryl ./cmd/pryl
```

Run the resulting binary:

```sh
./pryl time epoch 1712345678
```

### Cross-compile for another target

Set `GOOS` for the operating system and `GOARCH` for the CPU architecture:

```sh
# Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o pryl-darwin-arm64 ./cmd/pryl

# Intel Mac
GOOS=darwin GOARCH=amd64 go build -o pryl-darwin-amd64 ./cmd/pryl

# Linux Intel/AMD 64-bit
GOOS=linux GOARCH=amd64 go build -o pryl-linux-amd64 ./cmd/pryl

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o pryl-linux-arm64 ./cmd/pryl
```

Inspect the current native target with:

```sh
go env GOOS GOARCH
```

The resulting binary must match the target operating system and CPU architecture. Go is not required on the target machine to run a built binary.

## Install and uninstall

### Install

Choose the version to install and build the binary with that version embedded:

```sh
VERSION=0.0.1
go build -ldflags "-X pryl/internal/version.Value=${VERSION}" -o pryl ./cmd/pryl
```

Use a `-dev` suffix for a development build, for example `VERSION=0.0.1-dev`. Then install the binary system-wide:

```sh
sudo install -m 0755 pryl /usr/local/bin/pryl
```

You can now run it from any directory:

```sh
pryl help
```

Verify the installed version:

```sh
pryl --version
```

### Uninstall

```sh
sudo rm /usr/local/bin/pryl
```

## Development

Run the test suite with:

```sh
go test ./...
```

The command layer is intentionally separate from the utility packages. New functionality should generally be added under `internal/`, then registered in `internal/cli`.
