package main

import (
	"strings"
	"testing"
)

func makeSkill(name, description, body string) *SkillFile {
	raw := "---\nname: " + name + "\ndescription: " + description + "\n---\n" + body
	return &SkillFile{
		Path: "test/SKILL.md",
		Frontmatter: Frontmatter{
			Name:           name,
			Description:    description,
			HasName:        name != "",
			HasDescription: description != "",
		},
		HasFrontmatter: true,
		IsSkillMD:      true,
		Body:           body,
		BodyLines:      strings.Split(body, "\n"),
		Raw:            raw,
	}
}

func hasDiag(diags []Diagnostic, rule string) bool {
	for _, d := range diags {
		if d.Rule == rule {
			return true
		}
	}
	return false
}

func countDiag(diags []Diagnostic, rule string) int {
	n := 0
	for _, d := range diags {
		if d.Rule == rule {
			n++
		}
	}
	return n
}

// --- Name tests ---

func TestRuleName_Missing(t *testing.T) {
	s := makeSkill("", "", "body")
	s.Frontmatter.HasName = false
	diags := ruleName(s)
	if !hasDiag(diags, "name-required") {
		t.Error("expected name-required diagnostic")
	}
}

func TestRuleName_Empty(t *testing.T) {
	s := makeSkill("", "desc", "body")
	s.Frontmatter.HasName = true
	diags := ruleName(s)
	if !hasDiag(diags, "name-required") {
		t.Error("expected name-required for empty name")
	}
}

func TestRuleName_Valid(t *testing.T) {
	s := makeSkill("my-skill-1", "desc", "body")
	diags := ruleName(s)
	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}
}

func TestRuleNameMaxLength(t *testing.T) {
	long := strings.Repeat("a", 65)
	s := makeSkill(long, "desc", "body")
	diags := ruleNameMaxLength(s)
	if !hasDiag(diags, "name-max-length") {
		t.Error("expected name-max-length diagnostic")
	}
}

func TestRuleNameMaxLength_AtLimit(t *testing.T) {
	exact := strings.Repeat("a", 64)
	s := makeSkill(exact, "desc", "body")
	diags := ruleNameMaxLength(s)
	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics for 64-char name: %v", diags)
	}
}

func TestRuleNameCharset_Valid(t *testing.T) {
	for _, name := range []string{"pdf-processing", "my-skill-123", "a"} {
		s := makeSkill(name, "desc", "body")
		diags := ruleNameCharset(s)
		if len(diags) > 0 {
			t.Errorf("unexpected diagnostic for valid name %q: %v", name, diags)
		}
	}
}

func TestRuleNameCharset_Invalid(t *testing.T) {
	for _, name := range []string{"MySkill", "my_skill", "my skill", "pdf.processing"} {
		s := makeSkill(name, "desc", "body")
		diags := ruleNameCharset(s)
		if !hasDiag(diags, "name-charset") {
			t.Errorf("expected name-charset diagnostic for %q", name)
		}
	}
}

func TestRuleNameNoXML(t *testing.T) {
	s := makeSkill("<b>bold</b>", "desc", "body")
	diags := ruleNameNoXML(s)
	if !hasDiag(diags, "name-no-xml") {
		t.Error("expected name-no-xml diagnostic")
	}
}

func TestRuleNameNoReservedWords(t *testing.T) {
	for _, name := range []string{"anthropic-helper", "claude-tools", "my-claude-skill"} {
		s := makeSkill(name, "desc", "body")
		diags := ruleNameNoReservedWords(s)
		if !hasDiag(diags, "name-no-reserved") {
			t.Errorf("expected name-no-reserved diagnostic for %q", name)
		}
	}
}

func TestRuleNameNoReservedWords_Clean(t *testing.T) {
	s := makeSkill("pdf-processing", "desc", "body")
	diags := ruleNameNoReservedWords(s)
	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}
}

// --- Description tests ---

func TestRuleDescription_Missing(t *testing.T) {
	s := makeSkill("name", "", "body")
	s.Frontmatter.HasDescription = false
	diags := ruleDescription(s)
	if !hasDiag(diags, "description-required") {
		t.Error("expected description-required diagnostic")
	}
}

func TestRuleDescriptionMaxLength(t *testing.T) {
	long := strings.Repeat("x", 1025)
	s := makeSkill("name", long, "body")
	diags := ruleDescriptionMaxLength(s)
	if !hasDiag(diags, "description-max-length") {
		t.Error("expected description-max-length diagnostic")
	}
}

