package models

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Prompt struct {
	Display   string `json:"display"`
	Timestamp int64  `json:"timestamp"`
	Project   string `json:"project"`
}

func (p Prompt) Description() string {
	projectPath := p.ProjectPath()
	if timeAgo := p.TimeAgo(); timeAgo != "" {
		return projectPath + " • " + timeAgo
	}
	return projectPath
}

func (p Prompt) ProjectPath() string {
	if p.Project == "" {
		return "no project"
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return p.Project
	}

	home = strings.TrimSuffix(home, string(filepath.Separator))

	if p.Project == home {
		return "~"
	}
	if strings.HasPrefix(p.Project, home+string(filepath.Separator)) {
		return "~" + p.Project[len(home):]
	}

	return p.Project
}

type timeUnit struct {
	divisor  int64
	singular string
	plural   string
}

var timeUnits = []timeUnit{
	{31536000, "year", "years"},
	{2592000, "month", "months"},
	{604800, "week", "weeks"},
	{86400, "day", "days"},
	{3600, "hour", "hours"},
	{60, "minute", "minutes"},
}

func (p Prompt) TimeAgo() string {
	if p.Timestamp == 0 {
		return ""
	}

	elapsed := (time.Now().UnixMilli() - p.Timestamp) / 1000
	if elapsed < 60 {
		return "just now"
	}

	for _, unit := range timeUnits {
		if value := elapsed / unit.divisor; value >= 1 {
			if value == 1 {
				return "1 " + unit.singular + " ago"
			}
			return strconv.FormatInt(value, 10) + " " + unit.plural + " ago"
		}
	}

	return "just now"
}
