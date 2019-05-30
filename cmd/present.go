package cmd

import (
	"context"
	"encoding/json"
	"log"

	"github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(presentCmd)
}

var presentCmd = &cobra.Command{
  Use:   "present [acme-hostname] [TXT-record]",
  Short: "Update the dns record",
  Long:  `Update ACME TXT record on the specified domain`,
  DisableFlagsInUseLine: true,
  Args: cobra.ExactArgs(2),
  RunE: func(cmd *cobra.Command, args []string) error {
  	j := &Service{
        TTL:	30,
        Text:	args[1],
    }
    val, err := json.Marshal(j)
    if err != nil {
    	return err
    }

    var key = etcdKeyFor(args[0])
    var value = string(val)
    log.Println("etcd put", key, value)

    ctx, cancel := context.WithTimeout(context.Background(), globalFlags.CommandTimeOut)
	_, err = client.Put(ctx, key, value)
	cancel()
	return err
  },
}