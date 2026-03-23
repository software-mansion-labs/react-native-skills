package main

import "fmt"

type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// Diagnostic represents a single lint finding.
type Diagnostic struct {
	File     string
	Line     int // 0 means file-level (frontmatter)
	Severity Severity
	Rule     string
	Message  string
}

func (d Diagnostic) String() string {
	loc := d.File
	if d.Line > 0 {
		loc = fmt.Sprintf("%s:%d", d.File, d.Line)
	}
	return fmt.Sprintf("%s: %s [%s] %s", loc, d.Severity, d.Rule, d.Message)
}
