package server

import (
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	rfunc *RFunc
}

func (s *Server) Serve(lis net.Listener) error {
	s.rfunc = NewRFunc(lis)

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
