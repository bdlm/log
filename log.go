package log

import (
	"fmt"
	"log"
	"strings"

	stdLogger "github.com/bdlm/std/logger"
)

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func levelString(level stdLogger.Level) string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	}

	return "unknown"
}

// ParseLevel takes a string level and returns the log level constant.
func ParseLevel(lvl string) (stdLogger.Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	}

	var l stdLogger.Level
	return l, fmt.Errorf("not a valid log Level: %q", lvl)
}

// AllLevels is a constant exposing all logging levels.
var AllLevels = []stdLogger.Level{
	stdLogger.Panic,
	stdLogger.Fatal,
	stdLogger.Error,
	stdLogger.Warn,
	stdLogger.Info,
}

// These are the standard logging levels. You can set the logging level to log
// on your instance of logger, obtained with `New()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel = stdLogger.Panic
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel = stdLogger.Fatal
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel = stdLogger.Error
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel = stdLogger.Warn
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel = stdLogger.Info
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel = stdLogger.Debug
)

// Won't compile if StdLogger can't be realized by a log.Logger
var (
	_ StdLogger = &log.Logger{}
	_ StdLogger = &Entry{}
	_ StdLogger = &Logger{}
)

// StdLogger is what your bdlm/log-enabled library should take, that way
// it'll accept a stdlib logger and a bdlm/log logger. There's no standard
// interface, this is the closest we get, unfortunately.
type StdLogger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}

// The FieldLogger interface generalizes the Entry and Logger types
type FieldLogger interface {
	WithField(key string, value interface{}) *Entry
	WithFields(fields Fields) *Entry
	WithError(err error) *Entry

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
}
