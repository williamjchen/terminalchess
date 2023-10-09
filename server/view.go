package server

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
    "github.com/charmbracelet/ssh"
)

type (
	errMsg error
)

type model struct {
    choices  []string           // items on the to-do list
    cursor   int                // which to-do list item our cursor is pointing at
    choice string
	id string
	chosen bool
	viewport viewport.Model
	textinput textinput.Model
	idInput textinput.Model
    lobby *lobby
}

func GetModelOption(s ssh.Session, options []string) {
    p := tea.NewProgram(
        Model(options),
        tea.WithInput(s),
        tea.WithOutput(s),
    )
    _, err := p.Run()
    if err != nil {
        slog.Error("failed to run menu", err)
        return
    }
}

func Model(options []string) model {
	ti := textinput.New()
	ti.Placeholder = "Enter move in algebraic notation..."
	ti.CharLimit = 5
	ti.Width = 20

	ii := textinput.New()
	ii.Placeholder = "XXXXXX"
	ii.CharLimit = 6
	ii.Width = 20

	return model{
		choices:  options,
		chosen: false,
		textinput: ti,
		idInput: ii,
	}
}

func (m model) Init() tea.Cmd {
    return tea.Batch(tea.EnterAltScreen, textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	if m.chosen {
		return updateChosen(msg, m)
	} else {
		return updateChoices(msg, m)
	}
}

func (m model) View() string {
	if m.chosen {
		return chosenView(m)
	} else {
		return choicesViews(m)
	}
}
