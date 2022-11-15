package main

import "github.com/charmbracelet/lipgloss"

// answers
func correctAnswerStyle() lipgloss.Style {
	return lipgloss.
		NewStyle().
		Bold(true).
		Foreground(lightGreen)
}

// options
func optionsStyle() lipgloss.Style {
	return lipgloss.
		NewStyle().
		Border(lipgloss.NormalBorder()).
		PaddingLeft(1).
		PaddingRight(1)
}

// Status Bar
func statusBarStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
		Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})
}

func statusStyle(fc string) lipgloss.Style {
	return lipgloss.NewStyle().
		Inherit(statusBarStyle()).
		Foreground(lipgloss.Color(fc)).
		Padding(0, 1)
}

func creditStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#6124DF")).
		MarginRight(1).
		MarginLeft(1).
		Padding(0, 1)
}

// end messages
func getEndMessage(color string, message string) string {
	return lipgloss.
		NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true).
		Render(message)
}
