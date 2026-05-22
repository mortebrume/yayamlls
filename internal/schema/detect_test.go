package schema

import (
	"testing"

	"github.com/goccy/go-yaml/parser"
)

func TestDetectGVK(t *testing.T) {
	cases := []struct {
		name, in           string
		wantGroup, wantVer string
		wantKind           string
		wantOk             bool
	}{
		{"core pod", "apiVersion: v1\nkind: Pod\n", "", "v1", "Pod", true},
		{"grouped deployment", "apiVersion: apps/v1\nkind: Deployment\n", "apps", "v1", "Deployment", true},
		{"missing kind", "apiVersion: v1\nname: noisy\n", "", "", "", false},
		{"non-k8s doc", "name: Alice\nage: 30\n", "", "", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f, err := parser.ParseBytes([]byte(c.in), 0)
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			gvk, ok := DetectGVK(f.Docs[0].Body)
			if ok != c.wantOk {
				t.Fatalf("ok = %v, want %v", ok, c.wantOk)
			}
			if !ok {
				return
			}
			if gvk.Group != c.wantGroup || gvk.Version != c.wantVer || gvk.Kind != c.wantKind {
				t.Errorf("got %+v, want {%q %q %q}", gvk, c.wantGroup, c.wantVer, c.wantKind)
			}
		})
	}
}
