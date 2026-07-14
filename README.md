# pryl

A small CLI for everyday developer utilities.

## Quick start

Run a prebuilt binary from any directory:

```sh
pryl time epoch 1712345678
pryl secret generate --length 32
pryl encode base64 hello
```

Go is not required to run a prebuilt binary. The binary must match the target operating system and CPU architecture.

## Commands

### Time

Convert Unix timestamps to UTC ISO 8601:

```sh
pryl time epoch 1712345678
# 2024-04-05T19:34:38Z

pryl time epoch --unit milliseconds 1712345678000
# 2024-04-05T19:34:38Z
```

The supported units are `seconds` and `milliseconds`.

### Secrets

Secrets are copied to the clipboard by default:

```sh
pryl secret generate --length 32
```

Generate specific formats or explicitly print the secret:

```sh
pryl secret generate --length 32 --alphabet hex
pryl secret generate --length 32 --alphabet base64url
pryl secret generate --length 32 --print
```

macOS uses the built-in `pbcopy` command. Linux requires `wl-copy`, `xclip`, or `xsel` for clipboard support.

### Encoding and decoding

Input can be supplied as an argument or through stdin. Output is written without an extra newline, making these commands suitable for pipelines and binary data.

#### Base64

Supported modes are `standard`, `raw`, `url`, and `rawurl`:

```sh
pryl encode base64 --encoding standard hello
# aGVsbG8=
pryl decode base64 --encoding standard aGVsbG8=
# hello

pryl encode base64 --encoding raw hello
# aGVsbG8
pryl decode base64 --encoding raw aGVsbG8
# hello

pryl encode base64 --encoding url 'hello?'
# aGVsbG8_
pryl decode base64 --encoding url aGVsbG8_
# hello?

pryl encode base64 --encoding rawurl 'hello?'
# aGVsbG8_
pryl decode base64 --encoding rawurl aGVsbG8_
# hello?
```

Options may be placed before or after the input value:

```sh
pryl encode base64 --encoding raw hello
pryl encode base64 hello --encoding raw
```

#### Hexadecimal

```sh
pryl encode hex hello
# 68656c6c6f
pryl decode hex 68656c6c6f
# hello
```

#### URL components

URL encoding uses RFC 3986 percent-encoding:

```sh
pryl encode url 'hello world/?'
# hello%20world%2F%3F
pryl decode url 'hello%20world%2F%3F'
# hello world/?
```

Read from stdin when no value is provided:

```sh
printf 'hello' | pryl encode base64 --encoding rawurl
printf 'aGVsbG8' | pryl decode base64 --encoding rawurl
```

Invalid arguments and input return a non-zero exit code with relevant usage information. Unsupported encodings are reported as not implemented.

## Installation

### Prebuilt binary

After downloading a binary matching your operating system and architecture:

```sh
sudo install -m 0755 pryl /usr/local/bin/pryl
```

Verify the installation:

```sh
pryl --version
```

### Uninstall

```sh
sudo rm /usr/local/bin/pryl
```

## Build from source

Requirements:

- Go 1.26.5 or newer
- macOS or Linux

The exact Go toolchain is specified in `go.mod`.

Build for the current system:

```sh
go build -o pryl ./cmd/pryl
```

Cross-compile for another operating system and architecture by setting `GOOS` and `GOARCH`:

```sh
GOOS=darwin GOARCH=arm64 go build -o pryl-darwin-arm64 ./cmd/pryl
GOOS=darwin GOARCH=amd64 go build -o pryl-darwin-amd64 ./cmd/pryl
GOOS=linux GOARCH=amd64 go build -o pryl-linux-amd64 ./cmd/pryl
GOOS=linux GOARCH=arm64 go build -o pryl-linux-arm64 ./cmd/pryl
```

Embed a release version during a build:

```sh
VERSION=0.0.1
go build -ldflags "-X pryl/internal/version.Value=${VERSION}" -o pryl ./cmd/pryl
```

Use a `-dev` suffix for development builds, for example `VERSION=0.0.1-dev`. GitHub release builds use tags such as `v0.0.1` and embed `0.0.1` in the binary.

## Development

Run formatting, tests, and static analysis locally:

```sh
test -z "$(gofmt -l .)"
go test ./...
go vet ./...
```

The command layer is separate from the utility packages. New functionality should generally be added under `internal/`, then registered in `internal/cli`.
