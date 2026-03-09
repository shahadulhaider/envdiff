package sync

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/shahadulhaider/envdiff/internal/env"
)

var (
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	addedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	removedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	changedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	headerStyle   = lipgloss.NewStyle().Bold(true)
	footerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

type SyncItem struct {
	Entry    env.DiffEntry
	Selected bool
}

type Model struct {
	Source  string
	Target  string
	Items   []SyncItem
	Cursor  int
	Applied bool
	Quit    bool
}

func NewModel(source, target string, entries []env.DiffEntry) Model {
	items := make([]SyncItem, len(entries))
	for i, e := range entries {
		items[i] = SyncItem{Entry: e, Selected: true}
	}
	return Model{Source: source, Target: target, Items: items}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.Quit = true
			return m, tea.Quit
		case "enter":
			m.Applied = true
			return m, tea.Quit
		case " ":
			if len(m.Items) > 0 {
				m.Items[m.Cursor].Selected = !m.Items[m.Cursor].Selected
			}
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Items)-1 {
				m.Cursor++
			}
		}
	}
	return m, nil
}

func (m Model) View() tea.View {
	var sb strings.Builder
	sb.WriteString(headerStyle.Render(fmt.Sprintf("Sync: %s -> %s", m.Source, m.Target)))
	sb.WriteString("\n\n")

	if len(m.Items) == 0 {
		sb.WriteString("  No differences found.\n")
	} else {
		for i, item := range m.Items {
			cursor := "  "
			if i == m.Cursor {
				cursor = cursorStyle.Render("> ")
			}
			checkbox := "[ ]"
			if item.Selected {
				checkbox = selectedStyle.Render("[x]")
			}
			var desc string
			switch item.Entry.Type {
			case env.DiffAdded:
				val := ""
				if item.Entry.Right != nil {
					val = item.Entry.Right.Value
				}
				desc = addedStyle.Render(fmt.Sprintf("+ %s=%s", item.Entry.Key, val))
			case env.DiffRemoved:
				val := ""
				if item.Entry.Left != nil {
					val = item.Entry.Left.Value
				}
				desc = removedStyle.Render(fmt.Sprintf("- %s=%s", item.Entry.Key, val))
			case env.DiffChanged:
				leftVal, rightVal := "", ""
				if item.Entry.Left != nil {
					leftVal = item.Entry.Left.Value
				}
				if item.Entry.Right != nil {
					rightVal = item.Entry.Right.Value
				}
				desc = changedStyle.Render(fmt.Sprintf("~ %s: %s -> %s", item.Entry.Key, leftVal, rightVal))
			}
			sb.WriteString(fmt.Sprintf("%s%s %s\n", cursor, checkbox, desc))
		}
	}

	sb.WriteString("\n")
	sb.WriteString(footerStyle.Render("space: toggle  enter: apply  q: cancel"))
	sb.WriteString("\n")
	return tea.NewView(sb.String())
}

func (m Model) SelectedEntries() []env.DiffEntry {
	var result []env.DiffEntry
	for _, item := range m.Items {
		if item.Selected {
			result = append(result, item.Entry)
		}
	}
	return result
}
