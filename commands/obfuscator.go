package commands

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/buglloc/rip/pkg/ip_obfustacor"
)

var uglify = &cobra.Command{
	Use:   "uglify IP",
	Short: "Obfuscate IP",
	RunE:  runUglify,
}

func init() {
	RootCmd.AddCommand(uglify)
}

func runUglify(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("please provide IP")
	}

	obfuscated := ip_obfuscator.IPv4(args[0])
	for _, r := range obfuscated {
		fmt.Printf("http://%s\n", r)
	}
	return nil
}
