package etcd

import (
	"os"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"

	"google.golang.org/grpc/grpclog"
)

type Client struct {
	Client *clientv3.Client
	Config *Config
}

type Config struct {
	InsecureTransport bool
	Endpoints         []string
	DialTimeout       time.Duration
	CommandTimeOut    time.Duration

	TLS transport.TLSInfo

	User     string
	Password string

	Debug bool
}

func (c *Client) Init() error {
	tlsConfig, err := c.Config.TLS.ClientConfig()
	if err != nil {
		return err
	}

	var config = clientv3.Config{
		Endpoints:   c.Config.Endpoints,
		DialTimeout: c.Config.DialTimeout,
		Username:    c.Config.User,
		Password:    c.Config.Password,
	}

	if !c.Config.InsecureTransport {
		config.TLS = tlsConfig
	}

	if c.Config.Debug {
		clientv3.SetLogger(grpclog.NewLoggerV2(os.Stderr, os.Stderr, os.Stderr))
	}

	c.Client, err = clientv3.New(config)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}
