package server

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/url"

	pb "github.com/yukithm/rfunc/rfuncs"
	"github.com/yukithm/rfunc/server/clipboard"
	"github.com/yukithm/rfunc/server/shell"
	"github.com/yukithm/rfunc/text"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
	f.initServer()
	return f.grpcServer.Serve(f.Listener)
}

func (f *RFunc) initServer() {
	if f.grpcServer != nil {
		return
	}

	gs := grpc.NewServer()
	pb.RegisterRFuncsServer(gs, f)
	reflection.Register(gs)
	f.grpcServer = gs
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
	f.Log().Println("[gRPC] Copy")

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
	f.Log().Println("[gRPC] Paste")

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

func (f *RFunc) convertLineEnding(str string) string {
	if f.Config == nil || f.Config.EOL == "" {
		return str
	}

	return text.ConvertLineEnding(str, f.Config.EOL)
}
