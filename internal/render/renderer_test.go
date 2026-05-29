package render_test

import (
	"testing"

	"github.com/home-operations/yamlls/internal/render"
)

func TestMatchesKind_RequiresGroupBoundary(t *testing.T) {
	doc := func(group string) *render.SourceDocument {
		return &render.SourceDocument{Kind: "HelmRelease", APIGroup: group}
	}
	const group = "helm.toolkit.fluxcd.io"

	if !render.MatchesKind(doc("helm.toolkit.fluxcd.io/v2"), "HelmRelease", group) {
		t.Error("versioned group should match")
	}
	if !render.MatchesKind(doc("helm.toolkit.fluxcd.io"), "HelmRelease", group) {
		t.Error("bare group should match")
	}
	if render.MatchesKind(doc("helm.toolkit.fluxcd.iox/v1"), "HelmRelease", group) {
		t.Error("group sharing only a prefix must not match")
	}
}
