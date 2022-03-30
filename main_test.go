package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/open-policy-agent/opa/rego"
)

func TestBasicWriteRead(t *testing.T) {
	f := filepath.Join(t.TempDir(), "test.txt")
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

func TestRegoPolicyLoadAbsolutePath(t *testing.T) {
	policy := `
		package signature

		allow {
			input.predicateType == "https://slsa.dev/provenance/v0.2"
		}
	`
	policyFile := filepath.Join(t.TempDir(), "policy.rego")
	if err := os.WriteFile(policyFile, []byte(policy), 0644); err != nil {
		t.Fatal(err)
	}
	r := rego.New(rego.Query("data.signature.allow"), rego.Load([]string{policyFile}, nil))

	ctx := context.Background()

	// This fails on Windows Server 2022.
	_, err := r.PrepareForEval(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegoPolicyLoadRelativePath(t *testing.T) {
	policy := `
		package signature

		allow {
			input.predicateType == "https://slsa.dev/provenance/v0.2"
		}
	`
	policyFile := "policy.rego"
	if err := os.WriteFile(policyFile, []byte(policy), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(policyFile)
	r := rego.New(rego.Query("data.signature.allow"), rego.Load([]string{policyFile}, nil))

	ctx := context.Background()

	_, err := r.PrepareForEval(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
