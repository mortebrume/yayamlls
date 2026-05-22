package lsp

import (
	"fmt"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
)

func unifiedDiff(uri string, baseline, current []byte) string {
	a := string(baseline)
	b := string(current)
	if a == b {
		return ""
	}
	edits := myers.ComputeEdits(span.URIFromURI(uri), a, b)
	if len(edits) == 0 {
		return ""
	}
	return fmt.Sprint(gotextdiff.ToUnified("baseline", "current", a, edits))
}
