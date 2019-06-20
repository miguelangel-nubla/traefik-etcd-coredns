package coredns

import (
	"net"
	"strings"

	"github.com/miguelangel-nubla/traefik-etcd-coredns/client/spec"
)

func dnsNameFor(etcdKey string, coreDNSprefix string) string {
	domains := strings.Split(strings.TrimPrefix(etcdKey, coreDNSprefix), "/")
	reverse(domains)
	return strings.Join(domains, ".")
}

func etcdKeyFor(dnsName string, coreDNSprefix string) string {
	domains := strings.Split(dnsName, ".")
	reverse(domains)
	return "/" + strings.Trim(coreDNSprefix, "/") + "/" + strings.Join(domains, "/") + "/x1"
}

func reverse(slice []string) {
	for i := 0; i < len(slice)/2; i++ {
		j := len(slice) - i - 1
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func recordType(target string) string {
	if net.ParseIP(target) != nil {
		return spec.RecordTypeA
	}
	return spec.RecordTypeCNAME
}
