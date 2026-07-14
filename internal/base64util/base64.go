package base64util

import (
	"encoding/base64"
	"fmt"
	"io"
)

type Encoding string

const (
	Standard Encoding = "standard"
	Raw      Encoding = "raw"
	URL      Encoding = "url"
	RawURL   Encoding = "rawurl"
)

func (e Encoding) codec() (*base64.Encoding, error) {
	switch e {
	case Standard:
		return base64.StdEncoding, nil
	case Raw:
		return base64.RawStdEncoding, nil
	case URL:
		return base64.URLEncoding, nil
	case RawURL:
		return base64.RawURLEncoding, nil
	default:
		return nil, fmt.Errorf("Base64 encoding option %q is not implemented; choose standard, raw, url, or rawurl", e)
	}
}

// Encode reads bytes from input and writes Base64 to output without adding a newline.
func Encode(input io.Reader, output io.Writer, encoding Encoding) error {
	codec, err := encoding.codec()
	if err != nil {
		return err
	}

	encoder := base64.NewEncoder(codec, output)
	if _, err := io.Copy(encoder, input); err != nil {
		_ = encoder.Close()
		return fmt.Errorf("encode Base64: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("finish Base64 encoding: %w", err)
	}
	return nil
}

// Decode reads Base64 from input and writes decoded bytes to output.
func Decode(input io.Reader, output io.Writer, encoding Encoding) error {
	codec, err := encoding.codec()
	if err != nil {
		return err
	}
	codec = codec.Strict()

	decoder := base64.NewDecoder(codec, input)
	if _, err := io.Copy(output, decoder); err != nil {
		return fmt.Errorf("invalid Base64 input: %w", err)
	}
	return nil
}
