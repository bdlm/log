package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bdlm/std/v2/logger"
)

var bufferPool *sync.Pool

func init() {
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
}

// ErrorKey defines the key when adding errors using WithError.
var ErrorKey = "error"

// Entry is the final or intermediate logging entry. It contains all the
// fields passed with WithField{,s}. It's finally logged when Debug, Info,
// Warn, Error, Fatal or Panic is called on it. These objects can be reused
// and passed around as much as you wish to avoid field duplication.
type Entry struct {
	Logger *Logger

	// When formatter is called in entry.log(), an Buffer may be set to entry
	Buffer *bytes.Buffer

	// Contains all the fields set by the user.
	Data Fields

	// Contains the error passed by WithError(error)
	Err error

	// Level the log entry was logged at: Debug, Info, Warn, Error, Fatal or Panic
	// This field will be set on entry firing and the value will be equal to the one in Logger struct field.
	Level logger.Level

	// Message passed to Debug, Info, Warn, Error, Fatal or Panic
	Message string

	// Time at which the log entry was created
	Time time.Time
}

var sanitizeStrings = []string{}

// NewEntry returns a new logger entry.
func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: logger,
		// Default is five fields, give a little extra room
		Data: make(Fields, 5),
	}
}

// AddSecret adds strings to the sanitization list.
func AddSecret(secrets ...string) {
	newStrings := []string{}
	for _, secret := range secrets {
		duplicate := false
		for _, str := range sanitizeStrings {
			if str == secret {
				duplicate = true
			}
		}
		if !duplicate {
			newStrings = append(newStrings, secret)
		}
	}
	sanitizeStrings = append(sanitizeStrings, newStrings...)
}

// String returns the string representation from the reader and ultimately the
// formatter.
func (entry *Entry) String() (string, error) {
	serialized, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		return "", err
	}
	str := string(serialized)
	return str, nil
}

// WithError add an error as single field (using the key defined in ErrorKey) to the Entry.
func (entry *Entry) WithError(err error) *Entry {
	return &Entry{
		Data:    entry.Data,
		Err:     err,
		Level:   entry.Level,
		Logger:  entry.Logger,
		Message: entry.Message,
		Time:    entry.Time,
	}
}

// WithField add a single field to the Entry.
func (entry *Entry) WithField(key string, value interface{}) *Entry {
	return entry.WithFields(Fields{key: value})
}

// WithFields adds a map of fields to the Entry.
func (entry *Entry) WithFields(fields Fields) *Entry {

	data := Fields{}
	for k, v := range entry.Data {
		data[k] = v
	}
	for k, v := range fields {
		data[k] = v
	}

	return &Entry{
		Data:    data,
		Err:     entry.Err,
		Level:   entry.Level,
		Logger:  entry.Logger,
		Message: entry.Message,
		Time:    entry.Time,
	}
}

// WithTime overrides the time of the Entry.
func (entry *Entry) WithTime(t time.Time) *Entry {
	return &Entry{Logger: entry.Logger, Data: entry.Data, Err: entry.Err, Time: t}
}

// This function is not declared with a pointer value because otherwise
// race conditions will occur when using multiple goroutines
func (entry Entry) log(level logger.Level, msg string) {
	if nil == entry.Logger.Out || entry.Logger.Out == ioutil.Discard {
		return
	}

	var buffer *bytes.Buffer

	// Default to now, but allow users to override if they want.
	//
	// We don't have to worry about polluting future calls to Entry#log()
	// with this assignment because this function is declared with a
	// non-pointer receiver.
	if entry.Time.IsZero() {
		entry.Time = time.Now()
	}

	entry.Level = level
	entry.Message = msg

	entry.fireHooks()

	buffer = bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)
	entry.Buffer = buffer

	entry.write()

	entry.Buffer = nil

	// To avoid Entry#log() returning a value that only would make sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	if level <= PanicLevel && level != FatalLevel {
		panic(&entry)
	}
}

// This function is not declared with a pointer value because otherwise
// race conditions will occur when using multiple goroutines
func (entry Entry) fireHooks() {
	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()
	err := entry.Logger.Hooks.Fire(entry.Level, &entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fire hook: %v\n", err)
	}
}

func (entry *Entry) write() {
	serialized, err := entry.Logger.Formatter.Format(entry)
	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v\n", err)
	} else {
		for _, secret := range sanitizeStrings {
			// Sanitize secrets
			serialized = []byte(strings.Replace(
				string(serialized),
				secret,
				"[REDACTED]",
				-1,
			))

			// Sanitize JSON-encoded secrets
			jsonSecret, _ := json.Marshal(secret)
			// Trim " from json.Marshal
			jsonSecret = jsonSecret[1 : len(jsonSecret)-1]
			serialized = []byte(strings.Replace(
				string(serialized),
				string(jsonSecret),
				"[REDACTED]",
				-1,
			))
		}
		_, err = entry.Logger.Out.Write(serialized)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
		}
	}
}

// Debug logs a debug-level message using Println.
func (entry *Entry) Debug(args ...interface{}) {
	if entry.Logger.level() >= DebugLevel {
		entry.log(DebugLevel, fmt.Sprint(args...))
	}
}

