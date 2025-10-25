package models

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProjectPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	tests := []struct {
		name     string
		project  string
		expected string
	}{
		{
			name:     "empty project",
			project:  "",
			expected: "no project",
		},
		{
			name:     "exactly home directory",
			project:  home,
			expected: "~",
		},
		{
			name:     "path under home directory",
			project:  filepath.Join(home, "projects", "myapp"),
			expected: "~/projects/myapp",
		},
		{
			name:     "path not under home directory",
			project:  "/var/log/app",
			expected: "/var/log/app",
		},
		{
			name:     "path that starts with home but different user",
			project:  home + "smith/projects",
			expected: home + "smith/projects",
		},
		{
			name:     "root path",
			project:  "/",
			expected: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Prompt{
				Project: tt.project,
			}
			got := p.ProjectPath()
			if got != tt.expected {
				t.Errorf("ProjectPath() = %q, expected %q", got, tt.expected)
			}
		})
	}
}

func TestDescription(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	tests := []struct {
		name      string
		prompt    Prompt
		wantPath  string
		wantTime  bool
	}{
		{
			name: "with timestamp",
			prompt: Prompt{
				Project:   filepath.Join(home, "projects", "app"),
				Timestamp: 1000000000000,
			},
			wantPath: "~/projects/app",
			wantTime: true,
		},
		{
			name: "without timestamp",
			prompt: Prompt{
				Project:   filepath.Join(home, "code"),
				Timestamp: 0,
			},
			wantPath: "~/code",
			wantTime: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.prompt.Description()

			// Check that shortened path is in the description
			if !contains(got, tt.wantPath) {
				t.Errorf("Description() = %q, expected to contain %q", got, tt.wantPath)
			}

			// Check timestamp presence
			hasTime := contains(got, "•")
			if hasTime != tt.wantTime {
				t.Errorf("Description() has time = %v, expected %v", hasTime, tt.wantTime)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
