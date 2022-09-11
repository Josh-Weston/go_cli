package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	inputFile  = "./testdata/test1.md"
	goldenFile = "./testdata/test1.md.html"
)

// TestParseContent is a unit test
func TestParseContent(t *testing.T) {
	input, err := os.ReadFile(inputFile)
	require.NoError(t, err)

	result, err := parseContent(input, "")
	require.NoError(t, err)

	expected, err := os.ReadFile(goldenFile)
	require.NoError(t, err)

	for i, b := range result {
		if expected[i] != result[i] {
			fmt.Printf("Do not match at index: %d; expected char code %d (%s), result char code %d (%s)\n", i, expected[i], string(expected[i]), b, string(b))
		}
	}

	if !bytes.Equal(expected, result) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Error("Result content does not match golden file")
	}
}

// TestRun is an integration test
func TestRun(t *testing.T) {
	var mockStdOut bytes.Buffer
	if err := run(inputFile, "", &mockStdOut, true); err != nil {
		t.Fatal(err)
	}

	resultFile := strings.TrimSpace(mockStdOut.String())

	// compare the results file with the expected file
	result, err := os.ReadFile(resultFile)
	require.NoError(t, err)

	expected, err := os.ReadFile(goldenFile)
	require.NoError(t, err)

	if !bytes.Equal(expected, result) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Error("Result content does not match golden file")
	}
	os.Remove(resultFile)
}
