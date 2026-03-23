package main

import (
	"strings"
	"testing"
)

func TestSplitFrontmatter_Valid(t *testing.T) {
	content := "---\nname: my-skill\ndescription: Does things\n---\n# Body\nContent here"
	fm, body, hasFM := splitFrontmatter(content)
	if !hasFM {
		t.Fatal("expected frontmatter to be found")
	}
	if fm.Name != "my-skill" {
		t.Errorf("expected name 'my-skill', got %q", fm.Name)
	}
	if fm.Description != "Does things" {
		t.Errorf("expected description 'Does things', got %q", fm.Description)
	}
	if !fm.HasName || !fm.HasDescription {
		t.Error("expected HasName and HasDescription to be true")
	}
	if !strings.Contains(body, "# Body") {
		t.Errorf("body should contain '# Body', got %q", body)
	}
}

func TestSplitFrontmatter_NoFrontmatter(t *testing.T) {
	content := "# Just a regular markdown file"
	_, _, hasFM := splitFrontmatter(content)
	if hasFM {
		t.Error("expected no frontmatter")
	}
}

func TestSplitFrontmatter_NoClosing(t *testing.T) {
	content := "---\nname: test\ndescription: test\n# Body"
	_, _, hasFM := splitFrontmatter(content)
	if hasFM {
		t.Error("expected no frontmatter when closing --- is missing")
	}
}

func TestSplitFrontmatter_MissingFields(t *testing.T) {
	content := "---\nname: only-name\n---\n# Body"
	fm, _, hasFM := splitFrontmatter(content)
	if !hasFM {
		t.Fatal("expected frontmatter to be found")
	}
	if !fm.HasName {
		t.Error("expected HasName to be true")
	}
	if fm.HasDescription {
		t.Error("expected HasDescription to be false")
	}
}

func TestSplitFrontmatter_QuotedDescription(t *testing.T) {
	content := "---\nname: my-skill\ndescription: \"Quoted description with: colons\"\n---\n# Body"
	fm, _, hasFM := splitFrontmatter(content)
	if !hasFM {
		t.Fatal("expected frontmatter to be found")
	}
	if fm.Description != "Quoted description with: colons" {
		t.Errorf("expected unquoted description, got %q", fm.Description)
	}
}
