package server

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type joinModel struct {
	common *commonModel
	idInput textinput.Model
}

func NewJoinModel(com *commonModel) joinModel {
	ii := textinput.New()
	ii.Placeholder = "XXXXXX"
	ii.CharLimit = 6
	ii.Width = 20

	j := joinModel {
		common: com,
		idInput: ii,
	}

	return j
}

func (m joinModel) Init() tea.Cmd {
    return nil
}

func (m joinModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	m.idInput.Focus()

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.common.begin = true
			m.idInput.Reset()
			return m, cmd
		}

	case errMsg:
		return m, nil
	}

	m.idInput, cmd = m.idInput.Update(msg)
	return m, cmd
}

func (m joinModel) View() string {
	s := strings.Builder{}
	s.WriteString("Enter Room Code...\n")
	s.WriteString(m.idInput.View())
	return s.String()
}