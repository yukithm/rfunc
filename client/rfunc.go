package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/yukithm/rfunc/rfuncs"
	"google.golang.org/grpc"
)

const RPCTimeout = time.Second * 5

type RFunc struct {
	conn   *grpc.ClientConn
	rfuncs pb.RFuncsClient
}

func RunRFunc(network, addr string, f func(*RFunc) error) error {
	rfunc := &RFunc{}
	if err := rfunc.Connect(network, addr); err != nil {
		return err
	}
	defer rfunc.Close()

	return f(rfunc)
}

func (f *RFunc) Connect(network, addr string) error {
	conn, err := NewClientConn(network, addr)
	if err != nil {
		return err
	}
	f.conn = conn
	f.rfuncs = pb.NewRFuncsClient(f.conn)

	return nil
}

func (f *RFunc) Close() error {
	conn := f.conn
	if conn != nil {
		if err := conn.Close(); err != nil {
			return err
		}
		f.conn = nil
	}
	return nil
}

func (f *RFunc) Copy(text string) error {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()

	_, err := f.rfuncs.Copy(ctx, &pb.CopyRequest{
		ClipContent: pb.MakeTextClipboardContent(text),
	})
	return err
}

func (f *RFunc) Paste() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()

	res, err := f.rfuncs.Paste(ctx, &pb.PasteRequest{
		Accepts: []pb.ClipboardType{
			pb.ClipboardType_TEXT,
		},
	})
	if err != nil {
		return "", err
	}

	content := res.GetClipContent()
	switch content.GetType() {
	case pb.ClipboardType_TEXT:
		return content.GetText(), nil
	}

	return "", fmt.Errorf("Unsupported content: %s", content.GetType())
}

func (f *RFunc) OpenURL(url ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()

	_, err := f.rfuncs.OpenURL(ctx, &pb.OpenURLRequest{
		Url: url,
	})
	return err
}
