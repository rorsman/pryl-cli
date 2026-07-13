package secret

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	value, err := Generate(64, "hex")
	if err != nil {
		t.Fatal(err)
	}
	if len(value) != 64 {
		t.Fatalf("got length %d, want 64", len(value))
	}
	for _, char := range value {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
			t.Fatalf("unexpected character %q", char)
		}
	}
}

func TestGenerateBase64(t *testing.T) {
	value, err := Generate(32, "base64")
	if err != nil {
		t.Fatal(err)
	}
	if len(value) != 32 {
		t.Fatalf("got length %d, want 32", len(value))
	}
	if strings.ContainsAny(value, "+/=") {
		t.Fatalf("value is not raw URL-safe base64: %q", value)
	}
	if _, err := base64.RawURLEncoding.DecodeString(value); err != nil {
		t.Fatalf("generated value is not valid base64: %v", err)
	}
}

func TestGenerateRejectsExcessiveLength(t *testing.T) {
	if _, err := Generate(maxLength+1, "hex"); err == nil {
		t.Fatal("expected excessive length to be rejected")
	}
}

func TestGenerateRejectsInvalidInput(t *testing.T) {
	for _, test := range []struct {
		name, alphabet string
		length         int
	}{
		{name: "zero length", length: 0, alphabet: "hex"},
		{name: "negative length", length: -1, alphabet: "hex"},
		{name: "unknown alphabet", length: 16, alphabet: "unknown"},
		{name: "invalid base64 length", length: 5, alphabet: "base64url"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, err := Generate(test.length, test.alphabet); err == nil {
				t.Fatal("expected invalid input to be rejected")
			}
		})
	}
}
