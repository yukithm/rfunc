package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/yukithm/rfunc/options"
	pb "github.com/yukithm/rfunc/rfuncs"
	"github.com/yukithm/rfunc/server/clipboard"
	"github.com/yukithm/rfunc/server/shell"
	"github.com/yukithm/rfunc/text"
	"github.com/yukithm/rfunc/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Config struct {
	EOL       string
	AllowCmds []string
	TLS       *options.TLSOptions
}

type RFuncServer struct {
	Config     *Config
	Clipboard  clipboard.Clipboard
	Shell      shell.Shell
	Logger     *log.Logger
	grpcServer *grpc.Server
}

func (s *RFuncServer) Log() *log.Logger {
	if s.Logger == nil {
		s.Logger = log.New(ioutil.Discard, "", 0)
	}
	return s.Logger
}

func (s *RFuncServer) Run(lis net.Listener) error {
	quit := make(chan struct{})
	defer close(quit)

	errCh := make(chan error, 1)
	defer close(errCh)

	go func() {
		err := s.Serve(lis)
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
			s.GracefulStop()
			break
		}
	}
}

func (s *RFuncServer) Serve(lis net.Listener) error {
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

	if s.grpcServer == nil {
		server, err := s.newServer()
		if err != nil {
			return err
		}
		s.grpcServer = server
	}
	return s.grpcServer.Serve(lis)
}

func (s *RFuncServer) newServer() (*grpc.Server, error) {
	opts := make([]grpc.ServerOption, 0)
	if s.Config.TLS != nil && (s.Config.TLS.CertFile != "" || s.Config.TLS.CAFile != "") {
		tlsConfig, err := newTLSConfig(s.Config.TLS)
		if err != nil {
			return nil, err
		}
		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.Creds(creds))
	}

	gs := grpc.NewServer(opts...)
	pb.RegisterRFuncsServer(gs, s)
	reflection.Register(gs)

	return gs, nil
}

func (s *RFuncServer) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.Stop()
		s.grpcServer = nil
	}
}

func (s *RFuncServer) GracefulStop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
		s.grpcServer = nil
	}
}

func (s *RFuncServer) Copy(ctx context.Context, req *pb.CopyRequest) (*pb.CopyReply, error) {
	if s.allowed("copy") {
		s.Log().Println("[gRPC] Copy")
	} else {
		s.Log().Println("[gRPC] Copy is not allowed")
		return nil, status.Error(codes.PermissionDenied, "copy command is not allowd")
	}

	contentType := req.GetClipContent().GetType()
	switch contentType {
	case pb.ClipboardType_TEXT:
		str := s.convertLineEnding(req.GetClipContent().GetText())
		err := s.Clipboard.CopyText(str)
		if err != nil {
			s.Log().Println("[gRPC] Copy:", err)
			return nil, status.Error(codes.Internal, err.Error())
		}
		return &pb.CopyReply{}, nil
	}

	s.Log().Println("[gRPC] Copy: Unsupported content type: ", contentType)
	return nil, status.Error(codes.Unavailable, "Unsupported content type")
}

func (s *RFuncServer) Paste(ctx context.Context, req *pb.PasteRequest) (*pb.PasteReply, error) {
	if s.allowed("paste") {
		s.Log().Println("[gRPC] Paste")
	} else {
		s.Log().Println("[gRPC] Paste is not allowed")
		return nil, status.Error(codes.PermissionDenied, "paste command is not allowd")
	}

	if !req.Acceptable(pb.ClipboardType_TEXT) {
		s.Log().Println("[gRPC] Paste: Unsupported content type")
		return nil, status.Error(codes.Unavailable, "Unsupported content type")
	}

	content, err := s.Clipboard.PasteText()
	content = s.convertLineEnding(content)
	if err != nil {
		s.Log().Println("[gRPC] Paste:", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.PasteReply{
		ClipContent: pb.MakeTextClipboardContent(content),
	}, nil
}

func (s *RFuncServer) OpenURL(ctx context.Context, req *pb.OpenURLRequest) (*pb.OpenURLReply, error) {
	s.Log().Println("[gRPC] OpenURL")
	if s.allowed("open") {
		s.Log().Println("[gRPC] OpenURL")
	} else {
		s.Log().Println("[gRPC] OpenURL is not allowed")
		return nil, status.Error(codes.PermissionDenied, "open command is not allowd")
	}

	urls := req.GetUrl()
	for _, ref := range urls {
		if err := validateURL(ref); err != nil {
			s.Log().Println("[gRPC] OpenURL:", err)
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	if err := s.Shell.OpenURL(urls...); err != nil {
		s.Log().Println("[gRPC] OpenURL:", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.OpenURLReply{}, nil
}

func validateURL(ref string) error {
	u, err := url.Parse(ref)
	if err != nil {
		return err
	}

	if !u.IsAbs() {
		return errors.New("only full URL is allowed")
	}

	// restrict to HTTP only for security reason
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("only http and https schemes are allowed")
	}

	return nil
}

func (s *RFuncServer) allowed(name string) bool {
	if s.Config == nil || s.Config.AllowCmds == nil || len(s.Config.AllowCmds) == 0 {
		return true
	}

	return utils.FindString(s.Config.AllowCmds, name) != -1
}

func (s *RFuncServer) convertLineEnding(str string) string {
	if s.Config == nil || s.Config.EOL == "" {
		return str
	}

	return text.ConvertLineEnding(str, s.Config.EOL)
}

func newTLSConfig(opts *options.TLSOptions) (*tls.Config, error) {
	var certs []tls.Certificate
	if opts.CertFile != "" {
		cert, err := tls.LoadX509KeyPair(opts.CertFile, opts.KeyFile)
		if err != nil {
			return nil, err
		}
		certs = []tls.Certificate{cert}
	}

	var caCertPool *x509.CertPool
	if opts.CAFile != "" {
		caCertPem, err := ioutil.ReadFile(opts.CAFile)
		if err != nil {
			return nil, fmt.Errorf("Unable to read CA cert: %s", err)
		}
		caCertPool = x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCertPem); !ok {
			return nil, errors.New("Failed to append CA cert to the pool")
		}
	}

	tlsConfig := &tls.Config{
		Certificates:     certs,
		ClientAuth:       tls.RequireAndVerifyClientCert,
		ClientCAs:        caCertPool,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
	}
	tlsConfig.BuildNameToCertificate()

	return tlsConfig, nil
}
