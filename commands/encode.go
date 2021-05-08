package commands

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/spf13/cobra"

	"github.com/buglloc/rip/v2/pkg/iputil"
)

var ip2Hex = &cobra.Command{
	Use:   "encode [IP] IP",
	Short: "Encode IPs",
	RunE: func(_ *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide IP")
		}

		results := make([]string, len(args))
		for i, ip := range args {
			if strings.Contains(ip, ":") {
				results[i] = iputil.EncodeIP6(net.ParseIP(ip))
			} else {
				results[i] = iputil.EncodeIP4(net.ParseIP(ip))
			}
		}

		fmt.Println(strings.Join(results, "\t"))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(ip2Hex)
}
