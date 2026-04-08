# Go Table Driven Test t.Run for Visual Studio Code

A CLI tool and VS Code extension that statically analyzes Go table-driven tests to extract subtest names, enabling you to run or debug individual subtests.

## Features

- Statically extracts subtest names from `t.Run()` calls by parsing Go AST
- Zero external dependencies (Go standard library only)
- VS Code extension provides CodeLens for running and debugging individual subtests

### Supported Patterns

| Pattern | Example |
|---|---|
| Slice (key-value) | `[]struct{ name string; ... }{{name: "x"}, ...}` |
| Slice (positional) | `[]struct{ name string; ... }{{"x", 1}, ...}` |
| Map | `map[string]struct{ ... }{"name": { ... }, ...}` |

## Installation

The `extension/` directory contains a VS Code extension that displays **run subtest** / **debug subtest** CodeLens for each subtest in Go test files. To use it, install both the CLI and the extension:

### 1. Install the CLI

```bash
go install github.com/otakakot/vscode-go-tdt-trun/cmd/tdtls@latest
```

### 2. Install the VS Code Extension

```bash
cd extension
pnpm install
pnpm run compile
```

Then press `F5` in VS Code to launch an Extension Development Host, or package and install manually:

```bash
# Install vsce if you don't have it
pnpm install -g @vscode/vsce

# Package and install
cd extension
vsce package
code --install-extension tdtls-*.vsix
```

## VS Code Extension

The `extension/` directory contains a VS Code extension that displays **run subtest** / **debug subtest** CodeLens for each subtest in Go test files.

### Configuration

| Setting | Default | Description |
|---|---|---|
| `tdtls.cliPath` | `tdtls` | Path to the `tdtls` CLI binary |

## Limitations

- Nested `t.Run` calls (subtests within subtests) are not supported. Only top-level `t.Run` calls inside a test function are extracted.
