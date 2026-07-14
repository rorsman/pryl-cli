package hexutil

import (
	"bytes"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	var encoded bytes.Buffer
	if err := Encode(bytes.NewBufferString("hello"), &encoded); err != nil {
		t.Fatal(err)
	}
	if got, want := encoded.String(), "68656c6c6f"; got != want {
		t.Fatalf("encoded = %q; want %q", got, want)
	}

	var decoded bytes.Buffer
	if err := Decode(&encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if got, want := decoded.String(), "hello"; got != want {
		t.Fatalf("decoded = %q; want %q", got, want)
	}
}

func TestDecodeRejectsInvalidInput(t *testing.T) {
	var output bytes.Buffer
	if err := Decode(bytes.NewBufferString("abc"), &output); err == nil {
		t.Fatal("expected invalid hexadecimal input to fail")
	}
}

func TestDecodeAcceptsUppercaseInput(t *testing.T) {
	var output bytes.Buffer
	if err := Decode(bytes.NewBufferString("48656C6C6F"), &output); err != nil {
		t.Fatal(err)
	}
	if got, want := output.String(), "Hello"; got != want {
		t.Fatalf("decoded = %q; want %q", got, want)
	}
}
