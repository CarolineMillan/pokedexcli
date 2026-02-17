package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "hello        world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "HELLO world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "hello wOrld",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Hello World",
			expected: []string{"hello", "world"},
		},
		// add more cases here
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test

		if len(c.expected) != len(actual) {
			t.Errorf("Test Failed. Expected length: %v, Actual length: %v", len(c.expected), len(actual))
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("Test Failed. Expected: %v, Actual: %v", expectedWord, word)
				break
			}
		}
	}

}
