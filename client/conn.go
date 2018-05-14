package client

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/yukithm/rfunc/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewClientConn(network, addr string, tlsOpts *options.TLSOptions) (*grpc.ClientConn, error) {
	opts := make([]grpc.DialOption, 0)
	if network == "unix" {
		opts = append(opts, grpc.WithDialer(unixDomainSocketDialer))
	}

	if tlsOpts != nil && (tlsOpts.CertFile != "" || tlsOpts.CAFile != "") {
		tlsConfig, err := newTLSConfig(tlsOpts)
		if err != nil {
			return nil, err
		}
		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	return grpc.Dial(addr, opts...)
}

func unixDomainSocketDialer(addr string, t time.Duration) (net.Conn, error) {
	return net.DialTimeout("unix", addr, t)
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
		Certificates:       certs,
		RootCAs:            caCertPool,
		ServerName:         opts.ServerName,
		InsecureSkipVerify: opts.Insecure,
	}
	tlsConfig.BuildNameToCertificate()

	return tlsConfig, nil
}
