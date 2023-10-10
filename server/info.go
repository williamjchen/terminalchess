package server

import (
	"log/slog"
	"fmt"

	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/lipgloss"
)

type infoModel struct {
	common *commonModel
	style lipgloss.Style
	table *table.Table
	flipped bool
	row1 string
	row2 string
	row3 string
}

type infoMsg int

func NewInfoModel(com *commonModel) *infoModel {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 2:
				return lipgloss.NewStyle().Bold(true)
			default:
				return lipgloss.NewStyle().Bold(false)
			}
		}).
		BorderRow(true).
		BorderStyle(lipgloss.NewStyle())
	t.Row("fsfs")

	i := infoModel {
		common: com,
		style: lipgloss.NewStyle().MarginLeft(3),
		table: t,
		flipped: false,
	}

	return &i
}

func (m *infoModel) Flip() {
	slog.Info("Flip Info", "old:", m.flipped, "new:", !m.flipped)
	m.flipped = !m.flipped
	slog.Info("Flip Info", "cur:", m.flipped)

}

func (m *infoModel) View() string {
	var rows [][]string
	var whiteName, blackName, code string = "None", "Nil", "None"

	if m.common.player.lob.p1 != nil {
		whiteName = m.common.player.lob.p1.name
	}
	if m.common.player.lob.p2 != nil {
		blackName = m.common.player.lob.p2.name
	}
	if m.common.player.lob.p1 != nil {
		code = m.common.player.lob.id
	}

	if m.flipped {
		rows = [][]string{[]string{whiteName}, []string{fmt.Sprintf("Code: %s", code)}, []string{blackName}}
	} else {
		rows = [][]string{[]string{blackName}, []string{fmt.Sprintf("Code: %s", code)}, []string{whiteName}}
	}

	data := table.NewStringData(rows...)
	m.table.Data(data)

	return m.style.Render(m.table.String())
}