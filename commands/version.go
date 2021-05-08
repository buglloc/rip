package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/buglloc/rip/v2/pkg/cfg"
)

var version = &cobra.Command{
	Use:   "version",
	Short: "Print rip version",
	RunE: func(_ *cobra.Command, _ []string) error {
		fmt.Printf("RIP v%s\n", cfg.Version)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(version)
}
