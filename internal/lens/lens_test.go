package lens

import "testing"

func TestLenses_HelmRelease(t *testing.T) {
	text := `---
apiVersion: helm.toolkit.fluxcd.io/v2
kind: HelmRelease
metadata:
  name: foo
`
	got := Lenses("file:///tmp/hr.yaml", text)
	if len(got) != 2 {
		t.Fatalf("expected 2 lenses (view + diff), got %d", len(got))
	}
	if got[0].Command.Command != "yamlls.showRendered" {
		t.Errorf("first lens command = %q, want yamlls.showRendered", got[0].Command.Command)
	}
	if got[1].Command.Command != "yamlls.showRenderedDiff" {
		t.Errorf("second lens command = %q", got[1].Command.Command)
	}
}

func TestLenses_NonFluxDocNoLenses(t *testing.T) {
	if got := Lenses("file:///tmp/x.yaml", "apiVersion: v1\nkind: Pod\n"); len(got) != 0 {
		t.Errorf("Pod shouldn't get Flux lenses, got %d", len(got))
	}
}
