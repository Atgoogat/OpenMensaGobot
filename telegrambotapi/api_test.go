package telegrambotapi

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestSplitMessageShort(t *testing.T) {
	assert.Equal(t, []string{"hello"}, splitMessage("hello", 10))
}

func TestSplitMessagePrefersBlankLines(t *testing.T) {
	parts := splitMessage("aaaa\n\nbbbb\ncc", 9)
	assert.Equal(t, []string{"aaaa", "bbbb\ncc"}, parts)
}

func TestSplitMessageRespectsLimit(t *testing.T) {
	text := strings.Repeat("Schnitzel mit Pommes €\n", 500)
	parts := splitMessage(text, MaxMessageLength)
	assert.Greater(t, len(parts), 1)
	for _, p := range parts {
		assert.LessOrEqual(t, utf8.RuneCountInString(p), MaxMessageLength)
		assert.True(t, utf8.ValidString(p))
	}
}

func TestSplitMessageNoNewlines(t *testing.T) {
	parts := splitMessage(strings.Repeat("ü", 25), 10)
	assert.Equal(t, []string{strings.Repeat("ü", 10), strings.Repeat("ü", 10), strings.Repeat("ü", 5)}, parts)
}
