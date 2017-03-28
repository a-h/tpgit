package main

import (
	"reflect"
	"testing"
)

func TestExtractFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected []int
	}{
		{
			input:    "TP-1893, TP-1895, TP-1904 setup TLS certificates",
			expected: []int{1893, 1895, 1904},
		},
		{
			input:    "Tp 286\nUpdated API url on production to be the same origin.",
			expected: []int{286},
		},
		{
			input:    "[TP-450] uuid replaced by human readable shortID",
			expected: []int{450},
		},
		{
			input:    "TP:450 something or other",
			expected: []int{450},
		},
		{
			input:    "Merge pull request #14 from features/TP-1931",
			expected: []int{1931},
		},
		{
			input:    "#1931 - Merging pull request #34",
			expected: []int{1931},
		},
		{
			input:    "Merge pull request #3 from features/tp-1889-remove-access",
			expected: []int{1889},
		},
		{
			input:    "TP-404: Added the payment details text content",
			expected: []int{404},
		},
		{
			// Overflow
			input:    "TP-10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			expected: []int{},
		},
	}

	for _, test := range tests {
		actual := extract(test.input)

		if !reflect.DeepEqual(test.expected, actual) {
			t.Errorf("For %s, expected %v, but got %v", test.input, test.expected, actual)
		}
	}
}
