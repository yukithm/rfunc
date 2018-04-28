package server

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/url"

	pb "github.com/yukithm/rfunc/rfuncs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type RFunc struct {
	Logger     *log.Logger
	listener   net.Listener
	grpcServer *grpc.Server
	clipboard  Clipboard
	shell      Shell
}

func NewRFunc(lis net.Listener, clipboard Clipboard, shell Shell) *RFunc {
	s := &RFunc{
		listener:  lis,
		clipboard: clipboard,
		shell:     shell,
	}

	gs := grpc.NewServer()
	pb.RegisterRFuncsServer(gs, s)
	reflection.Register(gs)
	s.grpcServer = gs

	return s
}

func (f *RFunc) Log() *log.Logger {
	if f.Logger == nil {
		f.Logger = log.New(ioutil.Discard, "", 0)
	}
	return f.Logger
}

func (f *RFunc) Start() error {
	return f.grpcServer.Serve(f.listener)
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
		err := f.clipboard.CopyText(req.GetClipContent().GetText())
		if err != nil {
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
		return nil, status.Error(codes.Unavailable, "Unsupported content type")
	}

	content, err := f.clipboard.PasteText()
	if err != nil {
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
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	if err := f.shell.OpenURL(urls...); err != nil {
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
