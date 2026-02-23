package shortcode

import (
	"strings"
	"testing"
)

func TestGenerateLengthAndAlphabet(t *testing.T) {
	g := New()

	for i := 0; i < 100; i++ {
		code, err := g.Generate()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(code) != Length {
			t.Fatalf("unexpected length: %d", len(code))
		}
		for j := 0; j < len(code); j++ {
			if strings.IndexByte(Alphabet, code[j]) == -1 {
				t.Fatalf("unexpected char: %q", code[j])
			}
		}
	}
}
