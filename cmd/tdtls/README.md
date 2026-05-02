# tdtls

A CLI tool that statically extracts subtest names from Go table-driven tests.

## Installation

```bash
go install github.com/otakakot/vscode-go-tdt-trun/cmd/tdtls@latest
```

## Usage

```bash
# Single file
tdtls path/to/some_test.go

# All test files in a directory
tdtls path/to/dir

# Recursive search
tdtls ./...
```

### Output

Outputs subtest information as JSON to stdout.

```json
[
  {
    "func": "TestAdd",
    "name": "both positive",
    "file": "/absolute/path/to/calc_test.go",
    "line": 16
  },
  {
    "func": "TestAdd",
    "name": "positive and negative",
    "file": "/absolute/path/to/calc_test.go",
    "line": 17
  }
]
```
