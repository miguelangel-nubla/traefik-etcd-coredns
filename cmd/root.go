package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/miguelangel-nubla/traefik-etcd-coredns/client/etcd"
	"github.com/miguelangel-nubla/traefik-etcd-coredns/client/etcd/coredns"

	"github.com/miguelangel-nubla/traefik-etcd-coredns/client/spec"
)

type Config struct {
	Debug            bool
	UpdateDNSHost    string
	UpdateDNSHostTTL uint32
	Backend          string
}

const (
	cliName        = "traefik-etcd-coredns"
	cliDescription = "Traefik custom acme.dnsChallenge provider for CoreDNS servers with Etcd backend."

	defaultDialTimeout    = 2 * time.Second
	defaultCommandTimeOut = 5 * time.Second

	ACMEChallengePrefix = "_acme-challenge"
)

var configGlobal = Config{}
var configEtcd = etcd.Config{}

var coreDNS = &coredns.Client{
	Client: etcd.Client{Config: &configEtcd},
}

var (
	cli spec.Client

	rootCmd = &cobra.Command{
		Use:   cliName,
		Short: cliDescription,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if configGlobal.Backend == "etcd-coredns" {
				configEtcd.Debug = configGlobal.Debug
				cli = coreDNS
			} else {
				log.Fatalf("Unknown backend %s", configGlobal.Backend)
			}

			cli.Init()
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			cli.Close()
		},
	}
)

func init() {
	log.SetPrefix("[" + cliName + "] ")

	rootCmd.PersistentFlags().BoolVar(&configGlobal.Debug, "debug", false, "enable client-side debug logging")

	rootCmd.PersistentFlags().StringVar(&configGlobal.UpdateDNSHost, "update-dns-host", "", "hostname or IP to also update A or CNAME DNS records")
	rootCmd.PersistentFlags().Uint32Var(&configGlobal.UpdateDNSHostTTL, "update-dns-host-ttl", 60, "TTL for records created by --update-dns-host")

	rootCmd.PersistentFlags().StringVar(&configGlobal.Backend, "backend", "etcd-coredns", "select backend, currently only etcd-coredns available")

	rootCmd.PersistentFlags().StringSliceVar(&configEtcd.Endpoints, "etcd-endpoints", []string{"127.0.0.1:2379"}, "gRPC endpoints")

	rootCmd.PersistentFlags().DurationVar(&configEtcd.DialTimeout, "etcd-dial-timeout", defaultDialTimeout, "dial timeout for client connections")
	rootCmd.PersistentFlags().DurationVar(&configEtcd.CommandTimeOut, "etcd-command-timeout", defaultCommandTimeOut, "timeout for short running command (excluding dial timeout)")

	rootCmd.PersistentFlags().BoolVar(&configEtcd.InsecureTransport, "etcd-insecure-transport", true, "disable transport security for client connections")
	rootCmd.PersistentFlags().BoolVar(&configEtcd.TLS.InsecureSkipVerify, "etcd-insecure-skip-tls-verify", false, "skip server certificate verification")
	rootCmd.PersistentFlags().StringVar(&configEtcd.TLS.CertFile, "etcd-cert", "", "identify secure client using this TLS certificate file")
	rootCmd.PersistentFlags().StringVar(&configEtcd.TLS.KeyFile, "etcd-key", "", "identify secure client using this TLS key file")
	rootCmd.PersistentFlags().StringVar(&configEtcd.TLS.TrustedCAFile, "etcd-cacert", "", "verify certificates of TLS-enabled secure servers using this CA bundle")
	rootCmd.PersistentFlags().StringVar(&configEtcd.User, "etcd-user", "", "username[:password] for authentication")
	rootCmd.PersistentFlags().StringVar(&configEtcd.Password, "etcd-password", "", "password for authentication (if this option is used, --user option shouldn't include password)")

	rootCmd.PersistentFlags().StringVar(&coreDNS.CoreDNSPrefix, "etcd-coredns-prefix", "/skydns", "etcd key prefix for coredns records")

	rootCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		env := strings.ToUpper(strings.Replace(cliName+"_"+flag.Name, "-", "_", -1))
		flag.Usage = fmt.Sprintf("[env %v] %v", env, flag.Usage)
		if value := os.Getenv(env); value != "" {
			flag.Value.Set(value)
		}
	})

	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false

	cobra.EnableCommandSorting = false
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		//fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
