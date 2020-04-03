package wlog

import (
	"fmt"
	"log"
	"syscall/js"
	"time"

	"github.com/pkg/errors"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Raw(args ...interface{})
}

var logger Logger

func init() {
	logger = &DefaultLogger{}
}

type DefaultLogger struct {
}

func (l *DefaultLogger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO]  :"+format, args...)
}

func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR] :"+format, args...)
}

func (l *DefaultLogger) Raw(args ...interface{}) {
	log.Println(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Raw(args ...interface{}) {
	logger.Raw(args...)
}

type ConsoleLogger struct {
	jsCon js.Value
}

func InstallConsoleLogger() error {
	jsCon := js.Global().Get("console")
	if jsCon.Type() == js.TypeUndefined || jsCon.Type() == js.TypeNull {
		return errors.Errorf("failed to get console")
	}
	logger = &ConsoleLogger{
		jsCon: jsCon,
	}
	return nil
}

func formatNow() string {
	return time.Now().Round(1 * time.Millisecond).Format("2006-01-02T15:04:05.000")
}

func (l *ConsoleLogger) Infof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.jsCon.Call("log", formatNow()+" [INFO] "+msg)
}

func (l *ConsoleLogger) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.jsCon.Call("log", "%c"+formatNow()+" [ERROR] "+msg, "color: red;")
}

func (l *ConsoleLogger) Raw(args ...interface{}) {
	l.jsCon.Call("log", args...)
}
