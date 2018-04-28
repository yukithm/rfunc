package server

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	Clipboard Clipboard
	Shell     Shell
	Logger    *log.Logger
	rfunc     *RFunc
}

func (s *Server) Serve(lis net.Listener) error {
	if s.Clipboard == nil {
		clip, err := GetClipboard()
		if err != nil {
			return err
		}
		s.Clipboard = clip
	}

	if s.Shell == nil {
		shell, err := GetShell()
		if err != nil {
			return err
		}
		s.Shell = shell
	}

	s.rfunc = NewRFunc(lis, s.Clipboard, s.Shell)
	s.rfunc.Logger = s.Logger

	quit := make(chan struct{})
	defer close(quit)

	errCh := make(chan error, 1)
	defer close(errCh)

	go func() {
		err := s.rfunc.Start()
		errCh <- err
	}()
	defer s.rfunc.GracefulStop()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range sigCh {
			quit <- struct{}{}
		}
	}()

	select {
	case err, ok := <-errCh:
		if ok {
			return err
		}

	case <-quit:
		break
	}

	return nil
}

func (s *Server) Stop() {
	s.rfunc.Stop()
}

func (s *Server) GracefulStop() {
	s.rfunc.GracefulStop()
}
