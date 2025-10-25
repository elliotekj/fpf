package ui

import (
	"fmt"
	"io"
	"strings"

	"fpf/internal/matcher"
	"fpf/pkg/models"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultWidth      = 80
	defaultHeight     = 20
	minListHeight     = 4
	viewportPadding   = 8
	viewportOverhead  = 4
	filterInputMargin = 4
	itemPadding       = 4
	ellipsisWidth     = 1
)

var (
	accentColor = lipgloss.AdaptiveColor{
		Light: "#D75F00",
		Dark:  "#AF87FF",
	}
	mutedColor = lipgloss.AdaptiveColor{
		Light: "#949494",
		Dark:  "#6C6C6C",
	}
	lightMutedColor = lipgloss.AdaptiveColor{
		Light: "#A8A8A8",
		Dark:  "#7C7C7C",
	}
	separatorColor = lipgloss.AdaptiveColor{
		Light: "#CCCCCC",
		Dark:  "#444444",
	}
	borderColor = lipgloss.AdaptiveColor{
		Light: "#0087D7",
		Dark:  "#5FAFFF",
	}

	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(accentColor)
	projectStyle      = lipgloss.NewStyle().PaddingLeft(4).Foreground(mutedColor)
	filterInputStyle  = lipgloss.NewStyle().PaddingLeft(2)
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4).PaddingTop(1)
	helpStyle         = lipgloss.NewStyle().Foreground(mutedColor).PaddingLeft(4).PaddingTop(1)
	helpKeyStyle      = lipgloss.NewStyle().Foreground(mutedColor)
	helpDescStyle     = lipgloss.NewStyle().Foreground(lightMutedColor)
	helpSepStyle      = lipgloss.NewStyle().Foreground(separatorColor)
	previewStyle      = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(borderColor).
				Padding(1, 2).
				MarginLeft(2).
				MarginRight(2)
	previewTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(accentColor).MarginBottom(1)

	helpText = buildHelpText()
)

func buildHelpText() string {
	sep := " " + helpSepStyle.Render("•") + " "
	return helpStyle.Render(
		helpKeyStyle.Render("↑/↓") + " " + helpDescStyle.Render("navigate") + sep +
			helpKeyStyle.Render("ctrl+p") + " " + helpDescStyle.Render("preview") + sep +
			helpKeyStyle.Render("enter") + " " + helpDescStyle.Render("select") + sep +
			helpKeyStyle.Render("esc") + " " + helpDescStyle.Render("quit"),
	)
}

func firstLine(s string) string {
	if idx := strings.Index(s, "\n"); idx >= 0 {
		return s[:idx] + "…"
	}
	return s
}

type item struct {
	prompt models.Prompt
}

func (i item) FilterValue() string { return i.prompt.Display }
func (i item) Title() string       { return firstLine(i.prompt.Display) }
func (i item) Description() string { return i.prompt.Description() }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 2 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	title, desc := renderItem(i, index == m.Index(), m.Width())
	fmt.Fprintf(w, "%s\n%s", title, desc)
}

func renderItem(i item, selected bool, width int) (string, string) {
	titleText := i.Title()
	availableWidth := width - itemPadding

	if lipgloss.Width(titleText) > availableWidth {
		runes := []rune(titleText)
		targetWidth := availableWidth - ellipsisWidth
		for i := len(runes); i > 0; i-- {
			candidate := string(runes[:i])
			if lipgloss.Width(candidate) <= targetWidth {
				titleText = candidate + "…"
				break
			}
		}
	}

	var title string
	if selected {
		title = selectedItemStyle.Render("> " + titleText)
	} else {
		title = itemStyle.Render(titleText)
	}
	return title, projectStyle.Render(i.Description())
}

type Model struct {
	list        list.Model
	filterInput textinput.Model
	viewport    viewport.Model
	choice      string
	quitting    bool
	previewing  bool
	allPrompts  []models.Prompt
}

func promptsToItems(prompts []models.Prompt) []list.Item {
	items := make([]list.Item, len(prompts))
	for i, p := range prompts {
		items[i] = item{prompt: p}
	}
	return items
}

