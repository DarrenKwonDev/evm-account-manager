package main

import (
	"drkup/account-tracker/db"
	"drkup/account-tracker/service"
	"fmt"
	"log"
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
	case AccountCreatedMsg:

		// 1. create event
		account, err := service.GetAccountService().CreateAccount(
			msg.Alias,
			msg.Chain,
			msg.Label,
			msg.Memo,
		)
		if err != nil {
			log.Printf("계정 생성 실패: %v", err)
			// 실패 처리 (에러 메시지 표시 등)
			return m, nil
		}
		log.Printf("계정 생성 성공: %+v", account)
		// 2. update right panel to draw new account

		// 3. 폼 초기화
		m.leftPane = NewAccountForm()
		if m.width > 0 {
			m.leftPane.SetWidth(LeftWidthFromWnd(m.width) - SmallPadding)
		}

		return m, nil

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
				case CreateField:

					return m, func() tea.Msg {
						return AccountCreatedMsg{
							Alias: m.leftPane.AliasInput.Value(),
							Chain: m.leftPane.ChainInput.Value(),
							Label: m.leftPane.LabelInput.Value(),
							Memo:  m.leftPane.MemoInput.Value(),
						}
					}

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
	defer func() {
		log.Println("Exit")
		f.Close()
	}()

	database, err := db.New("app.db")
	if err != nil {
		log.Printf("DB 초기화 실패: %v", err)
		os.Exit(1)
	}
	defer database.Close()

	// AccountService 초기화
	service.InitAccountService(database)

	// create program and run
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	if p != nil {
		if _, err := p.Run(); err != nil {
			fmt.Println("Error while running program:", err)
			os.Exit(1)
		}
	}
}
