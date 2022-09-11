package main

import (
	"fmt"
	"os"
	"testing"
)

// unit test
func TestFilter(t *testing.T) {
	testCases := []struct {
		name     string
		file     string
		ext      string
		minSize  int64
		expected bool
	}{
		{"FilterNoExtension", "testdata/dir.log", "", 0, false},
		{"FilterExtensionMatch", "testdata/dir.log", ".log", 0, false},
		{"FilterExtensionNoMatch", "testdata/dir.log", ".sh", 0, true},
		{"FilterExtensionSizeMatch", "testdata/dir.log", ".log", 10, false},
		{"FilterExtensionSizeNoMatch", "testdata/dir.log", ".log", 30, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := os.Stat(tc.file) // get the file's stats so we can pass it to filterOut
			if err != nil {
				t.Fatal(err)
			}

			f := filterOut(tc.file, tc.ext, tc.minSize, info)
			if f != tc.expected {
				fmt.Println(tc)
				t.Errorf("Expected '%t', got '%t' instead\n", tc.expected, f)
			}
		})
	}
}
