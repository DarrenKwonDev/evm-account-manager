package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// AccountCreatedMsg is sent when user clicks Create button
type AccountCreatedMsg struct {
	Alias string
	Chain string
	Label string
	Memo  string
}

type AccountFormField int

const (
	NoFocus AccountFormField = iota
	AliasField
	ChainField
	LabelField
	MemoField
	CreateField
)

type AccountForm struct {
	AliasInput   textinput.Model
	ChainInput   textinput.Model
	LabelInput   textinput.Model
	MemoInput    textarea.Model
	CreateButton Button
	LastFocused  AccountFormField
}

func NewAccountForm() AccountForm {
	alias := textinput.New()
	alias.Placeholder = "Alias"
	alias.CharLimit = 50
	alias.Width = 20

	chain := textinput.New()
	chain.Placeholder = "Chain"
	chain.CharLimit = 50
	chain.Width = 20

	label := textinput.New()
	label.Placeholder = "Label"
	label.CharLimit = 100
	label.Width = 20

	memo := textarea.New()
	memo.Placeholder = "Memo"
	memo.CharLimit = 500
	memo.SetWidth(20)
	memo.SetHeight(4)

	// 버튼 초기화
	createBtn := NewButton("Create")
	createBtn.SetOnClick(func() tea.Cmd {
		return func() tea.Msg {
			return AccountCreatedMsg{
				Alias: alias.Value(),
				Chain: chain.Value(),
				Label: label.Value(),
				Memo:  memo.Value(),
			}
		}
	})

	return AccountForm{
		AliasInput:   alias,
		ChainInput:   chain,
		LabelInput:   label,
		MemoInput:    memo,
		CreateButton: createBtn,
		LastFocused:  NoFocus,
	}
}

func (af AccountForm) View() string {
	var builder strings.Builder

	builder.WriteString("Alias:\n")
	aliasView := af.AliasInput.View()
	builder.WriteString(aliasView)
	builder.WriteString("\n\n")

	builder.WriteString("Chain:\n")
	chainView := af.ChainInput.View()
	builder.WriteString(chainView)
	builder.WriteString("\n\n")

	builder.WriteString("Label:\n")
	labelView := af.LabelInput.View()
	builder.WriteString(labelView)
	builder.WriteString("\n\n")

	builder.WriteString("Memo:\n")
	memoView := af.MemoInput.View()
	builder.WriteString(memoView)
	builder.WriteString("\n\n")

	// 버튼 표시
	buttonView := af.CreateButton.View()
	builder.WriteString(buttonView)

	return builder.String()
}

func (af *AccountForm) Update(msg tea.Msg) tea.Cmd {
	// 현재 포커스된 필드만 업데이트
	var cmd tea.Cmd

	switch af.LastFocused {
	case AliasField:
		af.AliasInput, cmd = af.AliasInput.Update(msg)
	case ChainField:
		af.ChainInput, cmd = af.ChainInput.Update(msg)
	case LabelField:
		af.LabelInput, cmd = af.LabelInput.Update(msg)
	case MemoField:
		af.MemoInput, cmd = af.MemoInput.Update(msg)
	case CreateField:
		// 버튼은 포인터를 직접 수정하므로 cmd만 반환
		cmd = af.CreateButton.Update(msg)
	}

	return cmd
}

func (af *AccountForm) MoveFocus(direction int) {
	fields := []AccountFormField{AliasField, ChainField, LabelField,
		MemoField, CreateField}

	// 현재 포커스된 필드 인덱스 찾기
	currentIndex := -1
	for i, field := range fields {
		if field == af.LastFocused {
			currentIndex = i
			break
		}
	}

	// 포커스 해제
	af.blurAll()

	// 다음 인덱스 계산
	nextIndex := (currentIndex + direction + len(fields)) % len(fields)
	af.LastFocused = fields[nextIndex]

	// 새 포커스 설정
	af.focusField(af.LastFocused)
}

func (af *AccountForm) blurAll() {
	af.AliasInput.Blur()
	af.ChainInput.Blur()
	af.LabelInput.Blur()
	af.MemoInput.Blur()
	af.CreateButton.Blur()

	af.LastFocused = NoFocus
}

func (af *AccountForm) focusField(field AccountFormField) {
	switch field {
	case AliasField:
		af.AliasInput.Focus()
	case ChainField:
		af.ChainInput.Focus()
	case LabelField:
		af.LabelInput.Focus()
	case MemoField:
		af.MemoInput.Focus()
	case CreateField:
		af.CreateButton.Focus()
	}
}

func (af *AccountForm) ActivateFirst() {
	af.blurAll()
	af.LastFocused = AliasField
	af.AliasInput.Focus()
}

func (af *AccountForm) SetWidth(width int) {
	af.AliasInput.Width = width
	af.ChainInput.Width = width
	af.LabelInput.Width = width
	af.MemoInput.SetWidth(width)
	af.CreateButton.SetWidth(width)
}
