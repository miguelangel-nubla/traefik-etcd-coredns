# traefik-etcd-coredns

DNS provider for [Traefik Let's Encrypt support](https://docs.traefik.io/configuration/acme/#provider) that enables the use of wildcard domains.
Allows using a existing [CoreDNS with Etcd](https://coredns.io/plugins/etcd/) deployment to validate `DNS-01` challenges from [Let's Encrypt](https://letsencrypt.org/docs/challenge-types/#dns-01-challenge) and get the wildcard certificate issued and configured.


# Quickstart

* Build or download the binary:
 `curl -sL https://github.com/miguelangel-nubla/traefik-etcd-coredns/releases/latest/download/traefik-etcd-coredns_linux_amd64.tar.gz | tar zxvf - -C /path/to/`

* Add to your Traefik configuration file:
```
[acme]
acmeLogging = true # Recommended for debugging
email = "your-email-here@my-awesome-app.org"
storage = "acme.json"
entryPoint = "https"
onHostRule = true
[acme.dnsChallenge]
provider = "exec"
delayBeforeCheck = 30
```
* Provide the required environment variables:
```
EXEC_PATH=/path/to/traefik-etcd-coredns
TRAEFIK_ETCD_COREDNS_ENDPOINTS=[etcd-cluster-address]:[port]
TRAEFIK_ETCD_COREDNS_DEBUG=true # Recommended for debugging
```
* Done! Traefik should now be creating and submitting for validation `_acme-challenge.<DOMAIN>` records.

## Configuration
Depending on your Etcd connection requirements you can pass the corresponding environment variables:
```
$ ./traefik-etcd-coredns --help
Traefik custom acme.dnsChallenge provider for CoreDNS servers with Etcd backend.

Usage:
  traefik-etcd-coredns [command]

Available Commands:
  cleanup     Delete the dns record
  present     Update the dns record
  help        Help about any command

Flags:
      --endpoints strings          [env TRAEFIK_ETCD_COREDNS_ENDPOINTS] gRPC endpoints (default [127.0.0.1:2379])
      --debug                      [env TRAEFIK_ETCD_COREDNS_DEBUG] enable client-side debug logging
      --prefix string              [env TRAEFIK_ETCD_COREDNS_PREFIX] etcd key prefix (default "/skydns")
      --dial-timeout duration      [env TRAEFIK_ETCD_COREDNS_DIAL_TIMEOUT] dial timeout for client connections (default 2s)
      --command-timeout duration   [env TRAEFIK_ETCD_COREDNS_COMMAND_TIMEOUT] timeout for short running command (excluding dial timeout) (default 5s)
      --insecure-transport         [env TRAEFIK_ETCD_COREDNS_INSECURE_TRANSPORT] disable transport security for client connections (default true)
      --insecure-skip-tls-verify   [env TRAEFIK_ETCD_COREDNS_INSECURE_SKIP_TLS_VERIFY] skip server certificate verification
      --cert string                [env TRAEFIK_ETCD_COREDNS_CERT] identify secure client using this TLS certificate file
      --key string                 [env TRAEFIK_ETCD_COREDNS_KEY] identify secure client using this TLS key file
      --cacert string              [env TRAEFIK_ETCD_COREDNS_CACERT] verify certificates of TLS-enabled secure servers using this CA bundle
      --user string                [env TRAEFIK_ETCD_COREDNS_USER] username[:password] for authentication
      --password string            [env TRAEFIK_ETCD_COREDNS_PASSWORD] password for authentication (if this option is used, --user option shouldn't include password)
  -h, --help                       help for traefik-etcd-coredns

Use "traefik-etcd-coredns [command] --help" for more information about a command.
```