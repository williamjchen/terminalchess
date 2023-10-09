package server

import (
	"log/slog"
    "strings"

	tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/ssh"
)

type menu struct {
    choices  []string           // items on the to-do list
    cursor   int                // which to-do list item our cursor is pointing at
    choice string
}

func GetMenuOption(s ssh.Session, options []string) string {
    p := tea.NewProgram(
        Menu(options),
        tea.WithInput(s),
        tea.WithOutput(s),
    )
    m, err := p.Run()
    if err != nil {
        slog.Error("failed to run menu", err)
        return ""
    }

	// Assert the final tea.Model to our local model and print the choice.
	if m, ok := m.(menu); ok && m.choice != "" {
		return m.choice
	}
    return ""
}


func Menu(options []string) menu {
	return menu{
		// Our to-do list is a grocery list
		choices:  options,
	}
}

func (m menu) Init() tea.Cmd {
    // Just return `nil`, which means "no I/O right now, please."
    return tea.EnterAltScreen
}

func (m menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit

		case "enter":
			// Send the choice on the channel and exit.
			m.choice = m.choices[m.cursor]
			return m, tea.Quit

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

func (m menu) View() string {
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
