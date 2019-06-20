package coredns

import (
	"context"
	"encoding/json"
	"log"

	"github.com/miguelangel-nubla/traefik-etcd-coredns/client/etcd"

	"github.com/miguelangel-nubla/traefik-etcd-coredns/client/spec"
)

type Client struct {
	etcd.Client
	CoreDNSPrefix string
}

func (c *Client) GetName() string {
	return "etcd-coredns"
}

func (c *Client) Update(r spec.Record) error {
	val, err := json.Marshal(r)
	if err != nil {
		return err
	}

	var key = etcdKeyFor(r.DNSName, c.CoreDNSPrefix)
	var value = string(val)
	log.Println("etcd put", key, value)

	ctx, cancel := context.WithTimeout(context.Background(), c.Config.CommandTimeOut)
	_, err = c.Client.Client.Put(ctx, key, value)
	cancel()
	return err
}

func (c *Client) Delete(r spec.Record) error {
	var key = etcdKeyFor(r.DNSName, c.CoreDNSPrefix)
	log.Println("etcd del", key)

	ctx, cancel := context.WithTimeout(context.Background(), c.Config.CommandTimeOut)
	_, err := c.Client.Client.Delete(ctx, key)
	cancel()
	return err
}
