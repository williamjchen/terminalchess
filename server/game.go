package server

import (
	"log/slog"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

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
	s.WriteString(m.game.game.PrintBoard())
	s.WriteString("\n\n")
	
	if m.game.whiteTurn {
		s.WriteString("White to move\n")
	} else {
		s.WriteString("Black to Move\n")
	}

	s.WriteString(m.textinput.View())
	return s.String()
}
