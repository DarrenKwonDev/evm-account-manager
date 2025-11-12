package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"drkup/account-tracker/onchain"
)

type model struct {
	leftPane  string
	rightPane string
	keymap    keyMap
	help      help.Model
	width     int
	height    int
	focused   Pane
}

func NewModel() model {
	return model{
		leftPane:  "left left  left  left  left  left  left  left  left  left  left ",
		rightPane: "right",
		keymap:    keys,
		focused:   RightPane,
	}
}

type keyMap struct {
	quit  key.Binding
	focus key.Binding
}

var keys = keyMap{
	quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c", "esc"),
		key.WithHelp("q", "quit"),
	),
	focus: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "focus next"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.quit, k.focus}
}

// Init implements tea.Model.
func (m model) Init() tea.Cmd {
	return tea.WindowSize()
}

// Update implements tea.Model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keymap.focus):
			if m.focused == LeftPane {
				m.focused = RightPane
			} else {
				m.focused = LeftPane
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

// View implements tea.Model.
func (m model) View() string {

	leftWidth := m.width / 4
	rightWidth := m.width - leftWidth - 4

	basePaneStyle := lipgloss.NewStyle().MaxHeight(m.height - 20)
	focusedStyle := basePaneStyle.Border(lipgloss.NormalBorder()).BorderForeground(focusColor)
	notFocusedStyle := basePaneStyle.Border(lipgloss.NormalBorder())

	var left, right string

	// 덮어씌우기
	if m.focused == LeftPane {
		left = focusedStyle.Width(leftWidth).Render(m.leftPane)
		right = notFocusedStyle.Width(rightWidth).Render(m.rightPane)
	} else {
		left = notFocusedStyle.Width(leftWidth).Render(m.leftPane)
		right = focusedStyle.Width(rightWidth).Render(m.rightPane)
	}

	content := lipgloss.JoinHorizontal(lipgloss.Left, left, right)

	return lipgloss.NewStyle().
		MaxWidth(m.width).
		MaxHeight(90).
		Padding(0, 0).
		Margin(0, 0).
		Render(content)
}

func main() {
	// always leave log
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	addr, pk, err := onchain.CreateAccount()
	if err != nil {
		os.Exit(1)
	}
	fmt.Printf("%s \n%s \n", addr, pk)

	// create program and run
	// p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	// if p != nil {
	// 	if _, err := p.Run(); err != nil {
	// 		fmt.Println("Error while running program:", err)
	// 		os.Exit(1)
	// 	}
	// }
}
