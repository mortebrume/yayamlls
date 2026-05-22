// Package lens produces textDocument/codeLens results — inline command
// triggers shown above each Flux HelmRelease/Kustomization document so
// yamlls.showRendered is discoverable without a command palette.
package lens

import (
	"strings"

	"github.com/home-operations/yamlls/internal/schema"
	"github.com/home-operations/yamlls/internal/yamlast"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

const (
	commandShowRendered     = "yamlls.showRendered"
	commandShowRenderedDiff = "yamlls.showRenderedDiff"
)

func Lenses(uri, text string) []protocol.CodeLens {
	parsed := yamlast.Parse([]byte(text))
	if parsed.File == nil {
		return nil
	}
	var out []protocol.CodeLens
	for _, doc := range parsed.File.Docs {
		if doc == nil || doc.Body == nil {
			continue
		}
		gvk, ok := schema.DetectGVK(doc.Body)
		if !ok || !isFluxRenderable(gvk) {
			continue
		}
		r := yamlast.LocateRange(doc, "")
		lensRange := protocol.Range{
			Start: protocol.Position{Line: r.Start.Line, Character: 0},
			End:   protocol.Position{Line: r.Start.Line, Character: 0},
		}
		out = append(out,
			protocol.CodeLens{
				Range: lensRange,
				Command: &protocol.Command{
					Title:     "▶ View rendered",
					Command:   commandShowRendered,
					Arguments: []any{uri},
				},
			},
			protocol.CodeLens{
				Range: lensRange,
				Command: &protocol.Command{
					Title:     "⇄ Diff rendered",
					Command:   commandShowRenderedDiff,
					Arguments: []any{uri},
				},
			},
		)
	}
	return out
}

func isFluxRenderable(gvk schema.GVK) bool {
	if gvk.Kind == "HelmRelease" && strings.HasPrefix(gvk.Group, "helm.toolkit.fluxcd.io") {
		return true
	}
	if gvk.Kind == "Kustomization" && strings.HasPrefix(gvk.Group, "kustomize.toolkit.fluxcd.io") {
		return true
	}
	return false
}
