# Editor setup

`yamlls` speaks LSP 3.16 over stdio. Put the binary on `$PATH` or pass
an absolute path. Workspace config: `.yamlls.yaml` at the repo root
(see the [README](../README.md)).

## Neovim (nvim-lspconfig)

```lua
local lspconfig = require("lspconfig")
local configs = require("lspconfig.configs")

if not configs.yamlls then
  configs.yamlls = {
    default_config = {
      cmd = { "yamlls" },
      filetypes = { "yaml" },
      root_dir = lspconfig.util.find_git_ancestor,
      single_file_support = true,
    },
  }
end

lspconfig.yamlls.setup({})
```

## VSCode

Minimal extension:

```jsonc
// package.json
{
  "name": "yamlls",
  "engines": { "vscode": "^1.85.0" },
  "activationEvents": ["onLanguage:yaml"],
  "main": "./extension.js"
}
```

```js
// extension.js  (vscode-languageclient v9+)
const { LanguageClient } = require("vscode-languageclient/node");

let client;
exports.activate = function (ctx) {
  client = new LanguageClient(
    "yamlls",
    "yamlls",
    { command: "yamlls" },
    { documentSelector: [{ scheme: "file", language: "yaml" }] }
  );
  client.start();              // returns Promise<void>, fire-and-forget
  ctx.subscriptions.push(client);
};
exports.deactivate = () => client?.stop();
```

## Helix

```toml
# ~/.config/helix/languages.toml
[language-server.yamlls]
command = "yamlls"

[[language]]
name = "yaml"
language-servers = ["yamlls"]
```

## Zed

Zed's `lsp` key only accepts known language-server identifiers, so
`yamlls` as a top-level key triggers `Property yamlls is not allowed`.
Override Zed's bundled `yaml-language-server` binary instead:

```jsonc
// ~/.config/zed/settings.json
{
  "lsp": {
    "yaml-language-server": {
      "binary": {
        "ignore_system_version": true,
        "path": "yamlls"
      },
      "initialization_options": {
        "catalog": true
      }
    }
  }
}
```

## Flux rendering

Install flate:

```sh
go install github.com/home-operations/flate/cmd/flate@latest
```

Open a `HelmRelease` or `Kustomization`. Schema violations on rendered
manifests are surfaced on the source document as
`[rendered <kind>/<name> @ <jsonptr>]`. The `yamlls.showRendered`
command returns the rendered YAML; in Neovim:

```lua
vim.lsp.buf.execute_command({
  command = "yamlls.showRendered",
  arguments = { vim.uri_from_bufnr(0) },
})
```

## Debugging

```sh
yamlls --log-file /tmp/yamlls.log -v 2
```

`-v 0` is silent (default), `1` is info, `2+` is debug.
