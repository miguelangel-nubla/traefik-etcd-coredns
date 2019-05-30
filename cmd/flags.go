package cmd

import (
	"time"

	"go.etcd.io/etcd/pkg/transport"
)

type GlobalFlags struct {
	CoreDNSPrefix string

	InsecureTransport bool
	Endpoints         []string
	DialTimeout       time.Duration
	CommandTimeOut    time.Duration

	TLS transport.TLSInfo

	User     string
	Password string

	Debug bool
}

var (
	globalFlags = GlobalFlags{}
)
