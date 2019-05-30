package cmd

import (
	"strings"
)

type Service struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Weight   int    `json:"weight,omitempty"`
	Text     string `json:"text,omitempty"`
	Mail     bool   `json:"mail,omitempty"` // Be an MX record. Priority becomes Preference.
	TTL      uint32 `json:"ttl,omitempty"`
}

func dnsNameFor(etcdKey string) string {
	domains := strings.Split(strings.TrimPrefix(etcdKey, globalFlags.CoreDNSPrefix), "/")
	reverse(domains)
	return strings.Join(domains, ".")
}

func etcdKeyFor(dnsName string) string {
	domains := strings.Split(dnsName, ".")
	reverse(domains)
	return globalFlags.CoreDNSPrefix + strings.Join(domains, "/") + "/x1"
}

func reverse(slice []string) {
	for i := 0; i < len(slice)/2; i++ {
		j := len(slice) - i - 1
		slice[i], slice[j] = slice[j], slice[i]
	}
}