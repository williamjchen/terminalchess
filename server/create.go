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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String(){
		case "esc", "ctrl+c":
			return m, tea.Quit
		} // switch KeyMsg
		
	case lobMsg:
		if msg == nil {
			return m, nil
		} else {
			m.gs.lobby = msg
			return m, nil
		}

	case errMsg:
		return m, nil

	} // switch msg
	return m, nil
}

func (m createModel) View() string {
	return ""
}