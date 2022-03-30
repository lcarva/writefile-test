package main

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestRegoPolicy(t *testing.T) {
	policy := `
	package signature

	allow {
		input.predicateType == "https://slsa.dev/provenance/v0.2"
	}
	`
	jsonBody := []byte(`{
		"_type": "https://in-toto.io/Statement/v0.1",
		"predicateType": "https://slsa.dev/provenance/v0.2"
	}`)
	policyFile := filepath.Join(t.TempDir(), "policy.rego")
	if err := os.WriteFile(policyFile, []byte(policy), 0644); err != nil {
		t.Fatal(err)
	}
	r := rego.New(rego.Query("data.signature.allow"), rego.Load([]string{policyFile}, nil))
	if r == nil {
		t.Fatal("Expected *rego.Rego, got nil")
	}

	ctx := context.Background()

	query, err := r.PrepareForEval(ctx)
	if err != nil {
		t.Fatal(err)
	}

	var input interface{}
	dec := json.NewDecoder(bytes.NewBuffer(jsonBody))
	dec.UseNumber()
	if err := dec.Decode(&input); err != nil {
		t.Fatal(err)
	}

	rs, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		t.Fatal(err)
	}

	if !rs.Allowed() {
		t.Fatal("Not allowed")
	}

}
