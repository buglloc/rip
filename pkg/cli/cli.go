package cli

import (
	"os"
	"os/signal"
	"unicode"

	"github.com/buglloc/simplelog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func ListenInterrupt() {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	for {
		select {
		case <-sig:
			log.Info("Terminating...")
			return
		}
	}
}

func BindPFlags(flags *pflag.FlagSet) (err error) {
	flags.VisitAll(func(flag *pflag.Flag) {
		if err = viper.BindPFlag(transformFlagName(flag.Name), flag); err != nil {
			return
		}
	})
	return
}

func transformFlagName(name string) string {
	runes := []rune(name)
	length := len(runes)

	var out []rune
	nextUpper := true
	for i := 0; i < length; i++ {
		switch {
		case nextUpper:
			out = append(out, unicode.ToUpper(runes[i]))
			nextUpper = false
		case runes[i] == '-':
			nextUpper = true
		default:
			out = append(out, runes[i])
		}
	}

	return string(out)
}
