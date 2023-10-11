package server

import (
	"github.com/williamjchen/terminalchess/magic"
	"fmt"
	"os"
	"errors"
	"log/slog"
	"context"
	"time"
	"os/signal"
	"syscall"

	"github.com/williamjchen/terminalchess/magic"

	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/ssh"
)

type Server struct {
	host string
	path string
	port int
	mng *manager
	srv *ssh.Server
}

func NewServer(path, host string, port int) (*Server, error){
	server := Server{
		host: host,
		path: path,
		port: port,
		mng: NewManager(),
	}
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(path),
		wish.WithMiddleware(
			tui(&server),
		),
	)
	if err != nil {
		return nil, err
	}
	server.srv = s

	magic.Init() // want to initalize magic bitboadrs once

	return &server, nil
}

func (s *Server) Start() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)


	slog.Info("Starting SSH server", "host", s.host, "port", s.port)
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			slog.Error("could nto sttart", err)
		}
	}()

	<-done

	slog.Info("stopping server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.srv.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		slog.Error("could not stop server", "error", err)
	}
}
