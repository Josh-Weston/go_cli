package main

import (
	"bytes"
	"testing"
)

func TestCountWords(t *testing.T) {

	// Can do this
	// oob := new(bytes.Buffer)
	// oob.WriteString("word1 word2 word3 word4\n")

	// Or this
	// var ob bytes.Buffer
	// ob.WriteString("word1 word2 word3 word4\n")

	// Or this
	b := bytes.NewBufferString("word1 word2 word3 word4\n") // to create the io.Reader interface
	exp := 4
	res := count(b, false)
	if res != exp {
		t.Errorf("Expected %d, got %d instead.\n", exp, res)
	}
}

func TestCountLines(t *testing.T) {
	b := bytes.NewBufferString("word1 word2 word3\nline2\nline3 word1") // to create the io.Reader interface
	exp := 3
	res := count(b, true)
	if res != exp {
		t.Errorf("Expected %d, got %d instead.\n", exp, res)
	}
}
