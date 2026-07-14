package base64util

import (
	"bytes"
	"strings"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	for _, test := range []struct {
		name, encoding, encoded string
	}{
		{name: "standard", encoding: "standard", encoded: "aGVsbG8="},
		{name: "raw", encoding: "raw", encoded: "aGVsbG8"},
		{name: "url", encoding: "url", encoded: "aGVsbG8="},
		{name: "rawurl", encoding: "rawurl", encoded: "aGVsbG8"},
	} {
		t.Run(test.name, func(t *testing.T) {
			var encoded bytes.Buffer
			if err := Encode(bytes.NewBufferString("hello"), &encoded, Encoding(test.encoding)); err != nil {
				t.Fatal(err)
			}
			if encoded.String() != test.encoded {
				t.Fatalf("encoded = %q; want %q", encoded.String(), test.encoded)
			}

			var decoded bytes.Buffer
			if err := Decode(&encoded, &decoded, Encoding(test.encoding)); err != nil {
				t.Fatal(err)
			}
			if decoded.String() != "hello" {
				t.Fatalf("decoded = %q; want %q", decoded.String(), "hello")
			}
		})
	}
}

func TestDecodeRejectsInvalidInput(t *testing.T) {
	var output bytes.Buffer
	if err := Decode(bytes.NewBufferString("not valid base64!"), &output, Standard); err == nil {
		t.Fatal("expected invalid Base64 input to fail")
	}
}

func TestUnknownEncoding(t *testing.T) {
	var output bytes.Buffer
	err := Encode(bytes.NewBufferString("hello"), &output, "unknown")
	if err == nil {
		t.Fatal("expected unknown encoding to fail")
	}
	if !strings.Contains(err.Error(), `Base64 encoding option "unknown" is not implemented`) {
		t.Fatalf("error = %q; want an unimplemented-option message", err)
	}
}
