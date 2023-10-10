package server

import (
	tea "github.com/charmbracelet/bubbletea"
)

type createModel struct {
	common *commonModel
	gs *gameState
}

func NewCreateModel(com *commonModel, gs *gameState) createModel {
	c := createModel{
		common: com,
		gs: gs,
	}

	return c
}

func (m createModel) Init() tea.Cmd{
	return nil
}

func (m createModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m createModel) View() string {
	return ""
}