// Info logs a info-level message using Println.
func (entry *Entry) Info(args ...interface{}) {
	if entry.Logger.level() >= InfoLevel {
		entry.log(InfoLevel, fmt.Sprint(args...))
	}
}

// Print logs a info-level message using Println.
func (entry *Entry) Print(args ...interface{}) {
	entry.Info(args...)
}

// Warn logs a warn-level message using Println.
func (entry *Entry) Warn(args ...interface{}) {
	if entry.Logger.level() >= WarnLevel {
		entry.log(WarnLevel, fmt.Sprint(args...))
	}
}

// Warning logs a warn-level message using Println.
func (entry *Entry) Warning(args ...interface{}) {
	entry.Warn(args...)
}

// Error logs a error-level message using Println.
func (entry *Entry) Error(args ...interface{}) {
	if entry.Logger.level() >= ErrorLevel {
		entry.log(ErrorLevel, fmt.Sprint(args...))
	}
}

// Fatal logs a fatal-level message using Println.
func (entry *Entry) Fatal(args ...interface{}) {
	if entry.Logger.level() >= FatalLevel {
		entry.log(FatalLevel, fmt.Sprint(args...))
	}
	Exit(1)
}

// Panic logs a panic-level message using Println.
func (entry *Entry) Panic(args ...interface{}) {
	if entry.Logger.level() >= PanicLevel {
		entry.log(PanicLevel, fmt.Sprint(args...))
	}
	panic(fmt.Sprint(args...))
}

// Debugf logs a debug-level message using Printf.
func (entry *Entry) Debugf(format string, args ...interface{}) {
	if entry.Logger.level() >= DebugLevel {
		entry.Debug(fmt.Sprintf(format, args...))
	}
}

// Infof logs a info-level message using Printf.
func (entry *Entry) Infof(format string, args ...interface{}) {
	if entry.Logger.level() >= InfoLevel {
		entry.Info(fmt.Sprintf(format, args...))
	}
}

// Printf logs a info-level message using Printf.
func (entry *Entry) Printf(format string, args ...interface{}) {
	entry.Infof(format, args...)
}

// Warnf logs a warn-level message using Printf.
func (entry *Entry) Warnf(format string, args ...interface{}) {
	if entry.Logger.level() >= WarnLevel {
		entry.Warn(fmt.Sprintf(format, args...))
	}
}

// Warningf logs a warn-level message using Printf.
func (entry *Entry) Warningf(format string, args ...interface{}) {
	entry.Warnf(format, args...)
}

// Errorf logs a error-level message using Printf.
func (entry *Entry) Errorf(format string, args ...interface{}) {
	if entry.Logger.level() >= ErrorLevel {
		entry.Error(fmt.Sprintf(format, args...))
	}
}

// Fatalf logs a fatal-level message using Printf.
func (entry *Entry) Fatalf(format string, args ...interface{}) {
	if entry.Logger.level() >= FatalLevel {
		entry.Fatal(fmt.Sprintf(format, args...))
	}
	Exit(1)
}

// Panicf logs a panic-level message using Printf.
func (entry *Entry) Panicf(format string, args ...interface{}) {
	if entry.Logger.level() >= PanicLevel {
		entry.Panic(fmt.Sprintf(format, args...))
	}
}

// Debugln logs a debug-level message using Println.
func (entry *Entry) Debugln(args ...interface{}) {
	if entry.Logger.level() >= DebugLevel {
		entry.Debug(entry.sprintlnn(args...))
	}
}

// Infoln logs a info-level message using Println.
func (entry *Entry) Infoln(args ...interface{}) {
	if entry.Logger.level() >= InfoLevel {
		entry.Info(entry.sprintlnn(args...))
	}
}

// Println logs a info-level message using Println.
func (entry *Entry) Println(args ...interface{}) {
	entry.Infoln(args...)
}

// Warnln logs a warn-level message using Println.
func (entry *Entry) Warnln(args ...interface{}) {
	if entry.Logger.level() >= WarnLevel {
		entry.Warn(entry.sprintlnn(args...))
	}
}

// Warningln logs a warn-level message using Println.
func (entry *Entry) Warningln(args ...interface{}) {
	entry.Warnln(args...)
}

// Errorln logs a error-level message using Println.
func (entry *Entry) Errorln(args ...interface{}) {
	if entry.Logger.level() >= ErrorLevel {
		entry.Error(entry.sprintlnn(args...))
	}
}

// Fatalln logs a fatal-level message using Println.
func (entry *Entry) Fatalln(args ...interface{}) {
	if entry.Logger.level() >= FatalLevel {
		entry.Fatal(entry.sprintlnn(args...))
	}
	Exit(1)
}

// Panicln logs a panic-level message using Println.
func (entry *Entry) Panicln(args ...interface{}) {
	if entry.Logger.level() >= PanicLevel {
		entry.Panic(entry.sprintlnn(args...))
	}
}

// Sprintlnn => Sprint no newline. This is to get the behavior of how
// fmt.Sprintln where spaces are always added between operands, regardless of
// their type. Instead of vendoring the Sprintln implementation to spare a
// string allocation, we do the simplest thing.
func (entry *Entry) sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}
