package main

import (
	"os"
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
			input:    "TP-1893, TP-1893 no duplicates",
			expected: []int{1893},
		},
		{
			input:    "Title line\n\nTP-123",
			expected: []int{123},
		},
		{
			input:    "Testing version var in build step 2.",
			expected: []int{},
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

func TestBackendInMemory(t *testing.T) {
	be, err := getBackend("inmemory", "")
	if err != nil {
		t.Fatalf("failed to get the backend: %v", err)
	}
	testbackend(be, t)
}

func TestBackendFile(t *testing.T) {
	be, err := getBackend("localfile", "localfilebackend.test")
	defer os.Remove("localfilebackend.test")
	if err != nil {
		t.Fatalf("failed to get the backend: %v", err)
	}
	testbackend(be, t)
}

func TestNoBackend(t *testing.T) {
	_, err := getBackend("", "")
	if err == nil {
		t.Fatal("expected error, but got a backend")
	}
}

func testbackend(be Backend, t *testing.T) {
	id, err := be.GetLease()
	if err != nil {
		t.Errorf("error calling GetLease: %v", err)
	}
	ok, err := be.ExtendLease(id)
	if !ok {
		t.Error("failed to extend lease")
	}
	if err != nil {
		t.Errorf("failed to extend lease: %v", err)
	}
	ok, err = be.IsProcessed("abc")
	if ok {
		t.Errorf("a hash that hasn't been processed has been marked as processed")
	}
	err = be.MarkProcessed("abc")
	if err != nil {
		t.Errorf("failed to mark hash as processed: %v", err)
	}
	ok, err = be.IsProcessed("abc")
	if !ok {
		t.Errorf("failed to confirm the past processing of a hash")
	}
	if err != nil {
		t.Errorf("failed to call IsProcessed: %v", err)
	}
	err = be.CancelLease()
	if err != nil {
		t.Errorf("failed to cancel lease: %v", err)
	}
}
