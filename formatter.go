package log

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

// RFC3339Milli defines an RFC3339 date format with miliseconds
const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

const defaultTimestampFormat = RFC3339Milli

type FieldLabel string

// FieldMap allows customization of the key names for default fields.
type FieldMap map[FieldLabel]string

// Default key names for the default fields
const (
	LabelCaller = "caller"
	LabelData   = "data"
	LabelHost   = "host"
	LabelLevel  = "level"
	LabelMsg    = "msg"
	LabelTime   = "time"
)

func (f FieldMap) resolve(fieldLabel FieldLabel) string {
	if definedLabel, ok := f[fieldLabel]; ok {
		return definedLabel
	}
	return string(fieldLabel)
}

type logData struct {
	Timestamp string                 `json:"time,omitempty"`
	Level     string                 `json:"level,omitempty"`
	Hostname  string                 `json:"host,omitempty"`
	Message   string                 `json:"msg,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Caller    string                 `json:"caller,omitempty"`
	Color     string                 `json:"-"`
}

func getCaller() string {
	caller := ""
	a := 0
	for {
		if pc, file, line, ok := runtime.Caller(a); ok {
			if !strings.Contains(strings.ToLower(file), "github.com/bdlm/log") ||
				strings.HasSuffix(strings.ToLower(file), "_test.go") {
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

const (
	DEFAULTColor = "\033[38;5;46m"
	ERRColor     = "\033[38;5;196m"
	WARNColor    = "\033[38;5;226m"
	DEBUGColor   = "\033[38;5;245m"

	BOLDWHITEColor = "\033[0m\033[1m"
	BLUEColor      = "\033[38;5;4m" // -- blue
	GREYColor      = "\033[38;5;8m" // -- grey
	LTBLUEColor    = "\033[38;5;6m" // -- lt blue
	PURPLEColor    = "\033[38;5;5m" // -- purple
	WHITEColor     = "\033[38;5;7m" // -- white

	TSTColor = "\033[38;5;94m" // -- tst
)

/*
getData is a helper function that extracts log data from the Entry.
*/
func getData(entry *Entry, fieldMap FieldMap) *logData {
	var levelColor string
	switch entry.Level {
	case DebugLevel:
		levelColor = DEBUGColor
	case WarnLevel:
		levelColor = WARNColor
	case ErrorLevel, FatalLevel, PanicLevel:
		levelColor = ERRColor
	default:
		levelColor = DEFAULTColor
	}

	data := &logData{
		Caller:    strings.Trim(strconv.QuoteToASCII(getCaller()), `"`),
		Data:      make(map[string]interface{}),
		Hostname:  strings.Trim(strconv.QuoteToASCII(os.Getenv("HOSTNAME")), `"`),
		Level:     strings.Trim(strconv.QuoteToASCII(entry.Level.String()), `"`),
		Message:   entry.Message,
		Timestamp: entry.Time.Format(RFC3339Milli),
		Color:     levelColor,
	}

	keys := make([]string, 0)
	for k, v := range entry.Data {
		switch k {
		case fieldMap.resolve(LabelCaller):
			data.Caller = v.(string)
		case fieldMap.resolve(LabelHost):
			data.Hostname = v.(string)
		case fieldMap.resolve(LabelLevel):
			data.Level = v.(string)
		case fieldMap.resolve(LabelMsg):
			data.Message = v.(string)
		case fieldMap.resolve(LabelTime):
			data.Timestamp = v.(string)

		case fieldMap.resolve(LabelData):
			fallthrough
		default:
			keys = append(keys, k)
			switch v.(type) {
			case string:
				data.Data[strings.TrimPrefix(k, fieldMap.resolve(LabelData)+".")] = strings.Trim(strconv.QuoteToASCII(fmt.Sprintf("%v", v)), `"`)
			case error:
				data.Data[strings.TrimPrefix(k, fieldMap.resolve(LabelData)+".")] = v.(error).Error()
			default:
				data.Data[strings.TrimPrefix(k, fieldMap.resolve(LabelData)+".")] = v
			}
		}
	}
	sort.Strings(keys)

	return data
}

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
	var key string
	for _, field := range []FieldLabel{
		LabelCaller,
		LabelData,
		LabelHost,
		LabelLevel,
		LabelMsg,
		LabelTime,
	} {
		key = fieldMap.resolve(field)
		if t, ok := data[key]; ok {
			data[fieldMap.resolve(LabelData)+"."+key] = t
			delete(data, key)
		}
	}
}
