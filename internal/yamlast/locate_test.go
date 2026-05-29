package yamlast

import (
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestLocateRange_SpansMultilineMapping(t *testing.T) {
	text := `spec:
  containers:
    - name: web
      image: nginx
`
	parsed := Parse([]byte(text))
	if parsed.Err != nil {
		t.Fatalf("parse: %v", parsed.Err)
	}
	docs := parsed.Docs()
	if len(docs) != 1 {
		t.Fatalf("expected 1 doc, got %d", len(docs))
	}
	// /spec/containers resolves to the sequence node; its first token
	// is the `-` on line 2.
	r := LocateRange(docs[0], "/spec/containers", text)
	if r.Start.Line != 2 {
		t.Errorf("start line = %d, want 2 (the `-` line)", r.Start.Line)
	}
	if r.End.Line < 3 {
		t.Errorf("end line = %d, want >= 3", r.End.Line)
	}
}

func TestLocateRange_ScalarHasNonZeroWidth(t *testing.T) {
	text := "age: thirty\n"
	parsed := Parse([]byte(text))
	r := LocateRange(parsed.Docs()[0], "/age", text)
	want := protocol.Position{Line: 0, Character: 5 + 6}
	if r.End != want {
		t.Errorf("end = %+v, want %+v", r.End, want)
	}
}

func TestUTF16Position_TranslatesRuneColumnsToUTF16(t *testing.T) {
	// goccy reports the value token at rune column 7 (1-based): the six
	// preceding runes are 😀,k,e,y,:,space. 😀 is two UTF-16 code units, so
	// the LSP character is 7, not the 6 a rune count would give.
	line := "😀key: value"
	if got := UTF16Position(line, 1, 7); got.Character != 7 {
		t.Errorf("character = %d, want 7", got.Character)
	}
}

func TestUTF16Len_CountsCodeUnits(t *testing.T) {
	if got := UTF16Len("😀ab"); got != 4 { // 😀=2, a=1, b=1
		t.Errorf("UTF16Len = %d, want 4", got)
	}
}
