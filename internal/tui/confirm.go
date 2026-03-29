package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/inventage-ai/asylum/internal/term"
)

// Confirm shows a yes/no prompt and returns the user's choice.
// Returns defaultYes without prompting if stdin is not a TTY.
func Confirm(title string, defaultYes bool) (bool, error) {
	if !term.IsTerminal() {
		return defaultYes, nil
	}

	m := confirmModel{
		title: title,
		yes:   defaultYes,
	}

	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return false, err
	}

	final := result.(confirmModel)
	if final.cancelled {
		return false, ErrCancelled
	}
	return final.yes, nil
}

type confirmModel struct {
	title     string
	yes       bool
	cancelled bool
}

func (m confirmModel) Init() tea.Cmd { return nil }

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h", "tab":
			m.yes = !m.yes
		case "right", "l":
			m.yes = !m.yes
		case "y", "Y":
			m.yes = true
			return m, tea.Quit
		case "n", "N":
			m.yes = false
			return m, tea.Quit
		case "enter":
			return m, tea.Quit
		case "esc", "ctrl+c", "q":
			m.cancelled = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m confirmModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(m.title))
	b.WriteByte('\n')
	b.WriteByte('\n')

	yes := "  Yes  "
	no := "  No  "
	if m.yes {
		yes = selectedStyle.Render("▸ Yes ")
		no = labelStyle.Render("  No ")
	} else {
		yes = labelStyle.Render("  Yes ")
		no = selectedStyle.Render("▸ No ")
	}

	b.WriteString(fmt.Sprintf("%s    %s", yes, no))
	b.WriteByte('\n')
	b.WriteString(hintStyle.Render("←/→ toggle  •  y/n  •  enter confirm  •  esc cancel"))

	return "\n" + borderStyle.Render(b.String()) + "\n"
}
