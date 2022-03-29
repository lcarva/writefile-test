package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	td := t.TempDir()
	f := filepath.Join(td, "test.txt")
	if err := os.WriteFile(f, []byte("spam"), 0644); err != nil {
		t.Fatal(err)
	}
	contents, err := read(f)
	if err != nil {
		t.Fatal(err)
	}
	if contents != "spam" {
		t.Fatalf("Expected %q, got %q", "spam", contents)
	}
}
