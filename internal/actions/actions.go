// Package actions implements textDocument/codeAction quick-fixes derived
// from the structured CauseData attached to each schema diagnostic.
package actions

import (
	"encoding/json"
	"fmt"

	"github.com/home-operations/yamlls/internal/diagnostics"
	"github.com/home-operations/yamlls/internal/schema"
	"github.com/santhosh-tekuri/jsonschema/v5"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Compute returns the code actions for the diagnostics in params.Context
// that lie within params.Range. sch is the schema currently associated
// with the document (may be nil; we still try to recover enum values
// from cached schema state via sch.Resolve walks).
func Compute(uri string, sch *jsonschema.Schema, diags []protocol.Diagnostic) []protocol.CodeAction {
	var out []protocol.CodeAction
	for i := range diags {
		d := diags[i]
		data, ok := decode(d.Data)
		if !ok {
			continue
		}
		switch data.Kind {
		case "enum":
			out = append(out, enumActions(uri, sch, d, data)...)
		}
	}
	return out
}

func decode(raw any) (diagnostics.CauseData, bool) {
	if raw == nil {
		return diagnostics.CauseData{}, false
	}
	// LSP marshals Diagnostic.Data through JSON, so we may see either the
	// original Go struct or a generic map. Handle both.
	switch v := raw.(type) {
	case diagnostics.CauseData:
		return v, true
	case map[string]any:
		b, err := json.Marshal(v)
		if err != nil {
			return diagnostics.CauseData{}, false
		}
		var d diagnostics.CauseData
		if err := json.Unmarshal(b, &d); err != nil {
			return diagnostics.CauseData{}, false
		}
		return d, true
	}
	return diagnostics.CauseData{}, false
}

func enumActions(
	uri string,
	root *jsonschema.Schema,
	d protocol.Diagnostic,
	data diagnostics.CauseData,
) []protocol.CodeAction {
	if root == nil {
		return nil
	}
	target := schema.Resolve(root, data.InstanceLocation)
	if target == nil || len(target.Enum) == 0 {
		return nil
	}
	out := make([]protocol.CodeAction, 0, len(target.Enum))
	for _, v := range target.Enum {
		newText := renderEnumValue(v)
		kind := protocol.CodeActionKindQuickFix
		isPreferred := len(target.Enum) == 1
		action := protocol.CodeAction{
			Title:       fmt.Sprintf("Replace with %s", newText),
			Kind:        &kind,
			Diagnostics: []protocol.Diagnostic{d},
			IsPreferred: &isPreferred,
			Edit: &protocol.WorkspaceEdit{
				Changes: map[string][]protocol.TextEdit{
					uri: {{Range: d.Range, NewText: newText}},
				},
			},
		}
		out = append(out, action)
	}
	return out
}

func renderEnumValue(v any) string {
	switch v := v.(type) {
	case string:
		return v
	case bool, int, int64, float64, json.Number:
		return fmt.Sprintf("%v", v)
	}
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}
