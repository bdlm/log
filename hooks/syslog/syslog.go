// +build !windows,!nacl,!plan9

package syslog

import (
	"fmt"
	"log/syslog"
	"os"

	"github.com/bdlm/log/v2"
	stdLogger "github.com/bdlm/std/v2/logger"
)

// Hook to send logs via syslog.
type Hook struct {
	Writer        *syslog.Writer
	SyslogNetwork string
	SyslogRaddr   string
}

// NewHook creates a hook to be added to an instance of logger. This is called
// with `hook, err := NewHook("udp", "localhost:514", syslog.LOG_DEBUG, "")`
// `if err == nil { log.Hooks.Add(hook) }`
func NewHook(network, raddr string, priority syslog.Priority, tag string) (*Hook, error) {
	w, err := syslog.Dial(network, raddr, priority, tag)
	return &Hook{Writer: w, SyslogNetwork: network, SyslogRaddr: raddr}, err
}

// Fire executes the syslog hook.
func (hook *Hook) Fire(entry *log.Entry) error {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	switch entry.Level {
	case log.PanicLevel:
		return hook.Writer.Crit(line)
	case log.FatalLevel:
		return hook.Writer.Crit(line)
	case log.ErrorLevel:
		return hook.Writer.Err(line)
	case log.WarnLevel:
		return hook.Writer.Warning(line)
	case log.InfoLevel:
		return hook.Writer.Info(line)
	case log.DebugLevel:
		return hook.Writer.Debug(line)
	default:
		return nil
	}
}

// Levels returns all available log levels.
func (hook *Hook) Levels() []stdLogger.Level {
	return log.AllLevels
}
