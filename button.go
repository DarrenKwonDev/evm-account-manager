package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Button represents a simple clickable button component
type Button struct {
	text    string
	width   int
	focused bool
}

// NewButton creates a new button with the specified text
func NewButton(text string) Button {
	return Button{
		text:    text,
		width:   len(text) + 4, // padding
		focused: false,
	}
}

// Focus sets the button to focused state
func (b *Button) Focus() {
	b.focused = true
}

// Blur sets the button to normal state
func (b *Button) Blur() {
	b.focused = false
}

// IsFocused returns true if button is focused
func (b *Button) IsFocused() bool {
	return b.focused
}

// SetWidth sets the button width
func (b *Button) SetWidth(width int) {
	b.width = width
}

func (b Button) Init() tea.Cmd {
	return nil
}

// View renders the button
func (b Button) View() string {
	style := b.getStyle()
	return style.Render(b.text)
}

// getStyle returns the appropriate lipgloss style based on button focus
func (b Button) getStyle() lipgloss.Style {
	baseStyle := lipgloss.NewStyle().
		Width(b.width).
		Align(lipgloss.Center).
		Padding(0, 2).
		Border(lipgloss.RoundedBorder())

	if b.focused {
		return baseStyle.BorderForeground(focusColor)
	}

	return baseStyle.BorderForeground(lipgloss.Color("237"))
}
