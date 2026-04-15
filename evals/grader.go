package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func gradeAssertion(a Assertion, response, runDir string) GradingResult {
	text := a.Text
	if text == "" {
		text = a.Value
	}

	switch a.Type {
	case "contains":
		if strings.Contains(strings.ToLower(response), strings.ToLower(a.Value)) {
			return GradingResult{Text: text, Passed: true, Evidence: fmt.Sprintf("Found '%s' in response", a.Value)}
		}
		return GradingResult{Text: text, Passed: false, Evidence: fmt.Sprintf("'%s' not found in response", a.Value)}

	case "contains_any":
		lowerResponse := strings.ToLower(response)
		for _, v := range a.Values {
			if strings.Contains(lowerResponse, strings.ToLower(v)) {
				return GradingResult{Text: text, Passed: true, Evidence: fmt.Sprintf("Found '%s' in response", v)}
			}
		}
		return GradingResult{Text: text, Passed: false, Evidence: fmt.Sprintf("None of %v found in response", a.Values)}

	case "not_contains":
		if strings.Contains(strings.ToLower(response), strings.ToLower(a.Value)) {
			return GradingResult{Text: text, Passed: false, Evidence: fmt.Sprintf("Found '%s' in response (should not be present)", a.Value)}
		}
		return GradingResult{Text: text, Passed: true, Evidence: fmt.Sprintf("'%s' correctly absent from response", a.Value)}

	case "file_exists":
		path := a.Value
		if !filepath.IsAbs(path) {
			path = filepath.Join(runDir, "outputs", path)
		}
		if _, err := os.Stat(path); err == nil {
			return GradingResult{Text: text, Passed: true, Evidence: fmt.Sprintf("File exists: %s", a.Value)}
		}
		return GradingResult{Text: text, Passed: false, Evidence: fmt.Sprintf("File not found: %s", a.Value)}

	case "exit_code":
		metaBytes, err := os.ReadFile(filepath.Join(runDir, "metadata.json"))
		if err != nil {
			return GradingResult{Text: text, Passed: false, Evidence: "Cannot read metadata.json"}
		}
		var meta map[string]interface{}
		if err := json.Unmarshal(metaBytes, &meta); err != nil {
			return GradingResult{Text: text, Passed: false, Evidence: "Cannot parse metadata.json"}
		}
		if code, ok := meta["exit_code"]; ok {
			actual := fmt.Sprintf("%v", code)
			if actual == a.Value {
				return GradingResult{Text: text, Passed: true, Evidence: fmt.Sprintf("Exit code matches: %s", a.Value)}
			}
			return GradingResult{Text: text, Passed: false, Evidence: fmt.Sprintf("Exit code %s, expected %s", actual, a.Value)}
		}
		return GradingResult{Text: text, Passed: false, Evidence: "No exit_code in metadata.json"}

	default:
		return GradingResult{Text: text, Passed: false, Evidence: fmt.Sprintf("Unknown assertion type: %s", a.Type)}
	}
}