func TestRuleDescriptionNoXML(t *testing.T) {
	s := makeSkill("name", "Processes <b>files</b>", "body")
	diags := ruleDescriptionNoXML(s)
	if !hasDiag(diags, "description-no-xml") {
		t.Error("expected description-no-xml diagnostic")
	}
}

func TestRuleDescriptionThirdPerson(t *testing.T) {
	bad := []string{
		"I can help you process files",
		"You can use this to process files",
		"We help with file processing",
	}
	for _, desc := range bad {
		s := makeSkill("name", desc, "body")
		diags := ruleDescriptionThirdPerson(s)
		if !hasDiag(diags, "description-third-person") {
			t.Errorf("expected description-third-person diagnostic for %q", desc)
		}
	}
}

func TestRuleDescriptionThirdPerson_Good(t *testing.T) {
	s := makeSkill("name", "Processes PDF files and extracts text. Use when working with PDFs.", "body")
	diags := ruleDescriptionThirdPerson(s)
	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}
}

func TestRuleDescriptionIncludesTrigger(t *testing.T) {
	s := makeSkill("name", "Processes PDF files and extracts text", "body")
	diags := ruleDescriptionIncludesTrigger(s)
	if !hasDiag(diags, "description-trigger") {
		t.Error("expected description-trigger diagnostic")
	}
}

func TestRuleDescriptionIncludesTrigger_HasTrigger(t *testing.T) {
	good := []string{
		"Processes PDFs. Use when working with PDF files.",
		"Generates commit messages. Trigger on: git, commit, staged changes.",
		"Analyzes data. Use for spreadsheet analysis tasks.",
	}
	for _, desc := range good {
		s := makeSkill("name", desc, "body")
		diags := ruleDescriptionIncludesTrigger(s)
		if hasDiag(diags, "description-trigger") {
			t.Errorf("unexpected description-trigger diagnostic for %q", desc)
		}
	}
}

// --- Body tests ---

func TestRuleBodyMaxLines(t *testing.T) {
	body := strings.Repeat("line\n", 501)
	s := makeSkill("name", "desc", body)
	diags := ruleBodyMaxLines(s)
	if !hasDiag(diags, "body-max-lines") {
		t.Error("expected body-max-lines diagnostic")
	}
}

func TestRuleBodyMaxLines_AtLimit(t *testing.T) {
	body := strings.Repeat("line\n", 499) + "last line"
	s := makeSkill("name", "desc", body)
	diags := ruleBodyMaxLines(s)
	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}
}

func TestRuleNoWindowsPaths(t *testing.T) {
	body := `See scripts\helper.py for details`
	s := makeSkill("name", "desc", body)
	diags := ruleNoWindowsPaths(s)
	if !hasDiag(diags, "no-windows-paths") {
		t.Error("expected no-windows-paths diagnostic")
	}
}

func TestRuleNoWindowsPaths_ForwardSlash(t *testing.T) {
	body := `See scripts/helper.py for details`
	s := makeSkill("name", "desc", body)
	diags := ruleNoWindowsPaths(s)
	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}
}

func TestRuleNoTimeSensitive(t *testing.T) {
	body := "If you're doing this before August 2025, use the old API."
	s := makeSkill("name", "desc", body)
	diags := ruleNoTimeSensitive(s)
	if !hasDiag(diags, "no-time-sensitive") {
		t.Error("expected no-time-sensitive diagnostic")
	}
}

func TestRuleNoTimeSensitive_Clean(t *testing.T) {
	body := "Use the v2 API endpoint."
	s := makeSkill("name", "desc", body)
	diags := ruleNoTimeSensitive(s)
	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}
}

func TestRuleVagueName(t *testing.T) {
	for _, name := range []string{"helper", "utils", "tools", "misc", "common"} {
		s := makeSkill(name, "desc", "body")
		diags := ruleVagueName(s)
		if !hasDiag(diags, "name-not-vague") {
			t.Errorf("expected name-not-vague diagnostic for %q", name)
		}
	}
}

func TestRuleVagueName_Good(t *testing.T) {
	s := makeSkill("pdf-processing", "desc", "body")
	diags := ruleVagueName(s)
	if len(diags) > 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}
}
