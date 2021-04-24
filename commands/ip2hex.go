package commands

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/spf13/cobra"

	"github.com/buglloc/rip/pkg/iputil"
)

var ip2Hex = &cobra.Command{
	Use:   "ip2hex [IP] IP",
	Short: "Convert IPs to base-16 form",
	RunE:  runIp2HexCmd,
}

func init() {
	RootCmd.AddCommand(ip2Hex)
}

func runIp2HexCmd(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("please provide IP")
	}

	results := make([]string, len(args))
	for i, ip := range args {
		if strings.Contains(ip, ":") {
			results[i] = iputil.Ip6ToHex(net.ParseIP(ip))
		} else {
			results[i] = iputil.Ip4ToHex(net.ParseIP(ip))
		}
	}

	fmt.Println(strings.Join(results, "\t"))
	return nil
}
