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
	higherStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	lowerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	correctStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)
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
		Headers("Attempt", "Your Guess", "Hint/Result").
		Rows(rows...)

	return t.String()
}

type (
	errMsg error
)

type model struct {
	textInput textinput.Model
	secret    int
	attempts  int
	feedback  string
	history   [][]string
	err       error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		secret:    generateRandomNumber(),
		attempts:  0,
		feedback:  "",
		history:   [][]string{},
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			userInput, err := strconv.Atoi(m.textInput.Value())
			m.textInput.SetValue("")
			if err != nil {
				m.feedback = "Enter a number, dude!"
				return m, nil
			}
			m.attempts++
			if userInput < m.secret {
				m.feedback = "Higher"
				m.history = append(m.history, []string{strconv.Itoa(m.attempts), strconv.Itoa(userInput), "⬆️"})
			} else if userInput > m.secret {
				m.feedback = "Lower"
				m.history = append(m.history, []string{strconv.Itoa(m.attempts), strconv.Itoa(userInput), "⬇️"})
			} else {
				m.feedback = "Correct!"
				m.history = append(m.history, []string{strconv.Itoa(m.attempts), strconv.Itoa(userInput), "✅"})
				return m, tea.Quit
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

func (m model) View() string {
	styledFeedback := m.feedback

	switch m.feedback {
	case "Higher":
		styledFeedback = higherStyle.Render(m.feedback)
	case "Lower":
		styledFeedback = lowerStyle.Render(m.feedback)
	case "Correct!":
		styledFeedback = correctStyle.Render(m.feedback)
	}

	return fmt.Sprintf(
		"Welcome to the Guessing Game!\n\nGuess the number between 1 and 100\n\n%s\n\n%s\nAttempts so far: %d\n\n%s\n\n(enter to submit, esc to quit)\n",
		m.textInput.View(),
		styledFeedback,
		m.attempts,
		createTable(m.history),
	)
}
