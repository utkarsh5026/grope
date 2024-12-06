package table

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test struct for struct slice testing
type TestPerson struct {
	Name string
	Age  int
	City string
}

func TestPrintTable(t *testing.T) {
	t.Run("prints slice of structs", func(t *testing.T) {
		data := []TestPerson{
			{"Alice", 25, "New York"},
			{"Bob", 30, "San Francisco"},
		}

		err := PrintTable(data, Options{})
		assert.NoError(t, err)
	})

	t.Run("prints slice of slices", func(t *testing.T) {
		data := [][]string{
			{"Name", "Age", "City"},
			{"Alice", "25", "New York"},
			{"Bob", "30", "San Francisco"},
		}

		err := PrintTable(data, Options{})
		assert.NoError(t, err)
	})

	t.Run("prints with custom headers", func(t *testing.T) {
		data := []TestPerson{
			{"Alice", 25, "New York"},
			{"Bob", 30, "San Francisco"},
		}

		opts := Options{
			Headers: []string{"Full Name", "Years", "Location"},
		}
		err := PrintTable(data, opts)
		assert.NoError(t, err)
	})

	t.Run("returns error for invalid input", func(t *testing.T) {
		data := "not a slice"
		err := PrintTable(data, Options{})
		assert.Error(t, err)
	})

	t.Run("handles empty slice", func(t *testing.T) {
		var data []TestPerson
		err := PrintTable(data, Options{})
		assert.NoError(t, err)
	})
}

func TestParseHeaders(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "converts to title case",
			input:    []string{"first name", "last name"},
			expected: []string{"First_Name", "Last_Name"},
		},
		{
			name:     "handles empty strings",
			input:    []string{"", "header"},
			expected: []string{"", "Header"},
		},
		{
			name:     "handles mixed case",
			input:    []string{"FiRsT NaMe", "LaSt nAmE"},
			expected: []string{"First_Name", "Last_Name"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseHeaders(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "handles nil",
			input:    nil,
			expected: "",
		},
		{
			name:     "converts int",
			input:    42,
			expected: "42",
		},
		{
			name:     "converts string",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "converts bool",
			input:    true,
			expected: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
