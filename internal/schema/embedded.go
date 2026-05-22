package schema

import (
	"bytes"
	_ "embed"
	"errors"
	"io"
	"strings"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed embedded/yamlls.schema.json
var yamllsSchemaJSON []byte

// EmbeddedYamllsSchemaURL is the canonical URL for the .yamlls.yaml
// schema bundled into the binary. Served via the embedded:// loader.
const EmbeddedYamllsSchemaURL = "embedded://yamlls/yamlls.schema.json"

var installEmbeddedOnce sync.Once

func installEmbeddedLoader() {
	installEmbeddedOnce.Do(func() {
		jsonschema.Loaders["embedded"] = func(url string) (io.ReadCloser, error) {
			if url == EmbeddedYamllsSchemaURL {
				return io.NopCloser(bytes.NewReader(yamllsSchemaJSON)), nil
			}
			return nil, errors.New("unknown embedded resource: " + url)
		}
	})
}

// isYamllsConfigPath reports whether docPath looks like .yamlls.yaml.
// We match basename only so the file can live at any workspace root.
func isYamllsConfigPath(docPath string) bool {
	if docPath == "" {
		return false
	}
	// Use the last slash split to avoid pulling in filepath here.
	if i := strings.LastIndex(docPath, "/"); i >= 0 {
		docPath = docPath[i+1:]
	}
	return docPath == ".yamlls.yaml" || docPath == ".yamlls.yml"
}
