package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/inventage-ai/asylum/internal/term"
)

// Tab defines one tab in a tabbed TUI.
type Tab struct {
	Title       string
	Description string
	Kind        StepKind
	Options     []Option
	DefaultIdx  int   // default selection for StepSelect
	DefaultSel  []int // default selections for StepMultiSelect
}

// TabResult holds the outcome of a single tab.
type TabResult struct {
	SelectIdx int   // selected index for StepSelect
	MultiIdx  []int // selected indices for StepMultiSelect
	Completed bool
}

// RunTabs runs a tabbed TUI and returns results for each tab.
// Returns default selections without prompting if stdin is not a TTY.
func RunTabs(tabs []Tab) ([]TabResult, error) {
	results := make([]TabResult, len(tabs))

	if !term.IsTerminal() {
		for i, t := range tabs {
			results[i].Completed = true
			if t.Kind == StepSelect {
				results[i].SelectIdx = t.DefaultIdx
			} else {
				results[i].MultiIdx = append([]int(nil), t.DefaultSel...)
			}
		}
		return results, nil
	}

	m := tabsModel{
		tabs:    tabs,
		results: results,
	}
	m.initTab(0)

	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return nil, err
	}

	final := result.(tabsModel)
	if final.cancelled {
		return final.results, ErrCancelled
	}
	return final.results, nil
}

type tabsModel struct {
	tabs      []Tab
	results   []TabResult
	active    int
	cancelled bool

	selModel   selectModel
	multiModel multiModel
}

func (m *tabsModel) initTab(idx int) {
	t := m.tabs[idx]
	if t.Kind == StepSelect {
		m.selModel = selectModel{
			options: t.Options,
			cursor:  t.DefaultIdx,
		}
	} else {
		selected := map[int]bool{}
		for _, i := range t.DefaultSel {
			selected[i] = true
		}
		m.multiModel = multiModel{
			options:  t.Options,
			selected: selected,
		}
	}
}

// saveTab captures the current tab's state into results.
func (m *tabsModel) saveTab() {
	t := m.tabs[m.active]
	if t.Kind == StepSelect {
		m.results[m.active].SelectIdx = m.selModel.cursor
	} else {
		var indices []int
		for i := range m.multiModel.options {
			if m.multiModel.selected[i] {
				indices = append(indices, i)
			}
		}
		m.results[m.active].MultiIdx = indices
	}
}

func (m tabsModel) Init() tea.Cmd { return nil }

func (m tabsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			m.cancelled = true
			return m, tea.Quit

		case "left", "h":
			if m.active > 0 {
				m.saveTab()
				m.active--
				m.initTab(m.active)
			}
			return m, nil

		case "right", "l":
			if m.active < len(m.tabs)-1 {
				m.saveTab()
				m.active++
				m.initTab(m.active)
			}
			return m, nil

		case "enter":
			m.saveTab()
			for i := range m.results {
				m.results[i].Completed = true
			}
			return m, tea.Quit

		default:
			t := m.tabs[m.active]
			if t.Kind == StepSelect {
				switch msg.String() {
				case "up", "k":
					if m.selModel.cursor > 0 {
						m.selModel.cursor--
					}
				case "down", "j":
					if m.selModel.cursor < len(m.selModel.options)-1 {
						m.selModel.cursor++
					}
				}
			} else {
				switch msg.String() {
				case "up", "k":
					if m.multiModel.cursor > 0 {
						m.multiModel.cursor--
					}
				case "down", "j":
					if m.multiModel.cursor < len(m.multiModel.options)-1 {
						m.multiModel.cursor++
					}
				case " ":
					m.multiModel.selected[m.multiModel.cursor] = !m.multiModel.selected[m.multiModel.cursor]
				}
			}
		}
	}
	return m, nil
}

var (
	tabInactive = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Padding(0, 1)
	tabActive   = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true).Background(lipgloss.Color("63")).Padding(0, 1)
	tabSep      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	tabDesc     = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).MarginTop(2).MarginBottom(2)
)

func (m tabsModel) View() string {
	var b strings.Builder

	// Tab bar
	for i, t := range m.tabs {
		if i > 0 {
			b.WriteString(tabSep.Render(" │ "))
		}
		if i == m.active {
			b.WriteString(tabActive.Render(t.Title))
		} else {
			b.WriteString(tabInactive.Render(t.Title))
		}
	}
	b.WriteByte('\n')

	t := m.tabs[m.active]

	// Tab description
	if t.Description != "" {
		b.WriteString(tabDesc.Render(t.Description))
		b.WriteByte('\n')
	}

	// Tab content
	if t.Kind == StepSelect {
		b.WriteString(m.renderSelect())
	} else {
		b.WriteString(m.renderMulti())
	}

	return "\n" + borderStyle.Render(b.String()) + "\n"
}

func (m tabsModel) renderSelect() string {
	var b strings.Builder
	for i, opt := range m.selModel.options {
		if i > 0 {
			b.WriteByte('\n')
		}
		cursor := "  "
		label := labelStyle.Render(opt.Label)
		if i == m.selModel.cursor {
			cursor = selectedStyle.Render("▸ ")
			label = selectedStyle.Render(opt.Label)
		}
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, label))
		if opt.Description != "" {
			b.WriteString(descStyle.Render(opt.Description))
			b.WriteByte('\n')
		}
	}
	b.WriteString(hintStyle.Render("←/→ switch tab  •  ↑/↓ navigate  •  enter confirm  •  esc cancel"))
	return b.String()
}

func (m tabsModel) renderMulti() string {
	var b strings.Builder
	for i, opt := range m.multiModel.options {
		if i > 0 {
			b.WriteByte('\n')
		}
		cursor := "  "
		if i == m.multiModel.cursor {
			cursor = selectedStyle.Render("▸ ")
		}
		check := uncheckStyle.Render("[ ]")
		label := labelStyle.Render(opt.Label)
		if m.multiModel.selected[i] {
			check = checkStyle.Render("[✓]")
		}
		if i == m.multiModel.cursor {
			label = selectedStyle.Render(opt.Label)
		}
		b.WriteString(fmt.Sprintf("%s%s %s\n", cursor, check, label))
		if opt.Description != "" {
			b.WriteString(descStyle.Render(opt.Description))
			b.WriteByte('\n')
		}
	}
	b.WriteString(hintStyle.Render("←/→ switch tab  •  ↑/↓ navigate  •  space toggle  •  enter confirm  •  esc cancel"))
	return b.String()
}
