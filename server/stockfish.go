package server

import (
	tea "github.com/charmbracelet/bubbletea"
)

type stockfishModel struct {
	common *commonModel
	gs *gameState
}

func NewStockfishModel(com *commonModel, gs *gameState) stockfishModel {
	s := stockfishModel {
		common: com,
		gs: gs,
	}
	return s
}

func (m stockfishModel) Init() tea.Cmd {
	return nil
}

func (m stockfishModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m stockfishModel) View() string {
	return ""
}