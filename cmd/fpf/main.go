package main

import (
	"fmt"
	"os"

	"fpf/internal/history"
	"fpf/internal/ui"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	mutedStyle   = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#949494",
		Dark:  "#6C6C6C",
	})
)

func main() {
	prompts, err := history.ReadHistory()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading history: %v\n", err)
		os.Exit(1)
	}

	if len(prompts) == 0 {
		fmt.Fprintln(os.Stderr, "No prompts found in history")
		os.Exit(1)
	}

	m := ui.NewModel(prompts)
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}

	if m, ok := finalModel.(ui.Model); ok {
		choice := m.Choice()
		if choice != "" {
			if err := clipboard.WriteAll(choice); err != nil {
				fmt.Fprintf(os.Stderr, "Error copying to clipboard: %v\n", err)
				os.Exit(1)
			}

			fmt.Println(successStyle.Render("✔ Copied prompt to clipboard"))
			fmt.Println(mutedStyle.Render(choice))
		}
	}
}
