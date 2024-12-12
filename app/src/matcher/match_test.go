package matcher

import "testing"

func TestMatch(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		pattern string
		want    bool
	}{
		// Basic literal matching
		{
			name:    "simple literal match",
			line:    "hello",
			pattern: "hello",
			want:    true,
		},
		{
			name:    "simple literal no match",
			line:    "hello",
			pattern: "world",
			want:    false,
		},

		// Start/End anchors
		{
			name:    "starts with",
			line:    "hello world",
			pattern: "^hello",
			want:    true,
		},
		{
			name:    "ends with",
			line:    "hello world",
			pattern: "world$",
			want:    true,
		},

		// Character classes
		{
			name:    "character class match",
			line:    "abc123",
			pattern: "[abc]",
			want:    true,
		},
		{
			name:    "negated character class",
			line:    "xyz",
			pattern: "[^abc]",
			want:    true,
		},

		// Escape sequences
		{
			name:    "digit match",
			line:    "123",
			pattern: "\\d",
			want:    true,
		},
		{
			name:    "alphanumeric match",
			line:    "abc123",
			pattern: "\\w",
			want:    true,
		},

		// Quantifiers
		{
			name:    "zero or more",
			line:    "aaa",
			pattern: "a*",
			want:    true,
		},
		{
			name:    "one or more",
			line:    "aaa",
			pattern: "a+",
			want:    true,
		},
		{
			name:    "zero or one",
			line:    "ab",
			pattern: "a?b",
			want:    true,
		},

		// Wildcards
		{
			name:    "any character",
			line:    "abc",
			pattern: "a.c",
			want:    true,
		},

		// Alternation
		{
			name:    "alternation",
			line:    "cat",
			pattern: "(cat|dog)",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Match([]byte(tt.line), tt.pattern)
			if got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
