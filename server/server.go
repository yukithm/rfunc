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

type Server struct {
	listener   net.Listener
	grpcServer *grpc.Server
}

func NewServer(lis net.Listener) (*Server, error) {
	s := &Server{
		listener: lis,
	}

	gs := grpc.NewServer()
	pb.RegisterRFuncsServer(gs, s)
	reflection.Register(gs)
	s.grpcServer = gs

	return s, nil
}

func (s *Server) Start() error {
	return s.grpcServer.Serve(s.listener)
}

func (s *Server) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.Stop()
		s.grpcServer = nil
	}
}

func (s *Server) GracefulStop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
		s.grpcServer = nil
	}
}

func (s *Server) Copy(ctx context.Context, req *pb.CopyRequest) (*pb.CopyReply, error) {
	if req.GetClipContent().GetType() != pb.ClipboardType_TEXT {
		return nil, status.Error(codes.Unavailable, "Unsupported content type")
	}

	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}

func (s *Server) Paste(ctx context.Context, req *pb.PasteRequest) (*pb.PasteReply, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}

func (s *Server) OpenURL(ctx context.Context, req *pb.OpenURLRequest) (*pb.OpenURLReply, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}
