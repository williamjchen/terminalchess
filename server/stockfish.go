package server

import (
	tea "github.com/charmbracelet/bubbletea"
)

type stockfishModel struct {

}

func NewStockfishModel() stockfishModel {
	s := stockfishModel {}
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