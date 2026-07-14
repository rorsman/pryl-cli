package urlutil

import (
	"bytes"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	var encoded bytes.Buffer
	if err := Encode(bytes.NewBufferString("hello world/?"), &encoded); err != nil {
		t.Fatal(err)
	}
	if got, want := encoded.String(), "hello%20world%2F%3F"; got != want {
		t.Fatalf("encoded = %q; want %q", got, want)
	}

	var decoded bytes.Buffer
	if err := Decode(&encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if got, want := decoded.String(), "hello world/?"; got != want {
		t.Fatalf("decoded = %q; want %q", got, want)
	}
}

func TestDecodeRejectsInvalidEscape(t *testing.T) {
	var output bytes.Buffer
	if err := Decode(bytes.NewBufferString("hello%2"), &output); err == nil {
		t.Fatal("expected incomplete escape sequence to fail")
	}
	if err := Decode(bytes.NewBufferString("hello%GG"), &output); err == nil {
		t.Fatal("expected invalid escape sequence to fail")
	}
}

func TestEncodePreservesUnreservedCharacters(t *testing.T) {
	var output bytes.Buffer
	if err := Encode(bytes.NewBufferString("azAZ09-._~"), &output); err != nil {
		t.Fatal(err)
	}
	if got, want := output.String(), "azAZ09-._~"; got != want {
		t.Fatalf("encoded = %q; want %q", got, want)
	}
}

func TestDecodePreservesUnescapedCharacters(t *testing.T) {
	var output bytes.Buffer
	if err := Decode(bytes.NewBufferString("a+b"), &output); err != nil {
		t.Fatal(err)
	}
	if got, want := output.String(), "a+b"; got != want {
		t.Fatalf("decoded = %q; want %q", got, want)
	}
}
