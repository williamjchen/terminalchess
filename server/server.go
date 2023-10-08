package server

import (
	"fmt"
	"os"
	"errors"
	"log/slog"
	"context"
	"time"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/wish"
	"github.com/gliderlabs/ssh"
)


func NewServer(path, host string, port int) (error){
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(path),
	)
	if err != nil {
		return err
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)


	slog.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			slog.Error("could nto sttart", err)
		}
	}()

	<-done
	slog.Info("stopping server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		slog.Error("could not stop server", "error", err)
	}
	return nil
}
