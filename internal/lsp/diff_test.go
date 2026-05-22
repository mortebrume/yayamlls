package lsp

import (
	"strings"
	"testing"
)

func TestUnifiedDiff_ReturnsEmptyWhenEqual(t *testing.T) {
	if got := unifiedDiff("file:///x", []byte("a\n"), []byte("a\n")); got != "" {
		t.Errorf("expected empty diff, got %q", got)
	}
}

func TestUnifiedDiff_ProducesUnifiedFormat(t *testing.T) {
	got := unifiedDiff("file:///x", []byte("a\nb\nc\n"), []byte("a\nB\nc\n"))
	if !strings.Contains(got, "-b") || !strings.Contains(got, "+B") {
		t.Errorf("expected unified diff with -b/+B lines, got:\n%s", got)
	}
}

func TestUnifiedDiff_EmptyBaselineYieldsAdditions(t *testing.T) {
	got := unifiedDiff("file:///x", nil, []byte("apiVersion: v1\nkind: Pod\n"))
	if !strings.Contains(got, "+apiVersion: v1") {
		t.Errorf("expected addition lines, got:\n%s", got)
	}
}
