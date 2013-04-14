package mongowebstat

import (
	"log"
	"os"
)

var l *log.Logger

var debug debugging = false

type debugging bool

func (d debugging) Printf(format string, args ...interface{}) {
	if d {
		l.Printf(format, args...)
	}
}

func (d debugging) Print(args ...interface{}) {
	if d {
		l.Print(args...)
	}
}

func init() {
	out := os.Stdout
	l = log.New(out, "", 3)
}