func configureListKeyMap(l *list.Model) {
	l.KeyMap.CursorUp.SetKeys("up")
	l.KeyMap.CursorDown.SetKeys("down")
	l.KeyMap.NextPage.SetKeys("pgdown")
	l.KeyMap.PrevPage.SetKeys("pgup")
	l.KeyMap.GoToStart.SetKeys("home")
	l.KeyMap.GoToEnd.SetKeys("end")
	l.KeyMap.Filter.SetEnabled(false)
	l.KeyMap.ClearFilter.SetEnabled(false)
	l.KeyMap.CancelWhileFiltering.SetEnabled(false)
	l.KeyMap.AcceptWhileFiltering.SetEnabled(false)
	l.KeyMap.ShowFullHelp.SetEnabled(false)
	l.KeyMap.CloseFullHelp.SetEnabled(false)
	l.KeyMap.Quit.SetEnabled(false)
	l.KeyMap.ForceQuit.SetEnabled(false)
}

func NewModel(prompts []models.Prompt) Model {
	items := promptsToItems(prompts)
	l := list.New(items, itemDelegate{}, defaultWidth, defaultHeight)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	configureListKeyMap(&l)

	ti := textinput.New()
	ti.Placeholder = "Type to filter (use '%p project' to filter by project)..."
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(mutedColor)
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = defaultWidth

	vp := viewport.New(defaultWidth-2, defaultHeight)

	return Model{
		list:        l,
		filterInput: ti,
		viewport:    vp,
		allPrompts:  prompts,
		previewing:  false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		listHeight := msg.Height
		if listHeight%2 != 0 {
			listHeight--
		}
		if listHeight < minListHeight {
			listHeight = minListHeight
		}

		m.list.SetWidth(msg.Width)
		m.list.SetHeight(listHeight)
		m.filterInput.Width = msg.Width - filterInputMargin
		m.viewport.Width = msg.Width - viewportPadding
		viewportHeight := msg.Height - viewportOverhead
		if viewportHeight%2 != 0 {
			viewportHeight--
		}
		m.viewport.Height = viewportHeight
		return m, nil

	case tea.KeyMsg:
		if m.previewing {
			switch msg.String() {
			case "esc":
				m.previewing = false
				return m, nil
			case "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			default:
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}
		}

		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "ctrl+p":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.previewing = true
				title := previewTitleStyle.Render("Preview - Press 'esc' to exit")
				wrappedContent := lipgloss.NewStyle().Width(m.viewport.Width).Render(i.prompt.Display)
				content := title + "\n\n" + wrappedContent
				m.viewport.SetContent(content)
			}
			return m, nil

		case "esc":
			if m.filterInput.Value() != "" {
				m.filterInput.SetValue("")
				m.updateFilteredList()
			} else {
				m.quitting = true
				return m, tea.Quit
			}
			return m, nil

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.prompt.Display
			}
			return m, tea.Quit

		case "up", "down", "pgup", "pgdown":
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd

		default:
			var cmd tea.Cmd
			m.filterInput, cmd = m.filterInput.Update(msg)
			cmds = append(cmds, cmd)
			m.updateFilteredList()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) updateFilteredList() {
	query := m.filterInput.Value()
	filtered := matcher.MatchPrompts(m.allPrompts, query)
	m.list.SetItems(promptsToItems(filtered))
	m.list.Select(0)
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	if m.previewing {
		return "\n" + previewStyle.Render(m.viewport.View())
	}

	var s strings.Builder

	s.WriteString("\n")
	start, end := m.list.Paginator.GetSliceBounds(len(m.list.Items()))
	for i, listItem := range m.list.Items()[start:end] {
		if item, ok := listItem.(item); ok {
			actualIndex := start + i
			title, desc := renderItem(item, actualIndex == m.list.Index(), m.list.Width())
			s.WriteString(title)
			s.WriteString("\n")
			s.WriteString(desc)
			s.WriteString("\n")
		}
	}

	s.WriteString("\n")
	s.WriteString(filterInputStyle.Render(m.filterInput.View()))
	s.WriteString("\n")

	if len(m.list.Items()) > 0 {
		s.WriteString(paginationStyle.Render(m.list.Paginator.View()))
		s.WriteString("\n")
	}

	s.WriteString(helpText)

	return s.String()
}

func (m Model) Choice() string {
	return m.choice
}
