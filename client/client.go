package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/yukithm/rfunc/rfuncs"
	"google.golang.org/grpc"
)

const RPCTimeout = time.Second * 5

type Client struct {
	conn   *grpc.ClientConn
	rfuncs pb.RFuncsClient
}

func (c *Client) Connect(network, addr string) error {
	if c.conn == nil {
		conn, err := NewClientConn(network, addr)
		if err != nil {
			return err
		}
		c.conn = conn
	}

	c.rfuncs = pb.NewRFuncsClient(c.conn)

	return nil
}

func (c *Client) Close() error {
	conn := c.conn
	if conn != nil {
		if err := conn.Close(); err != nil {
			return err
		}
		c.conn = nil
	}
	return nil
}

func (c *Client) Copy(text string) error {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()

	_, err := c.rfuncs.Copy(ctx, &pb.CopyRequest{
		ClipContent: pb.MakeTextClipboardContent(text),
	})
	return err
}

func (c *Client) Paste() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()

	res, err := c.rfuncs.Paste(ctx, &pb.PasteRequest{
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

func (c *Client) OpenURL(url ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), RPCTimeout)
	defer cancel()

	_, err := c.rfuncs.OpenURL(ctx, &pb.OpenURLRequest{
		Url: url,
	})
	return err
}
