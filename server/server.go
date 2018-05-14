package server

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/yukithm/rfunc/options"
	"github.com/yukithm/rfunc/server/clipboard"
	"github.com/yukithm/rfunc/server/shell"
)

type Config struct {
	EOL       string
	AllowCmds []string
	TLS       *options.TLSOptions
}

type Server struct {
	Config    *Config
	Clipboard clipboard.Clipboard
	Shell     shell.Shell
	Logger    *log.Logger
	rfunc     *RFunc
}

func (s *Server) Serve(lis net.Listener) error {
	if s.Clipboard == nil {
		clip, err := clipboard.GetClipboard()
		if err != nil {
			return err
		}
		s.Clipboard = clip
	}

	if s.Shell == nil {
		shell, err := shell.GetShell()
		if err != nil {
			return err
		}
		s.Shell = shell
	}

	s.rfunc = &RFunc{
		Config:    s.Config,
		Logger:    s.Logger,
		Listener:  lis,
		Clipboard: s.Clipboard,
		Shell:     s.Shell,
	}

	quit := make(chan struct{})
	defer close(quit)

	errCh := make(chan error, 1)
	defer close(errCh)

	go func() {
		err := s.rfunc.Start()
		errCh <- err
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range sigCh {
			quit <- struct{}{}
		}
	}()

	for {
		select {
		case err := <-errCh:
			return err

		case <-quit:
			s.rfunc.GracefulStop()
			break
		}
	}
}

func (s *Server) Stop() {
	s.rfunc.Stop()
}

func (s *Server) GracefulStop() {
	s.rfunc.GracefulStop()
}
