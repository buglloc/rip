package log

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/inconshreveable/log15.v2"
)

const timeKey = "t"
const lvlKey = "lvl"
const msgKey = "msg"
const errorKey = "LOG15_ERROR"

const (
	timeFormat     = "2006-01-02T15:04:05-0700"
	termTimeFormat = "15:04:05"
	floatFormat    = 'f'
	termMsgJust    = 40
)

var colored = true

func init() {
	colored = terminal.IsTerminal(int(os.Stderr.Fd()))
}

func TextFormat() log15.Format {
	return log15.FormatFunc(func(r *log15.Record) []byte {
		var color = 0
		if colored {
			switch r.Lvl {
			case log15.LvlCrit:
				color = 35
			case log15.LvlError:
				color = 31
			case log15.LvlWarn:
				color = 33
			case log15.LvlInfo:
				color = 32
			case log15.LvlDebug:
				color = 36
			}
		}

		b := &bytes.Buffer{}
		lvl := strings.ToUpper(r.Lvl.String())
		if color > 0 {
			fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%s] %s ", color, lvl, r.Time.Format(termTimeFormat), r.Msg)
		} else {
			fmt.Fprintf(b, "[%s] [%s] %s ", lvl, r.Time.Format(termTimeFormat), r.Msg)
		}

		// try to justify the log output for short messages
		if len(r.Ctx) > 0 && len(r.Msg) < termMsgJust {
			b.Write(bytes.Repeat([]byte{' '}, termMsgJust-len(r.Msg)))
		}

		// print the keys logfmt style
		logfmt(b, r.Ctx, color)
		return b.Bytes()
	})
}

func logfmt(buf *bytes.Buffer, ctx []interface{}, color int) {
	for i := 0; i < len(ctx); i += 2 {
		if i != 0 {
			buf.WriteByte(' ')
		}

		k, ok := ctx[i].(string)
		v := formatLogfmtValue(ctx[i+1])
		if !ok {
			k, v = errorKey, formatLogfmtValue(k)
		}

		// XXX: we should probably check that all of your key bytes aren't invalid
		if color > 0 {
			fmt.Fprintf(buf, "\x1b[%dm%s\x1b[0m=%s", color, k, v)
		} else {
			fmt.Fprintf(buf, "%s=%s", k, v)
		}
	}

	buf.WriteByte('\n')
}

func needsQuoting(text string) bool {
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func formatLogfmtValue(value interface{}) string {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if needsQuoting(stringVal) {
		return fmt.Sprintf("%q", stringVal)
	}

	return stringVal
}
