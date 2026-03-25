package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// SkillCreatorGrading matches the skill-creator grading.json format.
type SkillCreatorGrading struct {
	Expectations []SCExpectation `json:"expectations"`
	Summary      SCSummary       `json:"summary"`
}

type SCExpectation struct {
	Text     string `json:"text"`
	Passed   bool   `json:"passed"`
	Evidence string `json:"evidence"`
}

type SCSummary struct {
	Passed   int     `json:"passed"`
	Failed   int     `json:"failed"`
	Total    int     `json:"total"`
	PassRate float64 `json:"pass_rate"`
}

// gradeWorkspace walks a skill-creator workspace directory and applies
// static assertions from evals.json to the outputs.
func gradeWorkspace(workspacePath string, suite *EvalSuite) {
	// Find the iteration directory. The user might point to the workspace root
	// or directly to an iteration directory.
	iterDir := resolveIterationDir(workspacePath)
	if iterDir == "" {
		fatal("no iteration directory found in %s", workspacePath)
	}

	fmt.Println("============================================")
	fmt.Printf("  Static grading: %s\n", iterDir)
	fmt.Println("============================================")

	// Build assertion lookup: eval ID -> assertions.
	assertionMap := buildAssertionMap(suite)
	triggerMap := buildTriggerMap(suite)

	// Print triggering overview.
	var triggerGroup TriggerGroup
	evalDirs := findEvalDirs(iterDir)
	fmt.Printf("\n  --- TRIGGERING ---\n\n")
	for _, evalDir := range evalDirs {
		evalID := extractEvalID(filepath.Base(evalDir))
		shouldTrigger, known := triggerMap[evalID]
		if !known {
			continue
		}
		triggerGroup.Total++
		if shouldTrigger {
			triggerGroup.ShouldTrigger++
			fmt.Printf("  [TRIGGER]     eval-%d: should trigger\n", evalID)
		} else {
			triggerGroup.ShouldNotTrigger++
			fmt.Printf("  [NO TRIGGER]  eval-%d: should NOT trigger\n", evalID)
		}
	}

	var withSkill, withoutSkill GradingGroup
	configs := []struct {
		label string
		dir   string
		group *GradingGroup
	}{
		{"WITH SKILL", "with_skill", &withSkill},
		{"WITHOUT SKILL", "without_skill", &withoutSkill},
	}

	for _, cfg := range configs {
		fmt.Printf("\n  --- %s (assertions) ---\n\n", cfg.label)

		for _, evalDir := range evalDirs {
			evalID := extractEvalID(filepath.Base(evalDir))

			// Skip assertion grading for evals that should not trigger.
			if shouldTrigger, ok := triggerMap[evalID]; ok && !shouldTrigger {
				continue
			}

			assertions, ok := assertionMap[evalID]
			if !ok || len(assertions) == 0 {
				continue
			}

			cfgDir := filepath.Join(evalDir, cfg.dir)

			outputText := readAllOutputs(filepath.Join(cfgDir, "outputs"))
			if outputText == "" {
				fmt.Printf("  [SKIP] eval-%d/%s: no output files\n", evalID, cfg.dir)
				continue
			}

			var expectations []SCExpectation
			for _, a := range assertions {
				result := gradeAssertion(a, outputText, cfgDir)
				exp := SCExpectation{
					Text:     result.Text,
					Passed:   result.Passed,
					Evidence: result.Evidence,
				}
				expectations = append(expectations, exp)

				cfg.group.TotalAssertions++
				icon := "FAIL"
				if result.Passed {
					cfg.group.Passed++
					icon = "PASS"
				} else {
					cfg.group.Failed++
				}
				fmt.Printf("  [%s] eval-%d/%s: %s\n", icon, evalID, cfg.dir, result.Text)
				fmt.Printf("         %s\n", result.Evidence)
			}

			writeStaticGrading(cfgDir, expectations)
		}
	}

	if withSkill.TotalAssertions > 0 {
		withSkill.PassRate = float64(withSkill.Passed) * 100 / float64(withSkill.TotalAssertions)
	}
	if withoutSkill.TotalAssertions > 0 {
		withoutSkill.PassRate = float64(withoutSkill.Passed) * 100 / float64(withoutSkill.TotalAssertions)
	}

	fmt.Println()
	fmt.Println("============================================")
	fmt.Printf("  Triggering:    %d should-trigger, %d should-not-trigger (%d total)\n",
		triggerGroup.ShouldTrigger, triggerGroup.ShouldNotTrigger, triggerGroup.Total)
	fmt.Println("  ---")
	fmt.Printf("  With skill:    %d/%d passed (%.1f%%)\n", withSkill.Passed, withSkill.TotalAssertions, withSkill.PassRate)
	fmt.Printf("  Without skill: %d/%d passed (%.1f%%)\n", withoutSkill.Passed, withoutSkill.TotalAssertions, withoutSkill.PassRate)
	diff := withSkill.PassRate - withoutSkill.PassRate
	if diff > 0 {
		fmt.Printf("  Skill lift:    +%.1f%%\n", diff)
	} else if diff < 0 {
		fmt.Printf("  Skill lift:    %.1f%%\n", diff)
	} else {
		fmt.Printf("  Skill lift:    0%%\n")
	}
	fmt.Println("============================================")

	summary := GradingSummary{
		WithSkill:    withSkill,
		WithoutSkill: withoutSkill,
		Triggering:   triggerGroup,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	}
	summaryBytes, _ := json.MarshalIndent(summary, "", "  ")
	outPath := filepath.Join(iterDir, "static_summary.json")
	_ = os.WriteFile(outPath, summaryBytes, 0o644)
	fmt.Printf("\nStatic summary written to: %s\n", outPath)
}

