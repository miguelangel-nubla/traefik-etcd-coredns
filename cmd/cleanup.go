package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cleanupCmd)
}

var cleanupCmd = &cobra.Command{
	Use:                   "cleanup [acme-hostname] [TXT-record]",
	Short:                 "Delete the dns record",
	Long:                  `Delete ACME TXT record on the specified domain`,
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var key = etcdKeyFor(args[0])
		log.Println("etcd del", key)

		ctx, cancel := context.WithTimeout(context.Background(), globalFlags.CommandTimeOut)
		_, err := client.Delete(ctx, key)
		cancel()
		return err
	},
}
