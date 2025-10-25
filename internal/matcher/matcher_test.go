package matcher

import (
	"fpf/pkg/models"
	"testing"
)

func TestParseQuery(t *testing.T) {
	tests := []struct {
		input           string
		expectedPrompt  string
		expectedProject string
	}{
		{
			input:           "fix bug",
			expectedPrompt:  "fix bug",
			expectedProject: "",
		},
		{
			input:           "fix bug %p website",
			expectedPrompt:  "fix bug",
			expectedProject: "website",
		},
		{
			input:           "%p regdelete",
			expectedPrompt:  "",
			expectedProject: "regdelete",
		},
		{
			input:           "implement feature %p myapp something else",
			expectedPrompt:  "implement feature something else",
			expectedProject: "myapp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseQuery(tt.input)
			if result.PromptQuery != tt.expectedPrompt {
				t.Errorf("ParseQuery(%q).PromptQuery = %q, want %q", tt.input, result.PromptQuery, tt.expectedPrompt)
			}
			if result.ProjectQuery != tt.expectedProject {
				t.Errorf("ParseQuery(%q).ProjectQuery = %q, want %q", tt.input, result.ProjectQuery, tt.expectedProject)
			}
		})
	}
}

func TestMatchPrompts(t *testing.T) {
	prompts := []models.Prompt{
		{Display: "fix the bug in authentication", Project: "/home/user/website"},
		{Display: "add new feature to dashboard", Project: "/home/user/webapp"},
		{Display: "refactor database code", Project: "/home/user/website"},
		{Display: "update documentation", Project: "/home/user/docs"},
	}

	tests := []struct {
		query         string
		expectedCount int
		description   string
	}{
		{
			query:         "",
			expectedCount: 4,
			description:   "empty query returns all",
		},
		{
			query:         "bug",
			expectedCount: 1,
			description:   "simple fuzzy match",
		},
		{
			query:         "fix",
			expectedCount: 1,
			description:   "match beginning",
		},
		{
			query:         "%p website",
			expectedCount: 2,
			description:   "project filter only",
		},
		{
			query:         "fix %p website",
			expectedCount: 1,
			description:   "prompt and project filter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := MatchPrompts(prompts, tt.query)
			if len(result) != tt.expectedCount {
				t.Errorf("MatchPrompts with query %q returned %d results, want %d", tt.query, len(result), tt.expectedCount)
			}
		})
	}
}
