package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// RFC3339Milli defines an RFC3339 date format with miliseconds
const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

const defaultTimestampFormat = RFC3339Milli

// FieldLabel is a type for defining label keys.
type FieldLabel string

// FieldMap allows customization of the key names for default fields.
type FieldMap map[FieldLabel]string

// Default key names for the default fields
const (
	LabelCaller = "caller"
	LabelData   = "data"
	LabelError  = "error"
	LabelHost   = "host"
	LabelLevel  = "level"
	LabelMsg    = "msg"
	LabelTime   = "time"
	LabelTrace  = "trace"
)

func (f FieldMap) resolve(fieldLabel FieldLabel) string {
	if definedLabel, ok := f[fieldLabel]; ok {
		return definedLabel
	}
	return string(fieldLabel)
}

type logData struct {
	LabelCaller string   `json:"-"`
	LabelData   string   `json:"-"`
	LabelError  string   `json:"-"`
	LabelHost   string   `json:"-"`
	LabelLevel  string   `json:"-"`
	LabelMsg    string   `json:"-"`
	LabelTime   string   `json:"-"`
	LabelTrace  string   `json:"-"`
	Color       colors   `json:"-"`
	ErrData     []string `json:"-"`

	Caller    string                 `json:"caller,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Err       interface{}            `json:"error,omitempty"`
	Hostname  string                 `json:"host,omitempty"`
	Level     string                 `json:"level,omitempty"`
	Message   string                 `json:"msg,omitempty"`
	Timestamp string                 `json:"time,omitempty"`
	Trace     []string               `json:"trace,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (l *logData) UnmarshalJSON(d []byte) error {
	data := map[string]interface{}{}

	err := json.Unmarshal(d, &data)
	if nil != err {
		return err
	}

	if _, ok := data["caller"]; ok {
		l.Caller = data["caller"].(string)
	}
	if _, ok := data["data"]; ok {
		if nil == l.Data {
			l.Data = map[string]interface{}{}
		}
		for k, v := range data["data"].(map[string]interface{}) {
			if e, ok := v.(error); ok {
				l.Data[k] = e.(error)
			} else {
				l.Data[k] = v
			}
		}
		//l.Data = data["data"].(map[string]interface{})
	}
	if _, ok := data["error"]; ok && "" != data["error"] {
		l.Err = fmt.Errorf(data["error"].(string))
	}
	if _, ok := data["host"]; ok {
		l.Hostname = data["host"].(string)
	}
	if _, ok := data["level"]; ok {
		l.Level = data["level"].(string)
	}
	if _, ok := data["msg"]; ok {
		l.Message = data["msg"].(string)
	}
	if _, ok := data["time"]; ok {
		l.Timestamp = data["time"].(string)
	}
	if _, ok := data["trace"]; ok {
		l.Trace = data["trace"].([]string)
	}

	return nil
}

// SetCallerLevel will adjust the relative caller level in log output.
func SetCallerLevel(level int) {
	callerLevel = level
}

var callerLevel int

func getCaller() string {
	caller := ""
	a := 0
	for {
		if pc, file, line, ok := runtime.Caller(a); ok {
			if !strings.Contains(strings.ToLower(file), "github.com/bdlm/log") ||
				strings.HasSuffix(strings.ToLower(file), "_test.go") {
				if 0 == callerLevel {
					caller = fmt.Sprintf("%s:%d %s", path.Base(file), line, runtime.FuncForPC(pc).Name())
				} else {
					if pc2, file2, line2, ok := runtime.Caller(a + callerLevel); ok {
						caller = fmt.Sprintf("%s:%d %s", path.Base(file2), line2, runtime.FuncForPC(pc2).Name())
					} else {
						caller = fmt.Sprintf("%s:%d %s", path.Base(file), line, runtime.FuncForPC(pc).Name())
					}
				}
				break
			}
		} else {
			break
		}
		a++
	}
	return caller
}

func getTrace() []string {
	trace := []string{}
	a := 0
	for {
		if pc, file, line, ok := runtime.Caller(a); ok {
			if !strings.Contains(strings.ToLower(file), "github.com/bdlm/log") ||
				strings.HasSuffix(strings.ToLower(file), "_test.go") {
				if 0 == callerLevel {
					trace = append(trace, fmt.Sprintf("%s:%d %s", path.Base(file), line, runtime.FuncForPC(pc).Name()))
				} else {
					if pc2, file2, line2, ok := runtime.Caller(a + callerLevel); ok {
						trace = append(trace, fmt.Sprintf("%s:%d %s", path.Base(file2), line2, runtime.FuncForPC(pc2).Name()))
					} else {
						trace = append(trace, fmt.Sprintf("%s:%d %s", path.Base(file), line, runtime.FuncForPC(pc).Name()))
					}
				}
			}
		} else {
			break
		}
		a++
	}
	if len(trace) > 2 {
		trace = trace[:len(trace)-2]
	}
	return trace
}

