package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	higherStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	lowerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
	wonStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	lostStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
)

var (
	purple    = lipgloss.Color("99")
	gray      = lipgloss.Color("245")
	lightGray = lipgloss.Color("241")

	headerStyle  = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
	cellStyle    = lipgloss.NewStyle().Padding(0, 1).Width(14)
	oddRowStyle  = cellStyle.Foreground(gray)
	evenRowStyle = cellStyle.Foreground(lightGray)
)

const maxAttempts int = 10

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

func generateRandomNumber() int {
	unixTime := time.Now().Unix()
	newSource := rand.NewSource(unixTime)
	newRand := rand.New(newSource)
	randomNumber := newRand.Intn(100) + 1
	return randomNumber
}

func createTable(rows [][]string) string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(purple)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return headerStyle
			case row%2 == 0:
				return evenRowStyle
			default:
				return oddRowStyle
			}
		}).
		Headers("Attempt", "Your Guess", "Result").
		Rows(rows...)

	return t.String()
}

type (
	errMsg error
)

type model struct {
	textInput    textinput.Model
	secret       int
	attemptsLeft int
	feedback     string
	history      [][]string
	gameOver     bool
	choice       int
	err          error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 3
	ti.Width = 3

	return model{
		textInput:    ti,
		secret:       generateRandomNumber(),
		attemptsLeft: maxAttempts,
		feedback:     "",
		history:      [][]string{},
		gameOver:     false,
		choice:       0,
		err:          nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.gameOver {
			switch msg.Type {
			case tea.KeyLeft, tea.KeyRight:
				m.choice = 1 - m.choice
				return m, nil
			case tea.KeyEnter:
				if m.choice == 0 {
					return initialModel(), nil
				}
				return m, tea.Quit
			}
			return m, nil
		}
		switch msg.Type {
		case tea.KeyEnter:
			userInput, err := strconv.Atoi(m.textInput.Value())
			m.textInput.SetValue("")
			if err != nil {
				m.feedback = "Enter a number, dude!"
				return m, nil
			}

			m.attemptsLeft--
			attemptsSoFar := maxAttempts - m.attemptsLeft

			if m.attemptsLeft == 0 {
				if userInput != m.secret {
					m.feedback = "You lost!"
					m.history = append(m.history, []string{strconv.Itoa(attemptsSoFar), strconv.Itoa(userInput), "❌"})
					m.gameOver = true
					m.choice = 0
					return m, nil
				} else {
					m.feedback = "You won!"
					m.history = append(m.history, []string{strconv.Itoa(attemptsSoFar), strconv.Itoa(userInput), "✅"})
					m.gameOver = true
					m.choice = 0
					return m, nil
				}
			} else {
				if userInput < m.secret {
					m.feedback = "Higher"
					m.history = append(m.history, []string{strconv.Itoa(attemptsSoFar), strconv.Itoa(userInput), "⬆️"})
				} else if userInput > m.secret {
					m.feedback = "Lower"
					m.history = append(m.history, []string{strconv.Itoa(attemptsSoFar), strconv.Itoa(userInput), "⬇️"})
				} else {
					m.feedback = "You won!"
					m.history = append(m.history, []string{strconv.Itoa(attemptsSoFar), strconv.Itoa(userInput), "✅"})
					m.gameOver = true
					m.choice = 0
					return m, nil
				}
			}

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) gameView() string {
	styledFeedback := m.feedback

	switch m.feedback {
	case "Higher":
		styledFeedback = higherStyle.Render(m.feedback)
	case "Lower":
		styledFeedback = lowerStyle.Render(m.feedback)
	case "You won!":
		styledFeedback = wonStyle.Render(m.feedback) + "😊"
	case "You lost!":
		revealedSecret := lipgloss.NewStyle().Foreground(purple).Bold(true).Render(strconv.Itoa(m.secret))
		styledFeedback = lostStyle.Render(m.feedback) + "😢" + "\nIt was " + revealedSecret + "!"
	}

	return fmt.Sprintf(
		"Welcome to the Guessing Game! 🎲 \n\nGuess the number between 1 and 100. You have %d attempts in total.\n\n%s\n\n%s\n\nAttempts left: %d\n\n%s\n\n(Press Enter to submit or Esc to quit)\n",
		maxAttempts,
		m.textInput.View(),
		styledFeedback,
		m.attemptsLeft,
		createTable(m.history),
	)
}

func replayBox(choice int) string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 3).
		BorderForeground(purple)

	btn := func(label string, selected bool) string {
		s := lipgloss.NewStyle().
			Padding(0, 3)
		if selected {
			s = s.Bold(true).Foreground(purple)
		}
		return s.Render(label)
	}

	yes := btn("Yes", choice == 0)
	no := btn("No", choice == 1)

	content := "Play again?" + lipgloss.JoinHorizontal(lipgloss.Center, yes, " ", no) + "\n\n(Use ← → to choose and Enter to confirm)"

	return box.Render(content)
}

func (m model) View() string {
	base := m.gameView()

	if m.gameOver {
		box := replayBox(m.choice)
		return base + "\n\n" + box + "\n"
	}

	return base
}
