package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/buglloc/rip/pkg/cfg"
)

var version = &cobra.Command{
	Use:   "version",
	Short: "Print rip version",
	RunE:  versionCmd,
}

func init() {
	RootCmd.AddCommand(version)
}

func versionCmd(cmd *cobra.Command, _ []string) error {
	fmt.Printf("RIP v%s\n", cfg.Version)
	return nil
}
