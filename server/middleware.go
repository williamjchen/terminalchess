package server

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func tui(server *Server) wish.Middleware {
	return func(sh ssh.Handler) ssh.Handler {
		lipgloss.SetColorProfile(termenv.ANSI256)
		return func(s ssh.Session) {
			slog.Info("middlware")
		
			options := []string{"Stockfish", "Join Room", "Create Room", "Basic AI"}
			p := GetModelOption(s, options, server, s)
			if p != nil {
				_, windowChanges, _ := s.Pty()
				go func() {
					for {
						select {
						case <-s.Context().Done():
							if p != nil {
								p.Quit()
								return
							}
						case w := <-windowChanges:
							if p != nil {
								p.Send(tea.WindowSizeMsg{Width: w.Width, Height: w.Height})
							}
						}
					}
				}()
				if _, err := p.Run(); err != nil {
					slog.Error("app exit with error", "error", err)
				}
				// p.Kill() will force kill the program if it's still running,
				// and restore the terminal to its original state in case of a
				// tui crash
				p.Kill()
			}			

			sh(s)
		}
	}
}
