package commands

import (
	"fmt"
	"net"
	"os"

	"github.com/buglloc/simplelog"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/buglloc/rip/pkg/cfg"
	"github.com/buglloc/rip/pkg/cli"
	"github.com/buglloc/rip/pkg/dns_server"
)

var (
	RootCmd = &cobra.Command{
		Use:          "dip --zone=example.com --zone=example1.com",
		Short:        "Wildcard DNS",
		SilenceUsage: false,
		PreRunE:      parseConfig,
		RunE:         runRootCmd,
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	flags := RootCmd.PersistentFlags()
	flags.String("config", "",
		"config file (default is $HOME/.rip.toml)")
	flags.String("listen", ":53",
		"address to listen on")
	flags.StringSlice("zone", []string{"."},
		"your zone name (e.g. 'buglloc.com')")
	flags.String("ipv4", "127.0.0.1",
		"default ipv4 address")
	flags.String("ipv6", "::1",
		"default ipv6 address")
	flags.Bool("verbose", false,
		"verbose output")

	viper.AutomaticEnv()
	cli.BindPFlags(flags)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile := viper.GetString("config"); cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".rip")
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
func runRootCmd(cmd *cobra.Command, args []string) error {
	dns_server.RunBackground()
	cli.ListenInterrupt()
	return nil
}

func parseConfig(cmd *cobra.Command, args []string) error {
	if viper.GetBool("verbose") {
		log.SetLevel(log.DebugLevel)
	}

	cfg.Zones = viper.GetStringSlice("zone")
	if len(cfg.Zones) == 0 {
		return errors.New("Empty zone list, please provide at leas one")
	}

	cfg.Addr = viper.GetString("listen")
	cfg.IPv4 = net.ParseIP(viper.GetString("ipv4"))
	cfg.IPv6 = net.ParseIP(viper.GetString("ipv6"))
	return nil
}
