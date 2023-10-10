package server

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/lipgloss"
)

type infoModel struct {
	common *commonModel
	style lipgloss.Style
	table *table.Table
	turnRow int
	row1 string
	row2 string
	row3 string
}

type infoMsg int

func NewInfoModel(com *commonModel) *infoModel {
	i := infoModel {
		common: com,
		style: lipgloss.NewStyle().MarginLeft(3),
	}

	iptr := &i

	t := table.New().
	Border(lipgloss.NormalBorder()).
	StyleFunc(func(row, col int) lipgloss.Style {
		switch {
		case row == 2:
			return lipgloss.NewStyle().Bold(true)
		case row == iptr.turnRow:
			return lipgloss.NewStyle().Foreground(lipgloss.Color("201"))
		default:
			return lipgloss.NewStyle().Bold(false)
		}
	}).
	BorderRow(true).
	BorderStyle(lipgloss.NewStyle())

	iptr.table = t
	return iptr
}

func (m *infoModel) View(flipped bool) string {
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

	whiteKing := "♚"
	blackKing := "♔"

	whiteName = fmt.Sprintf("%s %s", whiteKing, whiteName)
	blackName = fmt.Sprintf("%s %s", blackKing, blackName)

	if flipped {
		if m.common.player.lob.game.WhiteTurn() {
			m.turnRow = 1
		} else {
			m.turnRow = 3
		}
	
		rows = [][]string{[]string{whiteName}, []string{fmt.Sprintf("Code: %s", code)}, []string{blackName}}
	} else {
		if m.common.player.lob.game.WhiteTurn() {
			m.turnRow = 3
		} else {
			m.turnRow = 1
		}
	
		rows = [][]string{[]string{blackName}, []string{fmt.Sprintf("Code: %s", code)}, []string{whiteName}}
	}

	data := table.NewStringData(rows...)
	m.table.Data(data)

	return m.style.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.table.String() + "\n\n\n" + lipgloss.NewStyle().Faint(true).Render("ctrl+c / esc to exit\nctrl+f to flip board"),
		),
	)
}