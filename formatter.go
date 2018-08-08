package log

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"
	"time"
)

/*
RFC3339Milli defines an RFC3339 date format with miliseconds
*/
const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

type logData struct {
	Timestamp string            `json:"time"`
	Level     string            `json:"level"`
	Hostname  string            `json:"host"`
	Message   string            `json:"msg"`
	Data      map[string]string `json:"data"`
	Caller    string            `json:"caller"`
	Color     int               `json:"_"`
}

func getCaller() string {
	caller := ""
	a := 0
	for {
		if pc, file, line, ok := runtime.Caller(a + 1); ok {
			if !strings.Contains(strings.ToLower(file), "github.com/bdlm/log") {
				caller = fmt.Sprintf("%s:%d %s", path.Base(file), line, runtime.FuncForPC(pc).Name())
				break
			}
		} else {
			break
		}
		a++
	}
	return caller
}

/*
getData is a helper function that extracts log data from the logrus
entry.
*/
func getData(entry *Entry) *logData {
	var levelColor int
	switch entry.Level {
	case DebugLevel:
		levelColor = gray
	case WarnLevel:
		levelColor = yellow
	case ErrorLevel, FatalLevel, PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}

	data := &logData{
		Caller:    getCaller(),
		Data:      make(map[string]string),
		Hostname:  os.Getenv("HOSTNAME"),
		Level:     entry.Level.String(),
		Message:   entry.Message,
		Timestamp: entry.Time.Format(RFC3339Milli),
		Color:     levelColor,
	}

	keys := make([]string, 0)
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for k, v := range entry.Data {
		data.Data[k] = fmt.Sprintf("%v", v)
	}

	return data
}

const defaultTimestampFormat = time.RFC3339

// The Formatter interface is used to implement a custom Formatter. It takes an
// `Entry`. It exposes all the fields, including the default ones:
//
// * `entry.Data["msg"]`. The message passed from Info, Warn, Error ..
// * `entry.Data["time"]`. The timestamp.
// * `entry.Data["level"]. The level the entry was logged at.
//
// Any additional fields added with `WithField` or `WithFields` are also in
// `entry.Data`. Format is expected to return an array of bytes which are then
// logged to `logger.Out`.
type Formatter interface {
	Format(*Entry) ([]byte, error)
}

// This is to not silently overwrite `time`, `msg` and `level` fields when
// dumping it. If this code wasn't there doing:
//
//  WithField("level", 1).Info("hello")
//
// Would just silently drop the user provided level. Instead with this code
// it'll logged as:
//
//  {"level": "info", "fields.level": 1, "msg": "hello", "time": "..."}
//
// It's not exported because it's still using Data in an opinionated way. It's to
// avoid code duplication between the two default formatters.
func prefixFieldClashes(data Fields, fieldMap FieldMap) {
	timeKey := fieldMap.resolve(FieldKeyTime)
	if t, ok := data[timeKey]; ok {
		data["fields."+timeKey] = t
		delete(data, timeKey)
	}

	msgKey := fieldMap.resolve(FieldKeyMsg)
	if m, ok := data[msgKey]; ok {
		data["fields."+msgKey] = m
		delete(data, msgKey)
	}

	levelKey := fieldMap.resolve(FieldKeyLevel)
	if l, ok := data[levelKey]; ok {
		data["fields."+levelKey] = l
		delete(data, levelKey)
	}
}