// resolveIterationDir finds the iteration directory to grade.
// Accepts either a workspace root (picks latest iteration) or an iteration dir directly.
func resolveIterationDir(path string) string {
	// Check if path itself contains eval-* dirs (it's an iteration dir).
	if hasEvalDirs(path) {
		return path
	}

	// Look for iteration-N subdirectories.
	entries, err := os.ReadDir(path)
	if err != nil {
		return ""
	}

	var iterDirs []string
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), "iteration-") {
			iterDirs = append(iterDirs, filepath.Join(path, e.Name()))
		}
	}

	if len(iterDirs) == 0 {
		return ""
	}

	// Sort and pick the latest.
	sort.Strings(iterDirs)
	return iterDirs[len(iterDirs)-1]
}

func hasEvalDirs(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), "eval-") {
			return true
		}
	}
	return false
}

func findEvalDirs(iterDir string) []string {
	entries, err := os.ReadDir(iterDir)
	if err != nil {
		return nil
	}
	var dirs []string
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), "eval-") {
			dirs = append(dirs, filepath.Join(iterDir, e.Name()))
		}
	}
	sort.Strings(dirs)
	return dirs
}


func extractEvalID(dirName string) int {
	// eval-0, eval-1, eval-2, etc.
	parts := strings.SplitN(dirName, "-", 2)
	if len(parts) != 2 {
		return -1
	}
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return -1
	}
	return id
}

// buildAssertionMap creates a map from eval ID to assertions.
// Skill-creator uses 0-indexed eval IDs, evals.json uses 1-indexed.
// We store both mappings: original ID and ID-1 (0-indexed).
func buildAssertionMap(suite *EvalSuite) map[int][]Assertion {
	m := make(map[int][]Assertion)
	for _, skill := range suite.Skills {
		for i, eval := range skill.Evals {
			if len(eval.Assertions) > 0 {
				// Map by original ID from evals.json.
				m[eval.ID] = eval.Assertions
				// Also map by 0-based index for skill-creator compatibility.
				m[i] = eval.Assertions
			}
		}
	}
	return m
}

// buildTriggerMap creates a map from eval ID to should_trigger value.
func buildTriggerMap(suite *EvalSuite) map[int]bool {
	m := make(map[int]bool)
	for _, skill := range suite.Skills {
		for i, eval := range skill.Evals {
			trigger := eval.ShouldTriggerVal()
			m[eval.ID] = trigger
			m[i] = trigger
		}
	}
	return m
}

// readAllOutputs reads and concatenates all text files in an outputs directory.
func readAllOutputs(outputsDir string) string {
	entries, err := os.ReadDir(outputsDir)
	if err != nil {
		return ""
	}

	var parts []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		fp := filepath.Join(outputsDir, e.Name())
		data, err := os.ReadFile(fp)
		if err != nil {
			continue
		}

		// Skip binary-looking files.
		content := string(data)
		if isBinary(content) {
			continue
		}

		parts = append(parts, fmt.Sprintf("--- %s ---\n%s", e.Name(), content))
	}

	return strings.Join(parts, "\n\n")
}

func isBinary(s string) bool {
	for i := 0; i < len(s) && i < 512; i++ {
		if s[i] == 0 {
			return true
		}
	}
	return false
}

func writeStaticGrading(runDir string, expectations []SCExpectation) {
	passed := 0
	for _, e := range expectations {
		if e.Passed {
			passed++
		}
	}

	grading := SkillCreatorGrading{
		Expectations: expectations,
		Summary: SCSummary{
			Passed:   passed,
			Failed:   len(expectations) - passed,
			Total:    len(expectations),
			PassRate: 0,
		},
	}
	if grading.Summary.Total > 0 {
		grading.Summary.PassRate = float64(passed) / float64(grading.Summary.Total)
	}

	data, _ := json.MarshalIndent(grading, "", "  ")
	_ = os.WriteFile(filepath.Join(runDir, "static_grading.json"), data, 0o644)
}
