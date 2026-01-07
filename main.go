package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(initialModel())
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

type (
	errMsg error
)

type model struct {
	textInput textinput.Model
	secret    int
	attempts  int
	feedback  string
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
				m.feedback = "Please enter a number"
				return m, nil
			}
			m.attempts++
			if userInput < m.secret {
				m.feedback = "higher"
			} else if userInput > m.secret {
				m.feedback = "lower"
			} else {
				m.feedback = "correct"
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
	return fmt.Sprintf(
		"Guess the number (1-100)\n\n%s\n\n%s\nattempts: %d\n\n(enter to submit, esc to quit)\n",
		m.textInput.View(),
		m.feedback,
		m.attempts,
	)
}
