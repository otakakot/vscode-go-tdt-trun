package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/otakakot/vscode-go-tdt-trun/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: tdtls <file|dir|./...>\n")
		os.Exit(1)
	}

	target := os.Args[1]

	files, err := resolveTarget(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var allSubtests []parser.SubTest

	for _, f := range files {
		subs, err := parser.ExtractSubTests(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: %s: %v\n", f, err)
			continue
		}

		allSubtests = append(allSubtests, subs...)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	if err := enc.Encode(allSubtests); err != nil {
		fmt.Fprintf(os.Stderr, "error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

func resolveTarget(target string) ([]string, error) {
	if root, ok := strings.CutSuffix(target, "/..."); ok {
		if root == "" {
			root = "."
		}

		var files []string

		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && info.Name() == "testdata" {
				return filepath.SkipDir
			}

			if !info.IsDir() && strings.HasSuffix(path, "_test.go") {
				files = append(files, path)
			}

			return nil
		})

		return files, err
	}

	info, err := os.Stat(target)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return []string{target}, nil
	}

	var files []string

	entries, err := os.ReadDir(target)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), "_test.go") {
			files = append(files, filepath.Join(target, entry.Name()))
		}
	}

	return files, nil
}
