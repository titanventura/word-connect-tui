package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const darkGray = lipgloss.Color("#767676")
const lightGreen = lipgloss.Color("#04B575")

type game struct {
	guesses      []int
	options      []rune
	answers      []string
	attemptsLeft int
}

type model struct {
	width     int
	height    int
	textInput textinput.Model
	game
	status string
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			val := m.textInput.Value()
			m.textInput.Reset()
			for _, guessIdx := range m.game.guesses {
				if m.answers[guessIdx] == val {
					return m, nil
				}
			}
			exists := false
			for idx, answer := range m.game.answers {
				if answer == val {
					m.guesses = append(m.guesses, idx)
					exists = true
				}
			}

			if !exists {
				if m.game.attemptsLeft > 0 {
					m.game.attemptsLeft--
				}
			}

			m.updateStatus()
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *model) gameCompleted() bool {
	return len(m.game.answers) == len(m.game.guesses)
}

func (m *model) updateStatus() {
	if m.game.attemptsLeft == 0 {
		m.status = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("#DC3535")).
			Bold(true).
			Render("You've lost the game :(")
	} else if m.gameCompleted() {
		m.status = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color(lightGreen)).
			Render(":) You've won the game 🎉")
	} else {
		m.status = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color(darkGray)).
			Render(
				fmt.Sprintf(
					"%d words to go\n%d attempts left",
					len(m.game.answers)-len(m.game.guesses),
					m.game.attemptsLeft,
				),
			)
	}
}

func (m *model) wrapComponentWithBorder(component string) string {
	return lipgloss.
		NewStyle().
		Width(m.width/2 - m.width/10).
		Height(m.height / 2).
		Border(lipgloss.NormalBorder()).
		Render(component)
}

func (m *model) inputAndStatusArea() string {
	var inputArea string

	if m.game.attemptsLeft <= 0 || m.gameCompleted() {
		return m.wrapComponentWithBorder(m.status)
	}

	inputArea = fmt.Sprintf(
		"Guess a word from the given letters ...\n\n%s",
		m.textInput.View(),
	)

	avlHeight := m.height / 2
	avlHeight -= lipgloss.Height(inputArea)
	avlHeight -= lipgloss.Height(m.status)

	emptySpace := lipgloss.NewStyle().Height(avlHeight).Render("")

	return m.wrapComponentWithBorder(lipgloss.JoinVertical(lipgloss.Left, inputArea, emptySpace, m.status))
}

func (m *model) renderAnswers() string {
	ans := []string{}

	for _, guessIdx := range m.game.guesses {
		ans = append(
			ans,
			lipgloss.
				NewStyle().
				Bold(true).
				Foreground(lightGreen).
				Render(m.game.answers[guessIdx]),
		)
	}
	ansRender := lipgloss.JoinVertical(lipgloss.Left, ans...)

	return ansRender
}

func (m *model) optionsAndAnswersArea() string {
	options := []string{}

	for _, option := range m.game.options {
		options = append(
			options,
			lipgloss.
				NewStyle().
				Border(lipgloss.NormalBorder()).
				PaddingLeft(1).
				PaddingRight(1).
				Render(string(option)),
		)
	}

	optionsArea := lipgloss.JoinHorizontal(lipgloss.Center, options...)
	answersArea := m.renderAnswers()

	availableHeight := m.height / 2
	availableHeight -= (lipgloss.Height(optionsArea) + lipgloss.Height(answersArea))
	availableHeight -= m.height / 10
	emptySpace := lipgloss.NewStyle().Height(availableHeight).Render("")

	return m.wrapComponentWithBorder(
		lipgloss.Place(
			m.width/2-m.width/10,
			m.height/2,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.JoinVertical(
				lipgloss.Center,
				optionsArea,
				emptySpace,
				answersArea,
			),
		),
	)
}

func (m model) View() string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			m.inputAndStatusArea(),
			m.optionsAndAnswersArea(),
		),
	)
}

func initialModel() model {

	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 20

	m := model{
		textInput: ti,
		game: game{
			options: []rune{'a', 'e', 't'},
			guesses: []int{},
			answers: []string{
				"tea",
				"eat",
				"ate",
				"at",
			},
			attemptsLeft: 5,
		},
	}

	m.updateStatus()

	return m
}
