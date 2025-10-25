package history

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"fpf/pkg/models"
	"github.com/charlievieth/fastwalk"
)

type JSONLEntry struct {
	Type      string  `json:"type"`
	IsMeta    *bool   `json:"isMeta"`
	Cwd       string  `json:"cwd"`
	Message   Message `json:"message"`
	UUID      string  `json:"uuid"`
	Timestamp string  `json:"timestamp"`
}

type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

func GetProjectsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".claude", "projects"), nil
}

func ReadHistory() ([]models.Prompt, error) {
	projectsPath, err := GetProjectsPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(projectsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("projects directory does not exist: %s", projectsPath)
	}

	var prompts []models.Prompt

	conf := fastwalk.Config{
		Follow: false,
	}

	err = fastwalk.Walk(&conf, projectsPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".jsonl") {
			return nil
		}

		filePrompts, err := readJSONLFile(path)
		if err != nil {
			return err
		}

		prompts = append(prompts, filePrompts...)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking projects directory: %w", err)
	}

	return deduplicatePrompts(prompts), nil
}

func readJSONLFile(path string) ([]models.Prompt, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	prompts := make([]models.Prompt, 0, 64)
	scanner := bufio.NewScanner(file)

	const maxCapacity = 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		var entry JSONLEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue
		}

		if entry.Type != "user" || entry.Message.Role != "user" {
			continue
		}
		if entry.IsMeta != nil && *entry.IsMeta {
			continue
		}

		message := extractMessageContent(entry.Message.Content)
		if message == "" {
			continue
		}

		if shouldSkipMessage(message) {
			continue
		}

		timestamp := parseTimestamp(entry.Timestamp)

		prompts = append(prompts, models.Prompt{
			Display:   message,
			Timestamp: timestamp,
			Project:   entry.Cwd,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return prompts, nil
}

func extractMessageContent(content interface{}) string {
	switch v := content.(type) {
	case string:
		return v
	case []interface{}:
		var textParts []string
		for _, item := range v {
			if obj, ok := item.(map[string]interface{}); ok {
				if typ, ok := obj["type"].(string); ok && typ == "text" {
					if text, ok := obj["text"].(string); ok {
						textParts = append(textParts, text)
					}
				}
			}
		}
		return strings.Join(textParts, " ")
	default:
		return ""
	}
}

var skipPrefixes = []string{
	"<command-name>",
	"<local-command",
	"Caveat:",
	"/clear",
	"[Request interrupted",
}

func shouldSkipMessage(message string) bool {
	trimmed := strings.TrimSpace(message)

	if trimmed == "" || trimmed == "Warmup" {
		return true
	}

	for _, prefix := range skipPrefixes {
		if strings.HasPrefix(message, prefix) {
			return true
		}
	}

	return false
}

func parseTimestamp(timestamp string) int64 {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return time.Now().UnixMilli()
	}
	return t.UnixMilli()
}

func deduplicatePrompts(prompts []models.Prompt) []models.Prompt {
	seen := make(map[string]models.Prompt)

	for _, p := range prompts {
		if existing, exists := seen[p.Display]; !exists || p.Timestamp > existing.Timestamp {
			seen[p.Display] = p
		}
	}

	result := make([]models.Prompt, 0, len(seen))
	for _, p := range seen {
		result = append(result, p)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp > result[j].Timestamp
	})

	return result
}
