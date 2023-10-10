package server

import (
	tea "github.com/charmbracelet/bubbletea"
)

type stockfishModel struct {
	common *commonModel
}

func NewStockfishModel(com *commonModel) *stockfishModel {
	s := stockfishModel {
		common: com,
	}
	return &s
}

func (m *stockfishModel) Init() tea.Cmd {
	return nil
}

func (m *stockfishModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *stockfishModel) View() string {
	return ""
}