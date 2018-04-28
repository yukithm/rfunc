package server

import (
	"context"
	"net"

	pb "github.com/yukithm/rfunc/rfuncs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type RFunc struct {
	listener   net.Listener
	grpcServer *grpc.Server
}

func NewRFunc(lis net.Listener) *RFunc {
	s := &RFunc{
		listener: lis,
	}

	gs := grpc.NewServer()
	pb.RegisterRFuncsServer(gs, s)
	reflection.Register(gs)
	s.grpcServer = gs

	return s
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
	if req.GetClipContent().GetType() != pb.ClipboardType_TEXT {
		return nil, status.Error(codes.Unavailable, "Unsupported content type")
	}

	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}

func (f *RFunc) Paste(ctx context.Context, req *pb.PasteRequest) (*pb.PasteReply, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}

func (f *RFunc) OpenURL(ctx context.Context, req *pb.OpenURLRequest) (*pb.OpenURLReply, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}
