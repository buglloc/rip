package commands

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/spf13/cobra"

	"github.com/buglloc/rip/pkg/ip_utils"
)

var ip2Hex = &cobra.Command{
	Use:   "ip2hex IP",
	Short: "Start RIP server",
	RunE:  runIp2Hex,
}

func init() {
	RootCmd.AddCommand(ip2Hex)
}

func runIp2Hex(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("please provide IP")
	}

	results := make([]string, len(args))
	for i, ip := range args {
		if strings.Contains(ip, ":") {
			results[i] = ip_utils.Ip6ToHex(net.ParseIP(ip))
		} else {
			results[i] = ip_utils.Ip4ToHex(net.ParseIP(ip))
		}
	}

	fmt.Println(strings.Join(results, "\t"))
	return nil
}
