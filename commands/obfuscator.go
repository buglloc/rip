package commands

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	obfuscator "github.com/buglloc/rip/pkg/obfustacor"
)

var uglify = &cobra.Command{
	Use:   "uglify IP",
	Short: "Obfuscate IP",
	RunE:  runUglifyCmd,
}

func init() {
	RootCmd.AddCommand(uglify)
}

func runUglifyCmd(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("please provide IP")
	}

	obfuscated := obfuscator.IPv4(args[0])
	for _, r := range obfuscated {
		fmt.Printf("http://%s\n", r)
	}
	return nil
}
