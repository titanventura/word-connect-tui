package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// const darkGray = lipgloss.Color("#767676")
const lightGreen = lipgloss.Color("#04B575")
const dotSeparator = "•"

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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			val := m.textInput.Value()
			m.textInput.Reset()

			if len(val) == 0 {
				return m, nil
			}

			if m.game.isAlreadyGuessed(val) {
				return m, nil
			}

			idxAtAnswer := m.game.isCorrectGuess(val)
			if idxAtAnswer != -1 {
				m.game.guesses = append(m.game.guesses, idxAtAnswer)
			} else {
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

func (g *game) isAlreadyGuessed(val string) bool {
	for _, guessIdx := range g.guesses {
		if g.answers[guessIdx] == val {
			return true
		}
	}
	return false
}

func (g *game) isCorrectGuess(val string) int {
	for idx, answer := range g.answers {
		if val == answer {
			return idx
		}
	}
	return -1
}

func (m *model) inputAndStatusArea() string {
	if m.attemptsOver() || m.gameCompleted() {
		return m.wrapComponentWithBorder(m.status)
	}

	inputArea := fmt.Sprintf(
		"Guess a word from the given letters ...\n\n%s",
		m.textInput.View(),
	)

	h := lipgloss.Height
	avlHeight := m.height/2 - h(inputArea) - h(m.status)
	vEmptySpace := lipgloss.NewStyle().Height(avlHeight).Render("")

	credits := creditStyle().Render("word-connector")

	w := lipgloss.Width
	avlWidth := (m.width/2 - m.width/10) - w(m.status) - w(credits)
	hEmptySpace := statusBarStyle().Width(avlWidth).Render("")

	statusBar := lipgloss.JoinHorizontal(lipgloss.Center, m.status, hEmptySpace, credits)

	return m.wrapComponentWithBorder(
		lipgloss.JoinVertical(
			lipgloss.Left,
			inputArea,
			vEmptySpace,
			statusBar,
		),
	)
}

func (m *model) updateStatus() {
	if m.attemptsOver() {
		m.status = getEndMessage(
			"#DC3535",
			"You've lost the game :(",
		)
	} else if m.gameCompleted() {
		m.status = getEndMessage(
			string(lightGreen),
			":) You've won the game 🎉",
		)
	} else {
		progressMsg := getStatusMessage(len(m.game.answers)-len(m.game.guesses), "words to go")
		attemptsMsg := getStatusMessage(m.game.attemptsLeft, "attempts left")

		m.status = lipgloss.JoinHorizontal(
			lipgloss.Center,
			statusStyle("#FF5F87").MarginLeft(1).Render(progressMsg),
			statusStyle("#FFFDF5").Render(dotSeparator),
			statusStyle("#A550DF").Render(attemptsMsg),
		)
	}
}

func (m *model) attemptsOver() bool {
	return m.game.attemptsLeft <= 0
}

func (m *model) gameCompleted() bool {
	return len(m.game.answers) == len(m.game.guesses)
}

func (m *model) wrapComponentWithBorder(component string) string {
	return lipgloss.
		NewStyle().
		Width(m.width/2 - m.width/10).
		Height(m.height / 2).
		Border(lipgloss.NormalBorder()).
		Render(component)
}

func (m *model) optionsAndAnswersArea() string {
	options := []string{}

	for _, option := range m.game.options {
		options = append(
			options,
			optionsStyle().
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

func (m *model) renderAnswers() string {
	ans := []string{}

	for _, guessIdx := range m.game.guesses {
		ans = append(
			ans,
			correctAnswerStyle().
				Render(m.game.answers[guessIdx]),
		)
	}
	ansRender := lipgloss.JoinVertical(lipgloss.Left, ans...)

	return ansRender
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

func getStatusMessage(status int, message string) string {
	return fmt.Sprintf("%d %s", status, message)
}
