package secret

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
)

const maxLength = 1 << 20

const (
	alphanumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	hexadecimal  = "0123456789abcdef"
)

// Generate returns a cryptographically secure random string.
func Generate(length int, alphabetName string) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be greater than zero")
	}
	if length > maxLength {
		return "", fmt.Errorf("length must not exceed %d", maxLength)
	}
	if alphabetName == "base64" || alphabetName == "base64url" {
		return generateBase64(length)
	}

	alphabet, err := alphabetFor(alphabetName)
	if err != nil {
		return "", err
	}

	result := make([]byte, 0, length)
	// Discard values above the largest complete alphabet-sized range. This
	// avoids modulo bias while retaining crypto/rand as the entropy source.
	limit := 256 - (256 % len(alphabet))
	buffer := make([]byte, 128)
	for len(result) < length {
		if _, err := rand.Read(buffer); err != nil {
			return "", fmt.Errorf("read secure random data: %w", err)
		}
		for _, value := range buffer {
			if int(value) >= limit {
				continue
			}
			result = append(result, alphabet[int(value)%len(alphabet)])
			if len(result) == length {
				break
			}
		}
	}
	return string(result), nil
}

func generateBase64(length int) (string, error) {
	if length%4 == 1 {
		return "", errors.New("base64 length cannot be 1 modulo 4")
	}
	byteCount := (length*6 + 7) / 8
	random := make([]byte, byteCount)
	if _, err := rand.Read(random); err != nil {
		return "", fmt.Errorf("read secure random data: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(random)[:length], nil
}

func alphabetFor(name string) (string, error) {
	switch name {
	case "alphanumeric":
		return alphanumeric, nil
	case "hex":
		return hexadecimal, nil
	default:
		return "", fmt.Errorf("unknown alphabet %q; choose alphanumeric, hex, or base64url", name)
	}
}
