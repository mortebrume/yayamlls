package diagnostics_test

import (
	"testing"

	"github.com/home-operations/yamlls/internal/diagnostics"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestParseSuppressions(t *testing.T) {
	cases := []struct {
		name string
		text string
		// suppressed maps a 0-based line to whether it should be suppressed.
		suppressed map[uint32]bool
	}{
		{
			name:       "trailing disable-line suppresses its own line",
			text:       "name: Alice\nage: bad  # yamlls-disable-line\n",
			suppressed: map[uint32]bool{0: false, 1: true},
		},
		{
			name:       "standalone disable-line suppresses the next line",
			text:       "name: Alice\n# yamlls-disable-line\nage: bad\n",
			suppressed: map[uint32]bool{0: false, 1: false, 2: true},
		},
		{
			name:       "block disable/enable suppresses the lines between",
			text:       "a: 1\n# yamlls-disable\nb: 2\nc: 3\n# yamlls-enable\nd: 4\n",
			suppressed: map[uint32]bool{0: false, 2: true, 3: true, 5: false},
		},
		{
			name:       "disable-file suppresses every line",
			text:       "# yamlls-disable-file\na: 1\nb: 2\n",
			suppressed: map[uint32]bool{0: true, 1: true, 2: true, 99: true},
		},
		{
			name:       "hash inside a value is not a directive",
			text:       "url: http://x#yamlls-disable-line\nbad: 1\n",
			suppressed: map[uint32]bool{0: false, 1: false},
		},
		{
			name:       "disable-line token must not be a prefix match of a value comment",
			text:       "a: 1  # just a note\nb: 2\n",
			suppressed: map[uint32]bool{0: false, 1: false},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := diagnostics.ParseSuppressions(c.text)
			for line, want := range c.suppressed {
				if got := s.Suppressed(line); got != want {
					t.Errorf("line %d: Suppressed = %v, want %v", line, got, want)
				}
			}
		})
	}
}

func TestSuppressor_Filter(t *testing.T) {
	text := "a: 1\nb: 2  # yamlls-disable-line\nc: 3\n"
	diag := func(line uint32) protocol.Diagnostic {
		return protocol.Diagnostic{Range: protocol.Range{Start: protocol.Position{Line: line}}}
	}
	in := []protocol.Diagnostic{diag(0), diag(1), diag(2)}
	got := diagnostics.ParseSuppressions(text).Filter(in)
	if len(got) != 2 {
		t.Fatalf("expected 2 diagnostics after filtering, got %d", len(got))
	}
	for _, d := range got {
		if d.Range.Start.Line == 1 {
			t.Errorf("line 1 diagnostic should have been suppressed: %+v", d)
		}
	}
}
