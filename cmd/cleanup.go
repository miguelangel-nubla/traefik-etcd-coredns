package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/miguelangel-nubla/traefik-etcd-coredns/client/spec"
)

func init() {
	cleanupCmd.Flags().SortFlags = false
	cleanupCmd.Flags().SetInterspersed(false)

	rootCmd.AddCommand(cleanupCmd)
}

var cleanupCmd = &cobra.Command{
	Use:                   "cleanup [acme-hostname] [TXT-record]",
	Short:                 "Delete the dns record",
	Long:                  `Delete ACME TXT record on the specified domain`,
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var dnsName = args[0]

		r := spec.Record{
			DNSName: dnsName,
		}

		err := cli.Delete(r)
		if err != nil {
			return err
		}

		if len(configGlobal.UpdateDNSHost) > 0 {
			r := spec.Record{
				DNSName: strings.TrimPrefix(dnsName, ACMEChallengePrefix+"."),
			}
			return cli.Delete(r)
		}

		return nil
	},
}
