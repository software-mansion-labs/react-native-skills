package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Rule is a function that inspects a SkillFile and returns diagnostics.
type Rule func(*SkillFile) []Diagnostic

func runAllRules(skill *SkillFile) []Diagnostic {
	var rules []Rule

	if skill.IsSkillMD {
		if !skill.HasFrontmatter {
			return []Diagnostic{{
				File: skill.Path, Severity: SeverityError,
				Rule: "frontmatter-required", Message: "SKILL.md must have YAML frontmatter (start with ---)",
			}}
		}
		// SKILL.md files require name and description.
		rules = append(rules,
			ruleName,
			ruleNameMaxLength,
			ruleNameCharset,
			ruleNameNoXML,
			ruleNameNoReservedWords,
			ruleDescription,
			ruleDescriptionMaxLength,
			ruleDescriptionNoXML,
			ruleDescriptionThirdPerson,
			ruleDescriptionIncludesTrigger,
			ruleVagueName,
		)
	} else if skill.HasFrontmatter {
		// Reference files with frontmatter: validate fields only if present.
		if skill.Frontmatter.HasName {
			rules = append(rules,
				ruleNameMaxLength,
				ruleNameCharset,
				ruleNameNoXML,
				ruleNameNoReservedWords,
				ruleVagueName,
			)
		}
		if skill.Frontmatter.HasDescription {
			rules = append(rules,
				ruleDescriptionMaxLength,
				ruleDescriptionNoXML,
				ruleDescriptionThirdPerson,
			)
		}
	}

	// Body rules apply to all .md files.
	rules = append(rules,
		ruleBodyMaxLines,
		ruleNoWindowsPaths,
		ruleNoTimeSensitive,
		ruleDeepReferences,
	)

	var all []Diagnostic
	for _, r := range rules {
		all = append(all, r(skill)...)
	}
	return all
}

// --- Name rules ---

func ruleName(s *SkillFile) []Diagnostic {
	if !s.Frontmatter.HasName {
		return []Diagnostic{{
			File: s.Path, Severity: SeverityError,
			Rule: "name-required", Message: "frontmatter is missing required field 'name'",
		}}
	}
	if strings.TrimSpace(s.Frontmatter.Name) == "" {
		return []Diagnostic{{
			File: s.Path, Severity: SeverityError,
			Rule: "name-required", Message: "'name' must not be empty",
		}}
	}
	return nil
}

func ruleNameMaxLength(s *SkillFile) []Diagnostic {
	if len(s.Frontmatter.Name) > 64 {
		return []Diagnostic{{
			File: s.Path, Severity: SeverityError,
			Rule: "name-max-length",
			Message: fmt.Sprintf("'name' is %d characters (max 64)", len(s.Frontmatter.Name)),
		}}
	}
	return nil
}

var nameCharsetRe = regexp.MustCompile(`^[a-z0-9-]+$`)

func ruleNameCharset(s *SkillFile) []Diagnostic {
	name := s.Frontmatter.Name
	if name == "" {
		return nil
	}
	if !nameCharsetRe.MatchString(name) {
		return []Diagnostic{{
			File: s.Path, Severity: SeverityError,
			Rule: "name-charset",
			Message: "'name' must contain only lowercase letters, numbers, and hyphens",
		}}
	}
	return nil
}

var xmlTagRe = regexp.MustCompile(`<[^>]+>`)

func ruleNameNoXML(s *SkillFile) []Diagnostic {
	if xmlTagRe.MatchString(s.Frontmatter.Name) {
		return []Diagnostic{{
			File: s.Path, Severity: SeverityError,
			Rule: "name-no-xml", Message: "'name' must not contain XML tags",
		}}
	}
	return nil
}

var reservedWords = []string{"anthropic", "claude"}

func ruleNameNoReservedWords(s *SkillFile) []Diagnostic {
	name := strings.ToLower(s.Frontmatter.Name)
	for _, word := range reservedWords {
		if strings.Contains(name, word) {
			return []Diagnostic{{
				File: s.Path, Severity: SeverityError,
				Rule: "name-no-reserved",
				Message: fmt.Sprintf("'name' must not contain reserved word %q", word),
			}}
		}
	}
	return nil
}

