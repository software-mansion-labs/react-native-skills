package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	paths := os.Args[1:]
	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: skill-lint <path> [path...]")
		fmt.Fprintln(os.Stderr, "  Each path can be a .md file or a directory to search recursively.")
		os.Exit(1)
	}

	files, err := collectSkillFiles(paths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "No .md files found.")
		os.Exit(1)
	}

	totalDiagnostics := 0
	totalErrors := 0

	for _, f := range files {
		skill, err := parseSkillFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: parse error: %v\n", f, err)
			totalErrors++
			continue
		}

		diagnostics := runAllRules(skill)
		totalDiagnostics += len(diagnostics)

		for _, d := range diagnostics {
			fmt.Println(d)
		}
	}

	fmt.Fprintf(os.Stderr, "\n%d file(s) checked, %d diagnostic(s) found\n", len(files), totalDiagnostics)

	if totalDiagnostics > 0 || totalErrors > 0 {
		os.Exit(1)
	}
}

func collectSkillFiles(paths []string) ([]string, error) {
	var files []string
	seen := make(map[string]bool)

	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			return nil, fmt.Errorf("cannot access %q: %w", p, err)
		}

		if !info.IsDir() {
			abs, _ := filepath.Abs(p)
			if !seen[abs] {
				files = append(files, abs)
				seen[abs] = true
			}
			continue
		}

		err = filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			name := info.Name()
			lower := strings.ToLower(name)
			if !strings.HasSuffix(lower, ".md") {
				return nil
			}
			// Skip common non-skill markdown files.
			if lower == "readme.md" || lower == "changelog.md" || lower == "contributing.md" || lower == "license.md" {
				return nil
			}
			abs, _ := filepath.Abs(path)
			if !seen[abs] {
				files = append(files, abs)
				seen[abs] = true
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("walking %q: %w", p, err)
		}
	}
	return files, nil
}
