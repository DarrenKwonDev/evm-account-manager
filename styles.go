package main

import "github.com/charmbracelet/lipgloss"

var (
	focusColor = lipgloss.Color("#f0b90b")

	SmallPadding = 4
)

func LeftWidthFromWnd(wndWidth int) int {
	return wndWidth / 4
}
