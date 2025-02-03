package module

import (
	"errors"
	"testing"

	"github.com/d5/tengo/v2"
	"github.com/stretchr/testify/assert"
)

func TestExtractDomain(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		err      error
	}{
		{
			name:     "Valid URL",
			input:    "https://example.com/path/to/resource",
			expected: "example.com",
			err:      nil,
		},
		{
			name:     "Invalid URL",
			input:    "://invalid-url",
			expected: "",
			err:      errors.New("invalid url"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractDomain(&tengo.String{Value: tt.input})
			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result.(*tengo.String).Value)
			}
		})
	}
}

func TestExtractPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		err      error
	}{
		{
			name:     "Valid URL with path",
			input:    "https://example.com/path/to/resource",
			expected: "/path/to/resource",
			err:      nil,
		},
		{
			name:     "Valid URL without path",
			input:    "https://example.com",
			expected: "",
			err:      nil,
		},
		{
			name:     "Invalid URL",
			input:    "://invalid-url",
			expected: "",
			err:      errors.New("invalid url"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractPath(&tengo.String{Value: tt.input})
			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result.(*tengo.String).Value)
			}
		})
	}
}

func TestExtractQuery(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string][]string
		err      error
	}{
		{
			name:  "Valid URL with query",
			input: "https://example.com/path?foo=bar&baz=qux",
			expected: map[string][]string{
				"foo": {"bar"},
				"baz": {"qux"},
			},
			err: nil,
		},
		{
			name:     "Valid URL without query",
			input:    "https://example.com/path",
			expected: map[string][]string{},
			err:      nil,
		},
		{
			name:     "Invalid URL",
			input:    "://invalid-url",
			expected: nil,
			err:      errors.New("invalid url"),
		},
		{
			name:     "Invalid query",
			input:    "https://example.com/path?invalid=%zz",
			expected: nil,
			err:      errors.New("cannot decode query"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractQuery(&tengo.String{Value: tt.input})
			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
			} else {
				assert.NoError(t, err)

				// Convert the result to a map for comparison
				resultMap := make(map[string][]string)
				for key, value := range result.(*tengo.ImmutableMap).Value {
					var values []string
					for _, val := range value.(*tengo.Array).Value {
						values = append(values, val.(*tengo.String).Value)
					}
					resultMap[key] = values
				}

				assert.Equal(t, tt.expected, resultMap)
			}
		})
	}
}
