package tui

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/inventage-ai/asylum/internal/term"
)

// ErrCancelled is returned when the user cancels a prompt.
var ErrCancelled = errors.New("cancelled")

// Select shows a single-choice prompt and returns the selected index.
// Returns defaultIdx without prompting if stdin is not a TTY.
func Select(title string, options []Option, defaultIdx int) (int, error) {
	if !term.IsTerminal() {
		return defaultIdx, nil
	}

	m := selectModel{
		title:   title,
		options: options,
		cursor:  defaultIdx,
	}

	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return -1, err
	}

	final := result.(selectModel)
	if final.cancelled {
		return -1, ErrCancelled
	}
	return final.cursor, nil
}

type selectModel struct {
	title     string
	options   []Option
	cursor    int
	cancelled bool
}

func (m selectModel) Init() tea.Cmd { return nil }

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter":
			return m, tea.Quit
		case "esc", "ctrl+c", "q":
			m.cancelled = true
			return m, tea.Quit
		}
	}
	return m, nil
}

var (
	borderStyle   = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2)
	titleStyle    = lipgloss.NewStyle().Bold(true).MarginBottom(1)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	labelStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	descStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).MarginLeft(4)
	hintStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("243")).Italic(true).MarginTop(1)
)

func (m selectModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(m.title))
	b.WriteByte('\n')

	for i, opt := range m.options {
		if i > 0 {
			b.WriteByte('\n')
		}

		cursor := "  "
		label := labelStyle.Render(opt.Label)
		if i == m.cursor {
			cursor = selectedStyle.Render("▸ ")
			label = selectedStyle.Render(opt.Label)
		}

		b.WriteString(fmt.Sprintf("%s%s\n", cursor, label))
		if opt.Description != "" {
			b.WriteString(descStyle.Render(opt.Description))
			b.WriteByte('\n')
		}
	}

	b.WriteString(hintStyle.Render("↑/↓ navigate  •  enter select  •  esc cancel"))

	return "\n" + borderStyle.Render(b.String()) + "\n"
}
