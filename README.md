# yamlls

YAML language server in Go. Schema-driven diagnostics, completion, and
hover; pluggable rendering for Flux `HelmRelease` and `Kustomization`
sources via [home-operations/flate][flate].

Schema resolution per document — modeline →
workspace settings glob → Kubernetes `apiVersion`+`kind` auto-detect →
JSON Schema Store catalog. Multi-doc YAML files validate each document
against its own schema.

## Install

```sh
go install github.com/home-operations/yamlls/cmd/yamlls@latest
```

For Flux rendering:

```sh
go install github.com/home-operations/flate/cmd/flate@latest
```

## Editor setup

See [docs/SETUP.md](docs/SETUP.md).

## Configuration

`.yamlls.yaml` in the workspace root:

```yaml
schemas:
  "https://json.schemastore.org/github-workflow.json":
    - ".github/workflows/*.yml"
  "./schemas/local.json":
    - "k8s/**/*.yaml"

catalog: true
catalogUrl: ""

# Optional. Override the URL template used by Kubernetes auto-detect.
# Placeholders: {group}, {groupSeg}, {groupFirst}, {kind}, {kindLower},
# {version}, {versionLower}. Unset = yannh/kubernetes-json-schema layout.
# kubernetes:
#   schemaUrl: "https://schemas.example.com/{groupSeg}{kindLower}_{versionLower}.json"

# Optional. Defaults shown.
# renderers:
#   flate:
#     enabled: true
#     binary: flate
```

See [`.yamlls.yaml.example`](.yamlls.yaml.example) for a copyable starter.

Same shape works via `initializationOptions` or
`workspace/didChangeConfiguration`. Precedence (low → high):
`.yamlls.yaml` → `initializationOptions` → `didChangeConfiguration`.

## Commands

- `yamlls.showRendered` — returns the renderer's output for a
  `HelmRelease`/`Kustomization` URI.

## CLI flags

```
yamlls --version              print version and exit
yamlls --log-file PATH        append logs to PATH instead of stderr
yamlls -v N                   log verbosity (0=silent, 1=info, 2+=debug)
```

## Development

```sh
mise install   # toolchain
mise run test
mise run lint
mise run build
```

[flate]: https://github.com/home-operations/flate
