package server

import (
	"context"
	"io/ioutil"
	"log"
	"net"

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
}

func NewRFunc(lis net.Listener, clipboard Clipboard) *RFunc {
	s := &RFunc{
		listener:  lis,
		clipboard: clipboard,
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
			return nil, err
		}
		return &pb.CopyReply{}, nil
	}

	f.Log().Println("[gRPC] Copy: Unsupported content type: ", contentType)
	return nil, status.Error(codes.Unavailable, "Unsupported content type")
}

func (f *RFunc) Paste(ctx context.Context, req *pb.PasteRequest) (*pb.PasteReply, error) {
	f.Log().Println("[gRPC] Paste")
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}

func (f *RFunc) OpenURL(ctx context.Context, req *pb.OpenURLRequest) (*pb.OpenURLReply, error) {
	f.Log().Println("[gRPC] OpenURL")
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}
