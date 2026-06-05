package impl

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestSplitContentKeepsUTF8Boundaries(t *testing.T) {
	text := strings.Repeat("你", 360) + "\n" + strings.Repeat("好", 360)

	chunks := splitContent(text, 500)
	if len(chunks) < 2 {
		t.Fatalf("expected multiple chunks, got %d", len(chunks))
	}

	for i, chunk := range chunks {
		if !utf8.ValidString(chunk) {
			t.Fatalf("chunk %d is not valid UTF-8", i)
		}
	}
	if got := strings.Join(chunks, ""); got != text {
		t.Fatalf("split/join changed content")
	}
}
