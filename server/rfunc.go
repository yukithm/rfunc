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

type RFunc struct {
	Config    *Config
	Logger    *log.Logger
	Listener  net.Listener
	Clipboard clipboard.Clipboard
	Shell     shell.Shell

	grpcServer *grpc.Server
}

func (f *RFunc) Log() *log.Logger {
	if f.Logger == nil {
		f.Logger = log.New(ioutil.Discard, "", 0)
	}
	return f.Logger
}

func (f *RFunc) Start() error {
	if f.grpcServer == nil {
		server, err := f.newServer()
		if err != nil {
			return err
		}
		f.grpcServer = server
	}
	return f.grpcServer.Serve(f.Listener)
}

func (f *RFunc) newServer() (*grpc.Server, error) {
	opts := make([]grpc.ServerOption, 0)
	if f.Config.TLS != nil && (f.Config.TLS.CertFile != "" || f.Config.TLS.CAFile != "") {
		tlsConfig, err := newTLSConfig(f.Config.TLS)
		if err != nil {
			return nil, err
		}
		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.Creds(creds))
	}

	gs := grpc.NewServer(opts...)
	pb.RegisterRFuncsServer(gs, f)
	reflection.Register(gs)

	return gs, nil
}

func (f *RFunc) Stop() {
	if f.grpcServer != nil {
		f.grpcServer.Stop()
		f.grpcServer = nil
	}
}

func (f *RFunc) GracefulStop() {
	if f.grpcServer != nil {
		f.grpcServer.GracefulStop()
		f.grpcServer = nil
	}
}

func (f *RFunc) Copy(ctx context.Context, req *pb.CopyRequest) (*pb.CopyReply, error) {
	if f.allowed("copy") {
		f.Log().Println("[gRPC] Copy")
	} else {
		f.Log().Println("[gRPC] Copy is not allowed")
		return nil, status.Error(codes.PermissionDenied, "copy command is not allowd")
	}

	contentType := req.GetClipContent().GetType()
	switch contentType {
	case pb.ClipboardType_TEXT:
		str := f.convertLineEnding(req.GetClipContent().GetText())
		err := f.Clipboard.CopyText(str)
		if err != nil {
			f.Log().Println("[gRPC] Copy:", err)
			return nil, status.Error(codes.Internal, err.Error())
		}
		return &pb.CopyReply{}, nil
	}

	f.Log().Println("[gRPC] Copy: Unsupported content type: ", contentType)
	return nil, status.Error(codes.Unavailable, "Unsupported content type")
}

func (f *RFunc) Paste(ctx context.Context, req *pb.PasteRequest) (*pb.PasteReply, error) {
	if f.allowed("paste") {
		f.Log().Println("[gRPC] Paste")
	} else {
		f.Log().Println("[gRPC] Paste is not allowed")
		return nil, status.Error(codes.PermissionDenied, "paste command is not allowd")
	}

	if !req.Acceptable(pb.ClipboardType_TEXT) {
		f.Log().Println("[gRPC] Paste: Unsupported content type")
		return nil, status.Error(codes.Unavailable, "Unsupported content type")
	}

	content, err := f.Clipboard.PasteText()
	content = f.convertLineEnding(content)
	if err != nil {
		f.Log().Println("[gRPC] Paste:", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.PasteReply{
		ClipContent: pb.MakeTextClipboardContent(content),
	}, nil
}

func (f *RFunc) OpenURL(ctx context.Context, req *pb.OpenURLRequest) (*pb.OpenURLReply, error) {
	f.Log().Println("[gRPC] OpenURL")
	if f.allowed("open") {
		f.Log().Println("[gRPC] OpenURL")
	} else {
		f.Log().Println("[gRPC] OpenURL is not allowed")
		return nil, status.Error(codes.PermissionDenied, "open command is not allowd")
	}

	urls := req.GetUrl()
	for _, ref := range urls {
		if err := validateURL(ref); err != nil {
			f.Log().Println("[gRPC] OpenURL:", err)
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	if err := f.Shell.OpenURL(urls...); err != nil {
		f.Log().Println("[gRPC] OpenURL:", err)
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

func (f *RFunc) allowed(name string) bool {
	if f.Config == nil || f.Config.AllowCmds == nil || len(f.Config.AllowCmds) == 0 {
		return true
	}

	return utils.FindString(f.Config.AllowCmds, name) != -1
}

func (f *RFunc) convertLineEnding(str string) string {
	if f.Config == nil || f.Config.EOL == "" {
		return str
	}

	return text.ConvertLineEnding(str, f.Config.EOL)
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