var (
	// DEFAULTColor is the default TTY 'level' color.
	DEFAULTColor = "\033[38;5;46m"
	// ERRORColor is the TTY 'level' color for error messages.
	ERRORColor = "\033[38;5;166m"
	// FATALColor is the TTY 'level' color for fatal messages.
	FATALColor = "\033[38;5;124m"
	// PANICColor is the TTY 'level' color for panic messages.
	PANICColor = "\033[38;5;196m"
	// WARNColor is the TTY 'level' color for warning messages.
	WARNColor = "\033[38;5;226m"
	// DEBUGColor is the TTY 'level' color for debug messages.
	DEBUGColor = "\033[38;5;245m"

	// CallerColor is the TTY caller color.
	CallerColor = "\033[38;5;244m"
	// DataLabelColor is the TTY data label color.
	DataLabelColor = "\033[38;5;111m"
	// DataValueColor is the TTY data value color.
	DataValueColor = "\033[38;5;180m"
	// HostnameColor is the TTY hostname color.
	HostnameColor = "\033[38;5;39m"
	// TraceColor is the TTY trace color.
	TraceColor = "\033[38;5;244m"
	// TimestampColor is the TTY timestamp color.
	TimestampColor = "\033[38;5;72m"

	// ResetColor resets the TTY color scheme to it's default.
	ResetColor = "\033[0m"
)

type colors struct {
	Caller    string
	DataLabel string
	DataValue string
	Err       string
	Hostname  string
	Level     string
	Reset     string
	Timestamp string
	Trace     string
}

func escape(data interface{}, escapeHTML bool) string {
	var result string
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(escapeHTML)
	err := encoder.Encode(data)
	if nil == err {
		result = strings.Trim(buf.String(), "\n")
	}
	return result
}

// getData is a helper function that extracts log data from the Entry.
func getData(entry *Entry, fieldMap FieldMap, escapeHTML, isTTY bool) *logData {
	var levelColor string

	data := &logData{
		Caller:    getCaller(),
		Data:      map[string]interface{}{},
		Err:       entry.Err,
		ErrData:   []string{},
		Hostname:  os.Getenv("HOSTNAME"),
		Level:     LevelString(entry.Level),
		Message:   entry.Message,
		Timestamp: entry.Time.Format(RFC3339Milli),
		Trace:     getTrace(),
	}

	data.LabelCaller = fieldMap.resolve(LabelCaller)
	data.LabelData = fieldMap.resolve(LabelData)
	data.LabelError = fieldMap.resolve(LabelError)
	data.LabelHost = fieldMap.resolve(LabelHost)
	data.LabelLevel = fieldMap.resolve(LabelLevel)
	data.LabelMsg = fieldMap.resolve(LabelMsg)
	data.LabelTime = fieldMap.resolve(LabelTime)
	data.LabelTrace = fieldMap.resolve(LabelTrace)

	if isTTY {
		switch entry.Level {
		case DebugLevel:
			levelColor = DEBUGColor
		case WarnLevel:
			levelColor = WARNColor
		case ErrorLevel:
			levelColor = ERRORColor
		case FatalLevel:
			levelColor = FATALColor
		case PanicLevel:
			levelColor = PANICColor
		default:
			levelColor = DEFAULTColor
		}
		data.Color = colors{
			Caller:    CallerColor,
			DataLabel: DataLabelColor,
			DataValue: DataValueColor,
			Err:       ERRORColor,
			Hostname:  HostnameColor,
			Level:     levelColor,
			Reset:     ResetColor,
			Timestamp: TimestampColor,
			Trace:     TraceColor,
		}
	}

	remapData(entry, fieldMap, data)

	return data
}

func remapData(entry *Entry, fieldMap FieldMap, data *logData) {
	for k, v := range entry.Data {
		switch k {
		default:
			switch v := v.(type) {
			case string:
				data.Data[strings.TrimPrefix(k, fieldMap.resolve(LabelData)+".")] = strings.Trim(strconv.QuoteToASCII(fmt.Sprintf("%v", v)), `"`)
			default:
				data.Data[strings.TrimPrefix(k, fieldMap.resolve(LabelData)+".")] = v
			}
		}
	}
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
		LabelError,
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
