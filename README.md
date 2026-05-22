# yamlls

YAML language server in Go. Schema-driven diagnostics, completion, and
hover; pluggable rendering for Flux `HelmRelease` and `Kustomization`
sources via [home-operations/flate][flate].

Per-document schema resolution, highest priority first:

1. in-file modeline (`# yaml-language-server: $schema=…`)
2. workspace `schemas:` glob in `.yamlls.yaml`
3. JSON Schema Store catalog (filename match)
4. Kubernetes `apiVersion`+`kind` → `kubernetes.schemaUrl` template

Multi-doc files validate each document against its own schema. The
default `kubernetes.schemaUrl` is
`https://k8s-schemas.home-operations.com/{groupSeg}{kindLower}_{versionLower}.json`;
override in `.yamlls.yaml` to point elsewhere. 404s are silently skipped.

## Install

Homebrew:

```sh
brew install home-operations/tap/yamlls
```

Go:

```sh
go install github.com/home-operations/yamlls/cmd/yamlls@latest
```

Prebuilt binaries for linux/darwin/windows (amd64+arm64) are attached to
each [GitHub release](https://github.com/home-operations/yamlls/releases).

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

# Optional. Override the URL template for Kubernetes auto-detect.
# Placeholders: {group}, {groupSeg}, {groupFirst}, {groupFirstSeg},
# {kind}, {kindLower}, {version}, {versionLower}.
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

## Capabilities

`textDocument/`: diagnostics, completion, hover, foldingRange,
documentLink, documentSymbol, codeAction (enum quick-fix), codeLens.

`workspace/`: didChangeConfiguration, didChangeWorkspaceFolders, executeCommand.

## Commands

- `yamlls.showRendered <uri>` — rendered output for a Flux source.
- `yamlls.showRenderedDiff <uri>` — unified diff between the open-time
  render and the current render.

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
