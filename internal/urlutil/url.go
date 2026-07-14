package urlutil

import (
	"bufio"
	"fmt"
	"io"
)

const hex = "0123456789ABCDEF"

// Encode percent-encodes bytes that are not RFC 3986 unreserved characters.
func Encode(input io.Reader, output io.Writer) error {
	reader := bufio.NewReader(input)
	writer := bufio.NewWriter(output)
	buffer := make([]byte, 32*1024)

	for {
		count, err := reader.Read(buffer)
		for _, value := range buffer[:count] {
			if isUnreserved(value) {
				if err := writer.WriteByte(value); err != nil {
					return fmt.Errorf("encode URL: %w", err)
				}
				continue
			}
			if _, writeErr := fmt.Fprintf(writer, "%%%c%c", hex[value>>4], hex[value&0x0f]); writeErr != nil {
				return fmt.Errorf("encode URL: %w", writeErr)
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("encode URL: %w", err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("finish URL encoding: %w", err)
	}
	return nil
}

// Decode percent-decodes URL bytes. Non-percent-encoded bytes are preserved.
func Decode(input io.Reader, output io.Writer) error {
	reader := bufio.NewReader(input)
	for {
		value, err := reader.ReadByte()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("invalid URL input: %w", err)
		}
		if value != '%' {
			if _, err := output.Write([]byte{value}); err != nil {
				return fmt.Errorf("invalid URL input: %w", err)
			}
			continue
		}

		high, err := reader.ReadByte()
		if err != nil {
			return fmt.Errorf("invalid URL input: incomplete escape sequence")
		}
		low, err := reader.ReadByte()
		if err != nil {
			return fmt.Errorf("invalid URL input: incomplete escape sequence")
		}
		highValue, highOK := fromHex(high)
		lowValue, lowOK := fromHex(low)
		if !highOK || !lowOK {
			return fmt.Errorf("invalid URL input: invalid escape sequence %q", []byte{value, high, low})
		}
		if _, err := output.Write([]byte{highValue<<4 | lowValue}); err != nil {
			return fmt.Errorf("invalid URL input: %w", err)
		}
	}
}

func isUnreserved(value byte) bool {
	return value >= 'a' && value <= 'z' ||
		value >= 'A' && value <= 'Z' ||
		value >= '0' && value <= '9' ||
		value == '-' || value == '.' || value == '_' || value == '~'
}

func fromHex(value byte) (byte, bool) {
	switch {
	case value >= '0' && value <= '9':
		return value - '0', true
	case value >= 'a' && value <= 'f':
		return value - 'a' + 10, true
	case value >= 'A' && value <= 'F':
		return value - 'A' + 10, true
	default:
		return 0, false
	}
}
