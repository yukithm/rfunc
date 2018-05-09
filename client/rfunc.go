package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/yukithm/rfunc/rfuncs"
	"github.com/yukithm/rfunc/text"
	"google.golang.org/grpc"
)

const RPCTimeout = time.Second * 5

type Config struct {
	EOL string
}

type RFunc struct {
	Config *Config
	conn   *grpc.ClientConn
	rfuncs pb.RFuncsClient
}

func RunRFunc(network, addr string, config *Config, f func(*RFunc) error) error {
	rfunc := &RFunc{Config: config}
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

func (f *RFunc) Copy(str string) error {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()

	str = f.convertLineEnding(str)
	_, err := f.rfuncs.Copy(ctx, &pb.CopyRequest{
		ClipContent: pb.MakeTextClipboardContent(str),
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
		str := f.convertLineEnding(content.GetText())
		return str, nil
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

func (f *RFunc) convertLineEnding(str string) string {
	if f.Config == nil || f.Config.EOL == "" {
		return str
	}

	return text.ConvertLineEnding(str, f.Config.EOL)
}
