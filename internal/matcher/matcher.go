package matcher

import (
	"regexp"
	"strings"

	"fpf/pkg/models"
	"github.com/sahilm/fuzzy"
)

type Query struct {
	PromptQuery  string
	ProjectQuery string
}

var projectPattern = regexp.MustCompile(`%p\s+(\S+)`)

func ParseQuery(query string) Query {
	q := Query{}

	if matches := projectPattern.FindStringSubmatch(query); len(matches) > 1 {
		q.ProjectQuery = matches[1]
		q.PromptQuery = strings.Join(strings.Fields(projectPattern.ReplaceAllString(query, "")), " ")
	} else {
		q.PromptQuery = query
	}

	return q
}

func MatchPrompts(prompts []models.Prompt, query string) []models.Prompt {
	if query == "" {
		return prompts
	}

	parsedQuery := ParseQuery(query)

	filtered := prompts
	if parsedQuery.ProjectQuery != "" {
		filtered = make([]models.Prompt, 0, len(prompts))
		for _, p := range prompts {
			if matchesProject(p, parsedQuery.ProjectQuery) {
				filtered = append(filtered, p)
			}
		}
	}

	if parsedQuery.PromptQuery == "" {
		return filtered
	}

	texts := make([]string, len(filtered))
	for i, p := range filtered {
		texts[i] = p.Display
	}

	matches := fuzzy.Find(parsedQuery.PromptQuery, texts)

	result := make([]models.Prompt, len(matches))
	for i, match := range matches {
		result[i] = filtered[match.Index]
	}

	return result
}

func matchesProject(prompt models.Prompt, projectQuery string) bool {
	query := strings.ToLower(projectQuery)
	projectPath := strings.ToLower(prompt.ProjectPath())

	return len(fuzzy.Find(query, []string{projectPath})) > 0
}
