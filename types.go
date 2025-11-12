package main

type Pane int

const (
	LeftPane Pane = iota
	RightPane
)

func (p Pane) String() string {
	switch p {
	case LeftPane:
		return "left"
	case RightPane:
		return "right"
	default:
		return "unknown"
	}
}
