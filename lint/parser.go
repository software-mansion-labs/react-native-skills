package main

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// SkillFile represents a parsed .md skill or reference file.
type SkillFile struct {
	Path           string
	Frontmatter    Frontmatter
	HasFrontmatter bool // true if the file starts with ---
	IsSkillMD      bool // true if the filename is SKILL.md
	Body           string
	BodyLines      []string
	Raw            string
}

// Frontmatter holds the YAML frontmatter fields.
type Frontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`

	// Track which fields were explicitly set.
	HasName        bool
	HasDescription bool
}

func parseSkillFile(path string) (*SkillFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	raw := string(data)
	isSkillMD := strings.ToUpper(filepath.Base(path)) == "SKILL.MD"

	fm, body, hasFM := splitFrontmatter(raw)

	bodyLines := strings.Split(body, "\n")

	return &SkillFile{
		Path:           path,
		Frontmatter:    fm,
		HasFrontmatter: hasFM,
		IsSkillMD:      isSkillMD,
		Body:           body,
		BodyLines:      bodyLines,
		Raw:            raw,
	}, nil
}

// splitFrontmatter extracts YAML frontmatter if present.
// Returns the frontmatter, the body, and whether frontmatter was found.
func splitFrontmatter(content string) (Frontmatter, string, bool) {
	trimmed := strings.TrimSpace(content)
	if !strings.HasPrefix(trimmed, "---") {
		return Frontmatter{}, content, false
	}

	rest := trimmed[3:]
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return Frontmatter{}, content, false
	}

	yamlContent := rest[:idx]
	body := strings.TrimLeft(rest[idx+4:], "\n")

	var rawMap map[string]interface{}
	if err := yaml.Unmarshal([]byte(yamlContent), &rawMap); err != nil {
		return Frontmatter{}, content, false
	}

	var fm Frontmatter
	if err := yaml.Unmarshal([]byte(yamlContent), &fm); err != nil {
		return Frontmatter{}, content, false
	}

	_, fm.HasName = rawMap["name"]
	_, fm.HasDescription = rawMap["description"]

	return fm, body, true
}
