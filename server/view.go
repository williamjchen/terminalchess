package server

import (
	"log/slog"
    "strings"

	Game "github.com/williamjchen/terminalchess/game"	

	tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/ssh"
)

type model struct {
    choices  []string           // items on the to-do list
    cursor   int                // which to-do list item our cursor is pointing at
    choice string
	chosen bool
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
	return model{
		choices:  options,
		chosen: false,
		game: Game.NewGame(),
	}
}

func (m model) Init() tea.Cmd {
    return tea.EnterAltScreen
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
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
		case "ctrl+c", "q", "esc":
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
	return m, nil
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
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}

func chosenView(m model) string {
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
	return m.game.PrintBoard()
}
