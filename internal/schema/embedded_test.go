package schema

import (
	"io"
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

func TestEmbeddedYamllsSchema_LoadsAndCompiles(t *testing.T) {
	installEmbeddedLoader()
	loader := jsonschema.Loaders["embedded"]
	if loader == nil {
		t.Fatal("embedded loader not registered")
	}
	rc, err := loader(EmbeddedYamllsSchemaURL)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	body, _ := io.ReadAll(rc)
	_ = rc.Close()
	if !strings.Contains(string(body), "kubernetes") {
		t.Errorf("embedded schema body looks wrong: %s", body[:200])
	}

	c := jsonschema.NewCompiler()
	c.Draft = jsonschema.Draft2020
	sch, err := c.Compile(EmbeddedYamllsSchemaURL)
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	// Round-trip a valid .yamlls.yaml-shaped payload through the schema.
	good := map[string]any{
		"catalog": true,
		"kubernetes": map[string]any{
			"schemaUrl": "https://example.com/{kindLower}.json",
		},
	}
	if err := sch.Validate(good); err != nil {
		t.Errorf("valid config rejected: %v", err)
	}

	// Reject an unknown top-level key.
	bad := map[string]any{"schmas": map[string]any{"a": "b"}}
	if err := sch.Validate(bad); err == nil {
		t.Errorf("typo'd key 'schmas' should be rejected")
	}
}

func TestResolver_RecognizesYamllsConfigPath(t *testing.T) {
	r := NewResolver()
	if got := r.Resolve("", "/repo/.yamlls.yaml"); got != EmbeddedYamllsSchemaURL {
		t.Errorf("got %q, want %q", got, EmbeddedYamllsSchemaURL)
	}
	if got := r.Resolve("", "/repo/other.yaml"); got == EmbeddedYamllsSchemaURL {
		t.Errorf("unrelated file should NOT match the embedded schema")
	}
}
