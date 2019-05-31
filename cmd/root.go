package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"

	"google.golang.org/grpc/grpclog"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	cliName        = "traefik-etcd-coredns"
	cliDescription = "Traefik custom acme.dnsChallenge provider for CoreDNS servers with etcd backend"

	defaultDialTimeout    = 2 * time.Second
	defaultCommandTimeOut = 5 * time.Second
)

var client *clientv3.Client

var (
	rootCmd = &cobra.Command{
		Use:   cliName,
		Short: cliDescription,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initClient()
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			client.Close()
		},
	}
)

func init() {
	log.SetPrefix("[" + cliName + "] ")

	cobra.EnableCommandSorting = false
	rootCmd.PersistentFlags().SortFlags = false

	rootCmd.PersistentFlags().StringVar(&globalFlags.CoreDNSPrefix, "prefix", "/skydns", "etcd key prefix")

	rootCmd.PersistentFlags().StringSliceVar(&globalFlags.Endpoints, "endpoints", []string{"127.0.0.1:2379"}, "gRPC endpoints")
	rootCmd.PersistentFlags().BoolVar(&globalFlags.Debug, "debug", false, "enable client-side debug logging")

	rootCmd.PersistentFlags().DurationVar(&globalFlags.DialTimeout, "dial-timeout", defaultDialTimeout, "dial timeout for client connections")
	rootCmd.PersistentFlags().DurationVar(&globalFlags.CommandTimeOut, "command-timeout", defaultCommandTimeOut, "timeout for short running command (excluding dial timeout)")

	rootCmd.PersistentFlags().BoolVar(&globalFlags.InsecureTransport, "insecure-transport", true, "disable transport security for client connections")
	rootCmd.PersistentFlags().BoolVar(&globalFlags.TLS.InsecureSkipVerify, "insecure-skip-tls-verify", false, "skip server certificate verification")
	rootCmd.PersistentFlags().StringVar(&globalFlags.TLS.CertFile, "cert", "", "identify secure client using this TLS certificate file")
	rootCmd.PersistentFlags().StringVar(&globalFlags.TLS.KeyFile, "key", "", "identify secure client using this TLS key file")
	rootCmd.PersistentFlags().StringVar(&globalFlags.TLS.TrustedCAFile, "cacert", "", "verify certificates of TLS-enabled secure servers using this CA bundle")
	rootCmd.PersistentFlags().StringVar(&globalFlags.User, "user", "", "username[:password] for authentication")
	rootCmd.PersistentFlags().StringVar(&globalFlags.Password, "password", "", "password for authentication (if this option is used, --user option shouldn't include password)")

	rootCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		env := strings.ToUpper(strings.Replace(cliName+"_"+flag.Name, "-", "_", -1))
		flag.Usage = fmt.Sprintf("[env %v] %v", env, flag.Usage)
		if value := os.Getenv(env); value != "" {
			flag.Value.Set(value)
		}
	})
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		//fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func initClient() {
	tlsConfig, err := globalFlags.TLS.ClientConfig()
	if err != nil {
		log.Fatal(err)
	}

	var config = clientv3.Config{
		Endpoints:   globalFlags.Endpoints,
		DialTimeout: globalFlags.DialTimeout,
		Username:    globalFlags.User,
		Password:    globalFlags.Password,
	}

	if !globalFlags.InsecureTransport {
		config.TLS = tlsConfig
	}

	if globalFlags.Debug {
		clientv3.SetLogger(grpclog.NewLoggerV2(os.Stderr, os.Stderr, os.Stderr))
	}

	client, err = clientv3.New(config)
	if err != nil {
		log.Fatal(err)
	}
}
