package log

import (
	"fmt"
	"os"

	"github.com/juju/loggo"
)

// TailLoggoLogger implements the logger interface of the `tail` package,
//  by calling the loggo.Logger internally
type TailLoggoLogger struct {
	loggo.Logger
}

func (l TailLoggoLogger) Fatal(v ...interface{}) {
	l.Criticalf(v[0].(string), v[1:]...)
	os.Exit(1)
}
func (l TailLoggoLogger) Fatalf(format string, v ...interface{}) {
	l.Criticalf(format, v...)
	os.Exit(1)
}
func (l TailLoggoLogger) Fatalln(v ...interface{}) {
	l.Criticalf(v[0].(string), v[1:]...)
	os.Exit(1)
}
func (l TailLoggoLogger) Panic(v ...interface{}) {
	l.Criticalf(v[0].(string), v[1:]...)
	panic(fmt.Sprintf(v[0].(string), v[1:]...))
}
func (l TailLoggoLogger) Panicf(format string, v ...interface{}) {
	l.Criticalf(format, v...)
	panic(fmt.Sprintf(format, v...))
}
func (l TailLoggoLogger) Panicln(v ...interface{}) {
	l.Criticalf(v[0].(string), v[1:]...)
	panic(fmt.Sprintf(v[0].(string), v[1:]...))
}
func (l TailLoggoLogger) Print(v ...interface{}) {
	l.Infof(v[0].(string), v[1:]...)
}
func (l TailLoggoLogger) Printf(format string, v ...interface{}) {
	l.Infof(format, v...)
}
func (l TailLoggoLogger) Println(v ...interface{}) {
	l.Infof(v[0].(string), v[1:]...)
}
