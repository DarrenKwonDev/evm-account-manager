package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	// "drkup/account-tracker/onchain"
)

type model struct {
	// pane
	leftPane  AccountForm
	rightPane string

	keymap keyMap
	width  int
	height int

	// components
	help help.Model

	// utils
	focusedPane Pane
}

func NewModel() model {
	h := help.New()
	h.ShowAll = true
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Color("237")) // 회색
	h.Styles.ShortDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("251"))

	return model{
		leftPane:    NewAccountForm(),
		rightPane:   "right",
		keymap:      keys,
		help:        h,
		focusedPane: RightPane,
	}
}

type keyMap struct {
	quit  key.Binding
	esc   key.Binding
	tab   key.Binding
	down  key.Binding
	up    key.Binding
	enter key.Binding
}

var keys = keyMap{
	quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "blur"),
	),
	tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "focus/next"),
	),
	down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "down"),
	),
	up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "up"),
	),
	enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.quit, k.tab, k.down, k.up, k.enter}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

// Init implements tea.Model.
func (m model) Init() tea.Cmd {
	return tea.WindowSize()
}

// Update implements tea.Model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			return m, tea.Quit

		case key.Matches(msg, m.keymap.esc):
			// 왼쪽에 포커스가 있을 때 포커스를 전부 해제한다
			if m.focusedPane == LeftPane {
				m.leftPane.blurAll()
			}

			return m, nil

		case key.Matches(msg, m.keymap.enter):
			// focus가 없는 상태라면 첫번째 textinput에 포커스를 주고, 이미 포커스가 있으면 다음 input이동
			if m.focusedPane == LeftPane {
				switch m.leftPane.LastFocused {
				case NoFocus:
					m.leftPane.ActivateFirst()
					return m, nil
				case MemoField:
					// textarea의 경우엔 enter의 줄바꿈 동작을 그대로 유지해야 한다.
					return m, m.leftPane.Update(msg)
				default:
					m.leftPane.MoveFocus(1)
					return m, nil
				}
			}

		case key.Matches(msg, m.keymap.tab):
			// 왼쪽에서 포커스가 없다면 다음 pane으로 이동, 있다면 다음 input으로 이동
			if m.focusedPane == LeftPane {
				switch m.leftPane.LastFocused {
				case NoFocus:
					m.focusedPane = RightPane
					return m, nil
				default:
					m.leftPane.MoveFocus(1)
					return m, nil
				}
			} else {
				m.focusedPane = LeftPane
				return m, nil
			}

		case key.Matches(msg, m.keymap.down):
			if m.focusedPane == LeftPane {
				if m.leftPane.LastFocused != MemoField {
					m.leftPane.MoveFocus(1)
					return m, nil
				}
			}

		case key.Matches(msg, m.keymap.up):
			if m.focusedPane == LeftPane {
				if m.leftPane.LastFocused != MemoField {
					m.leftPane.MoveFocus(-1)
					return m, nil
				}
			}
		}

		cmds = append(cmds, m.leftPane.Update(msg))

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.leftPane.SetWidth(LeftWidthFromWnd(m.width) - SmallPadding)
		return m, nil
	}

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

// View implements tea.Model.
func (m model) View() string {

	leftWidth := LeftWidthFromWnd(m.width)
	rightWidth := m.width - leftWidth - SmallPadding

	basePaneStyle := lipgloss.NewStyle().MaxHeight(m.height - 20)
	focusedStyle := basePaneStyle.Border(lipgloss.NormalBorder()).BorderForeground(focusColor)
	notFocusedStyle := basePaneStyle.Border(lipgloss.NormalBorder())

	var left, right string

	if m.focusedPane == LeftPane {
		left = focusedStyle.Width(leftWidth).Render(m.leftPane.View())
		right = notFocusedStyle.Width(rightWidth).Render(m.rightPane)
	} else {
		left = notFocusedStyle.Width(leftWidth).Render(m.leftPane.View())
		right = focusedStyle.Width(rightWidth).Render(m.rightPane)
	}

	content := lipgloss.JoinHorizontal(lipgloss.Left, left, right)

	// draw help
	helpView := m.help.View(m.keymap)
	layout := lipgloss.JoinVertical(lipgloss.Top, content, helpView)

	return lipgloss.NewStyle().
		MaxWidth(m.width).
		MaxHeight(90).
		Padding(0, 0).
		Margin(0, 0).
		Render(layout)
}

func main() {
	// always leave log
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	// addr, pk, err := onchain.CreateAccount()
	// if err != nil {
	// os.Exit(1)
	// }
	// fmt.Printf("%s \n%s \n", addr, pk)

	// create program and run
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	if p != nil {
		if _, err := p.Run(); err != nil {
			fmt.Println("Error while running program:", err)
			os.Exit(1)
		}
	}
}
