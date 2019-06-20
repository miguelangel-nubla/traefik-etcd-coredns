package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/miguelangel-nubla/traefik-etcd-coredns/client/spec"
)

func init() {
	presentCmd.Flags().SortFlags = false
	presentCmd.Flags().SetInterspersed(false)

	rootCmd.AddCommand(presentCmd)
}

var presentCmd = &cobra.Command{
	Use:                   "present [acme-hostname] [TXT-record]",
	Short:                 "Update the dns record",
	Long:                  `Update ACME TXT record on the specified domain`,
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var dnsName = args[0]
		var dnsTxt = args[1]

		r := spec.Record{
			DNSName: dnsName,
			TTL:     30,
			Text:    dnsTxt,
		}

		err := cli.Update(r)
		if err != nil {
			return err
		}

		if len(configGlobal.UpdateDNSHost) > 0 {
			r := spec.Record{
				DNSName: strings.TrimPrefix(dnsName, ACMEChallengePrefix+"."),
				Host:    configGlobal.UpdateDNSHost,
				TTL:     configGlobal.UpdateDNSHostTTL,
			}
			return cli.Update(r)
		}

		return nil
	},
}
