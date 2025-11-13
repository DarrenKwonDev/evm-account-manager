package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type AccountFormField int

const (
	NoFocus AccountFormField = iota
	AliasField
	ChainField
	LabelField
	MemoField
)

type AccountForm struct {
	AliasInput  textinput.Model
	ChainInput  textinput.Model
	LabelInput  textinput.Model
	MemoInput   textarea.Model
	LastFocused AccountFormField
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

	return AccountForm{
		AliasInput:  alias,
		ChainInput:  chain,
		LabelInput:  label,
		MemoInput:   memo,
		LastFocused: NoFocus,
	}
}

func (af AccountForm) View() string {
	var builder strings.Builder

	builder.WriteString("Alias:\n")
	builder.WriteString(af.AliasInput.View())
	builder.WriteString("\n\n")

	builder.WriteString("Chain:\n")
	builder.WriteString(af.ChainInput.View())
	builder.WriteString("\n\n")

	builder.WriteString("Label:\n")
	builder.WriteString(af.LabelInput.View())
	builder.WriteString("\n\n")

	builder.WriteString("Memo:\n")
	builder.WriteString(af.MemoInput.View())

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
	}

	return cmd
}

func (af *AccountForm) MoveFocus(direction int) {
	fields := []AccountFormField{AliasField, ChainField, LabelField,
		MemoField}

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
}
