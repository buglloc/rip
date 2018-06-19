package commands

import (
	"errors"
	"net"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/buglloc/rip/pkg/cfg"
	"github.com/buglloc/rip/pkg/cli"
	"github.com/buglloc/rip/pkg/dns_server"
)

var serverCmd = &cobra.Command{
	Use:     "server --zone=example.com --zone=example1.com",
	Short:   "Start RIP server",
	PreRunE: parseServerConfig,
	RunE:    runServerCmd,
}

func init() {
	flags := serverCmd.PersistentFlags()
	flags.String("listen", ":53",
		"address to listen on")
	flags.StringSlice("zone", []string{"."},
		"your zone name (e.g. 'buglloc.com')")
	flags.String("ipv4", "127.0.0.1",
		"default ipv4 address")
	flags.String("ipv6", "::1",
		"default ipv6 address")
	flags.String("upstream", "77.88.8.8:53",
		"upstream DNS server")
	flags.Bool("strict", false,
		"don't return default IPs for not supported requests")
	flags.Bool("no-proxy", false,
		"disable proxy mode")

	cli.BindPFlags(flags)
	RootCmd.AddCommand(serverCmd)
}

func runServerCmd(cmd *cobra.Command, args []string) error {
	dns_server.RunBackground()
	cli.ListenInterrupt()
	return nil
}

func parseServerConfig(cmd *cobra.Command, args []string) error {
	cfg.Zones = viper.GetStringSlice("Zone")
	if len(cfg.Zones) == 0 {
		return errors.New("empty zone list, please provide at leas one")
	}

	cfg.Addr = viper.GetString("Listen")
	cfg.IPv4 = net.ParseIP(viper.GetString("Ipv4"))
	cfg.IPv6 = net.ParseIP(viper.GetString("Ipv6"))
	cfg.AllowProxy = !viper.GetBool("NoProxy")
	cfg.StrictMode = viper.GetBool("Strict")
	cfg.Upstream = viper.GetString("Upstream")
	return nil
}
