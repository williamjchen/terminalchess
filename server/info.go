package server

import (
	"fmt"

	Game "github.com/williamjchen/terminalchess/game"	

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
	var statusMessage string = ""

	green := lipgloss.NewStyle().Foreground(lipgloss.Color("154"))
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	whiteKing := "♚"
	blackKing := "♔"

	if m.common.player.lob.p1 != nil {
		whiteName = m.common.player.lob.p1.name
	}
	if m.common.player.lob.p2 != nil {
		blackName = m.common.player.lob.p2.name
	}
	if m.common.player.lob.p1 != nil {
		code = m.common.player.lob.id
	}

	var bullet string
	if m.common.player.lob.p1Pres {
		bullet = "•"
	} else {
		bullet = "◦"
	}
	whiteName = fmt.Sprintf("%s %s %s", bullet, whiteKing, whiteName)

	if m.common.player.lob.p2Pres {
		bullet = "•"
	} else {
		bullet = "◦"
	}
	blackName = fmt.Sprintf("%s %s %s", bullet, blackKing, blackName)

	if flipped {
		if m.common.player.lob.game.Turn() == Game.WhiteTurn {
			m.turnRow = 1
		} else {
			m.turnRow = 3
		}
	
		rows = [][]string{[]string{whiteName}, []string{fmt.Sprintf("Code: %s", code)}, []string{blackName}}
	} else {
		if m.common.player.lob.game.Turn() == Game.WhiteTurn {
			m.turnRow = 3
		} else {
			m.turnRow = 1
		}
	
		rows = [][]string{[]string{blackName}, []string{fmt.Sprintf("Code: %s", code)}, []string{whiteName}}
	}

	data := table.NewStringData(rows...)
	m.table.Data(data)

	var col lipgloss.Style
	if m.common.player.playerType == white && m.common.player.lob.status == whiteWin ||  m.common.player.playerType == black && m.common.player.lob.status == blackWin {
		col = green
	} else if m.common.player.playerType == white && m.common.player.lob.status == blackWin || m.common.player.playerType == black && m.common.player.lob.status == whiteWin {
		col = red
	} else {
		col = lipgloss.NewStyle().Foreground(lipgloss.Color("254"))
	}

	switch m.common.player.lob.status {
	case inProgres:
		statusMessage = ""
	case whiteWin:
		statusMessage = col.Render("White Win")
	case blackWin:
		statusMessage = col.Render("Black Win")
	case stalemate:
		statusMessage = col.Render("Stalemate")
	}

	return m.style.Render(
		fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			m.table.String(),
			statusMessage,
			lipgloss.NewStyle().Faint(true).Render("ctrl+c / esc to exit\nctrl+f to flip board\nctrl+n to return to menu"),
		),
	)
}