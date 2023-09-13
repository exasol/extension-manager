package extensionAPI

import (
	"github.com/dop251/goja_nodejs/console"
	log "github.com/sirupsen/logrus"
)

// createJavaScriptLogger creates a new [console.Printer] for handling `console.log()`, `console.warn()` and `console.error()` calls in JavaScript.
// This implementation forwards all log messages to logrus using the appropriate log method `Print()`, `Warn()` and `Error()`.
func createJavaScriptLogger(logPrefix string) console.Printer {
	return &jsLogger{logPrefix: logPrefix}
}

type jsLogger struct {
	logPrefix string
}

func (l jsLogger) Log(s string) {
	log.Print(l.logPrefix + s)
}

func (l jsLogger) Warn(s string) {
	log.Warn(l.logPrefix + s)
}

func (l jsLogger) Error(s string) {
	log.Error(l.logPrefix + s)
}