// --- Description rules ---

func ruleDescription(s *SkillFile) []Diagnostic {
	if !s.Frontmatter.HasDescription {
		return []Diagnostic{{
			File: s.Path, Severity: SeverityError,
			Rule: "description-required", Message: "frontmatter is missing required field 'description'",
		}}
	}
	if strings.TrimSpace(s.Frontmatter.Description) == "" {
		return []Diagnostic{{
			File: s.Path, Severity: SeverityError,
			Rule: "description-required", Message: "'description' must not be empty",
		}}
	}
	return nil
}

func ruleDescriptionMaxLength(s *SkillFile) []Diagnostic {
	if len(s.Frontmatter.Description) > 1024 {
		return []Diagnostic{{
			File: s.Path, Severity: SeverityError,
			Rule: "description-max-length",
			Message: fmt.Sprintf("'description' is %d characters (max 1024)", len(s.Frontmatter.Description)),
		}}
	}
	return nil
}

func ruleDescriptionNoXML(s *SkillFile) []Diagnostic {
	if xmlTagRe.MatchString(s.Frontmatter.Description) {
		return []Diagnostic{{
			File: s.Path, Severity: SeverityError,
			Rule: "description-no-xml", Message: "'description' must not contain XML tags",
		}}
	}
	return nil
}

// thirdPersonBadStarts catches first-person and second-person description openings.
var thirdPersonBadStarts = []string{
	"i can ", "i will ", "i help ", "i am ",
	"you can ", "you will ", "you should ",
	"we can ", "we will ", "we help ",
	"this helps you ", "use this to ",
}

func ruleDescriptionThirdPerson(s *SkillFile) []Diagnostic {
	desc := strings.ToLower(strings.TrimSpace(s.Frontmatter.Description))
	if desc == "" {
		return nil
	}
	for _, prefix := range thirdPersonBadStarts {
		if strings.HasPrefix(desc, prefix) {
			return []Diagnostic{{
				File: s.Path, Severity: SeverityWarning,
				Rule: "description-third-person",
				Message: fmt.Sprintf(
					"description should be written in third person; starts with %q",
					strings.TrimSpace(prefix),
				),
			}}
		}
	}
	return nil
}

// triggerPhrases are indicators that the description explains when to activate the skill.
var triggerPhrases = []string{
	"use when", "trigger on", "trigger when",
	"use for", "use this when",
	"when the user", "activate when",
}

func ruleDescriptionIncludesTrigger(s *SkillFile) []Diagnostic {
	desc := strings.ToLower(s.Frontmatter.Description)
	if desc == "" {
		return nil
	}
	for _, phrase := range triggerPhrases {
		if strings.Contains(desc, phrase) {
			return nil
		}
	}
	return []Diagnostic{{
		File: s.Path, Severity: SeverityWarning,
		Rule: "description-trigger",
		Message: "description should explain when to use this Skill (e.g. 'Use when...', 'Trigger on:')",
	}}
}

// --- Body rules ---

const maxBodyLines = 500

func ruleBodyMaxLines(s *SkillFile) []Diagnostic {
	count := len(s.BodyLines)
	if count > maxBodyLines {
		return []Diagnostic{{
			File: s.Path, Severity: SeverityWarning,
			Rule: "body-max-lines",
			Message: fmt.Sprintf("body is %d lines (recommended max %d); split into reference files", count, maxBodyLines),
		}}
	}
	return nil
}

// windowsPathRe matches backslash-separated path segments like `foo\bar.md` or `scripts\helper.py`.
// Excludes common non-path uses: escaped chars (\n, \t, \\), regex patterns, and line continuations.
var windowsPathRe = regexp.MustCompile(`[a-zA-Z0-9_-]+\\[a-zA-Z0-9_-]+\.[a-zA-Z]{1,5}`)

