package hexutil

import (
	"encoding/hex"
	"fmt"
	"io"
)

// Encode reads bytes from input and writes lowercase hexadecimal to output.
func Encode(input io.Reader, output io.Writer) error {
	encoder := hex.NewEncoder(output)
	if _, err := io.Copy(encoder, input); err != nil {
		return fmt.Errorf("encode hexadecimal: %w", err)
	}
	return nil
}

// Decode reads hexadecimal from input and writes decoded bytes to output.
func Decode(input io.Reader, output io.Writer) error {
	decoder := hex.NewDecoder(input)
	if _, err := io.Copy(output, decoder); err != nil {
		return fmt.Errorf("invalid hexadecimal input: %w", err)
	}
	return nil
}
