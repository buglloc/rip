package commands

import (
	"fmt"
	"os"

	log "github.com/buglloc/simplelog"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/buglloc/rip/pkg/cli"
)

var (
	RootCmd = &cobra.Command{
		Use:          "rip",
		Short:        "Wildcard DNS",
		SilenceUsage: false,
		PreRunE:      parseRootConfig,
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	flags := RootCmd.PersistentFlags()
	flags.String("config", "",
		"config file (default is $HOME/.rip.toml)")
	flags.Bool("verbose", false,
		"verbose output")

	viper.AutomaticEnv()
	_ = cli.BindPFlags(flags)
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

func parseRootConfig(cmd *cobra.Command, args []string) error {
	if viper.GetBool("Verbose") {
		log.SetLevel(log.DebugLevel)
	}
	return nil
}
