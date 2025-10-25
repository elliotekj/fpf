package matcher

import (
	"fpf/pkg/models"
	"testing"
)

func BenchmarkMatchPrompts(b *testing.B) {
	prompts := make([]models.Prompt, 100)
	for i := 0; i < 100; i++ {
		prompts[i] = models.Prompt{
			Display:   "fix the bug in authentication module",
			Project:   "/home/user/project",
			Timestamp: 1234567890,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MatchPrompts(prompts, "bug auth")
	}
}

func BenchmarkParseQuery(b *testing.B) {
	query := "implement feature %p myapp something else"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ParseQuery(query)
	}
}
