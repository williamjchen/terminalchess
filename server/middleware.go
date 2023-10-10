package server

import (
	"log/slog"
	//tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/ssh"
)

func tui(server *Server) wish.Middleware {
	return func(sh ssh.Handler) ssh.Handler {
		return func(s ssh.Session) {
			slog.Info("middlware")
		
			options := []string{"Stockfish", "Join Room", "Create Room"}
			GetModelOption(s, options, server, s)

			sh(s)
		}
	}
}