func ruleNoWindowsPaths(s *SkillFile) []Diagnostic {
	var diags []Diagnostic
	for i, line := range s.BodyLines {
		// Skip code blocks.
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			continue
		}
		if windowsPathRe.MatchString(line) {
			diags = append(diags, Diagnostic{
				File: s.Path, Line: i + bodyLineOffset(s), Severity: SeverityWarning,
				Rule: "no-windows-paths",
				Message: fmt.Sprintf("possible Windows-style path on line %d; use forward slashes", i+1),
			})
		}
	}
	return diags
}

// timeSensitiveRe matches phrases like "before August 2025", "after March 2026",
// "until Q3 2025", "starting January 2025".
var timeSensitiveRe = regexp.MustCompile(
	`(?i)(before|after|until|starting|by)\s+(January|February|March|April|May|June|July|August|September|October|November|December|Q[1-4])\s+\d{4}`,
)

func ruleNoTimeSensitive(s *SkillFile) []Diagnostic {
	var diags []Diagnostic
	for i, line := range s.BodyLines {
		if timeSensitiveRe.MatchString(line) {
			diags = append(diags, Diagnostic{
				File: s.Path, Line: i + bodyLineOffset(s), Severity: SeverityWarning,
				Rule: "no-time-sensitive",
				Message: fmt.Sprintf("time-sensitive language on line %d; consider using an 'old patterns' section instead", i+1),
			})
		}
	}
	return diags
}

// mdLinkRe captures markdown links: [text](path)
var mdLinkRe = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)

// ruleDeepReferences checks that referenced .md files do not themselves reference further .md files,
// which would create deeply nested references that Claude may only partially read.
func ruleDeepReferences(s *SkillFile) []Diagnostic {
	var diags []Diagnostic
	dir := filepath.Dir(s.Path)

	for _, match := range mdLinkRe.FindAllStringSubmatch(s.Body, -1) {
		target := match[2]
		// Only check relative .md file references.
		if strings.HasPrefix(target, "http") || !strings.HasSuffix(strings.ToLower(target), ".md") {
			continue
		}

		refPath := filepath.Join(dir, target)
		data, err := readFileIfExists(refPath)
		if err != nil || data == "" {
			continue
		}

		// Check if the referenced file itself contains .md links.
		subLinks := mdLinkRe.FindAllStringSubmatch(data, -1)
		for _, sub := range subLinks {
			subTarget := sub[2]
			if !strings.HasPrefix(subTarget, "http") && strings.HasSuffix(strings.ToLower(subTarget), ".md") {
				diags = append(diags, Diagnostic{
					File: s.Path, Severity: SeverityWarning,
					Rule: "shallow-references",
					Message: fmt.Sprintf(
						"referenced file %q links to %q; keep references one level deep from SKILL.md",
						target, subTarget,
					),
				})
			}
		}
	}
	return diags
}

var vagueNames = []string{"helper", "helpers", "utils", "utility", "utilities", "tools", "misc", "common"}

func ruleVagueName(s *SkillFile) []Diagnostic {
	name := strings.ToLower(s.Frontmatter.Name)
	for _, v := range vagueNames {
		if name == v {
			return []Diagnostic{{
				File: s.Path, Severity: SeverityWarning,
				Rule: "name-not-vague",
				Message: fmt.Sprintf("'name' %q is too vague; use a descriptive name", name),
			}}
		}
	}
	return nil
}

// --- Helpers ---

// bodyLineOffset returns the line number offset for body lines (accounting for frontmatter).
func bodyLineOffset(s *SkillFile) int {
	if !s.HasFrontmatter {
		return 1
	}
	idx := strings.Index(s.Raw, s.Body)
	if idx < 0 {
		return 1
	}
	return strings.Count(s.Raw[:idx], "\n") + 1
}

func readFileIfExists(path string) (string, error) {
	data, err := readFileBytes(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func readFileBytes(path string) ([]byte, error) {
	info, err := statFile(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("is a directory")
	}
	return readFile(path)
}

// Thin wrappers to make testing easier if needed.
var statFile = os.Stat
var readFile = os.ReadFile
