package server

import (
	"log/slog"
    "strings"

	Game "github.com/williamjchen/terminalchess/game"	

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
    "github.com/charmbracelet/ssh"
)

type (
	errMsg error
)
type model struct {
    choices  []string           // items on the to-do list
    cursor   int                // which to-do list item our cursor is pointing at
    choice string
	chosen bool
	viewport viewport.Model
	textinput textinput.Model
    game *Game.Game
}

func GetModelOption(s ssh.Session, options []string) {
    p := tea.NewProgram(
        Model(options),
        tea.WithInput(s),
        tea.WithOutput(s),
    )
    _, err := p.Run()
    if err != nil {
        slog.Error("failed to run menu", err)
        return
    }
}


func Model(options []string) model {
	ti := textinput.New()
	ti.Placeholder = "Enter move in algebraic notation..."
	ti.Focus()
	ti.CharLimit = 5
	ti.Width = 20

	return model{
		choices:  options,
		chosen: false,
		textinput: ti,
		game: Game.NewGame(),
	}
}

func (m model) Init() tea.Cmd {
    return tea.Batch(tea.EnterAltScreen, textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	if m.chosen {
		return updateChosen(msg, m)
	} else {
		return updateChoices(msg, m)
	}
}

func (m model) View() string {
	if m.chosen {
		return chosenView(m)
	} else {
		return choicesViews(m)
	}
}

func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			// Send the choice on the channel and exit.
			m.choice = m.choices[m.cursor]
			m.chosen = true	
			return m, nil

		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
		}
	}
	return m, nil
}

func updateChosen(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.textinput.Reset()
			return m, cmd
		}

	case errMsg:
		return m, nil
	}

	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func choicesViews(m model) string {
	s := strings.Builder{}
	s.WriteString("What chess mode would you like to play?\n\n")

	for i := 0; i < len(m.choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(m.choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press esc to quit)\n")

	return s.String()
}

func chosenView(m model) string {
	s := strings.Builder{}
	switch m.choice {
	case "":
		slog.Info("optiion", m.choice)
	case m.choices[0]:
		slog.Info("option", m.choice)
	case m.choices[1]:
		slog.Info("option", m.choice)
	case m.choices[2]:
		slog.Info("option", m.choice)
	}
	s.WriteString(m.game.PrintBoard())
	s.WriteString("\n")
	s.WriteString(m.textinput.View())
	return s.String()
}
