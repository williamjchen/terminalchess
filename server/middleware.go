package server

import (
	"log/slog"

	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/ssh"
)

func tui(server *Server) wish.Middleware {
	return func(sh ssh.Handler) ssh.Handler {
		return func(s ssh.Session) {
			slog.Info("middlware")
		
			options := []string{"Stockfish", "Join Room", "Create Room"}
			opt := GetMenuOption(s, options)

			switch opt {
			case "":
				slog.Info("optiion", opt)
			case options[0]:
				slog.Info("option", opt)
			case options[1]:
				slog.Info("option", opt)
			case options[2]:
				slog.Info("option", opt)
			}

			sh(s)
		}
	}
}
