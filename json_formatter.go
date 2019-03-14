package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"text/template"
)

var funcMap = template.FuncMap{
	// The name "title" is what the function will be called in the template text.
	"newCounter": func(data interface{}) map[string]int {
		var max int
		switch v := data.(type) {
		case map[string]interface{}:
			max = len(v) - 1
		case []string:
			max = len(v) - 1
		}
		return map[string]int{
			"max":   max,
			"cnt":   0,
			"comma": 1,
		}
	},
	"inc": func(counter map[string]int) map[string]int {
		counter["cnt"] = counter["cnt"] + 1
		if counter["cnt"] >= counter["max"] {
			counter["comma"] = 0
		}
		return counter
	},
	"json": func(v interface{}, prefix, indent string) string {
		byts, err := json.MarshalIndent(v, prefix, indent)
		if nil == err {
			return string(byts)
		}
		return ""
	},
}
var jsonTermTemplate = template.Must(template.New("tty").Funcs(funcMap).Parse(
	"{{$color := .Color}}{{$caller := .Caller}}{\n" +
		// Level
		"    \"{{$color.Level}}{{.LabelLevel}}{{$color.Reset}}\": \"{{$color.Level}}{{printf \"%5s\" .Level}}{{$color.Reset}}\",\n" +
		// Hostname
		"{{if .Hostname}}" +
		"    \"{{$color.Level}}{{.LabelHost}}{{$color.Reset}}\": \"{{$color.Hostname}}{{.Hostname}}{{$color.Reset}}\",\n" +
		"{{end}}" +
		// Timestamp
		"{{if .Timestamp}}" +
		"    \"{{$color.Level}}{{.LabelTime}}{{$color.Reset}}\": \"{{$color.Timestamp}}{{.Timestamp}}{{$color.Reset}}\",\n" +
		"{{end}}" +
		// Message
		"    \"{{$color.Level}}{{.LabelMsg}}{{$color.Reset}}\": \"{{printf \"%s\" .Message}}\",\n" +
		// Data fields
		"{{if .Data}}" +
		"{{$counter := newCounter .Data}}" +
		"    \"{{$color.Level}}{{.LabelData}}{{$color.Reset}}\": {\n{{range $k, $v := .Data}}" +
		"        \"{{$color.DataLabel}}{{$k}}{{$color.Reset}}\": {{$color.DataValue}}{{json $v (printf \"%s        \" $color.DataValue) \"    \"}}{{$color.Reset}}{{if eq 1 $counter.comma}},{{end}}\n" +
		"{{$_ := inc $counter}}" +
		"{{end}}    },\n" +
		"{{end}}" +
		// Caller
		"{{if and (.Caller) (not .Trace)}}" +
		"    \"{{$color.Level}}{{.LabelCaller}}{{$color.Reset}}\": \"{{$color.Caller}}{{.Caller}}{{$color.Reset}}\"\n" +
		"{{end}}" +
		// Trace
		"{{if .Trace}}" +
		"{{$counter := newCounter .Trace}}" +
		"    \"{{$color.Level}}{{.LabelTrace}}{{$color.Reset}}\": [\n{{range $k, $v := .Trace}}" +
		"        \"{{$color.Trace}}{{$v}}{{$color.Reset}}\"{{if eq 1 $counter.comma}},{{end}}\n" +
		"{{$_ := inc $counter}}" +
		"{{end}}    ]\n" +
		"{{end}}" +
		"}",
))

// JSONFormatter formats logs into parsable json.
type JSONFormatter struct {
	// DataKey allows users to put all the log entry parameters into a
	// nested dictionary at a given key.
	DataKey string

	// DisableCaller disables caller data output.
	DisableCaller bool

	// DisableHostname disables hostname output.
	DisableHostname bool

	// DisableLevel disables level output.
	DisableLevel bool

	// DisableMessage disables message output.
	DisableMessage bool

	// DisableTimestamp disables timestamp output.
	DisableTimestamp bool

	// DisableTTY disables TTY formatted output.
	DisableTTY bool

	// Enable full backtrace output.
	EnableTrace bool

	// EscapeHTML is a flag that notes whether HTML characters should be
	// escaped.
	EscapeHTML bool

	// ForceTTY forces TTY formatted output.
	ForceTTY bool

	// FieldMap allows users to customize the names of keys for default
	// fields.
	//
	// For example:
	// 	formatter := &TextFormatter{FieldMap: FieldMap{
	//      LabelCaller: "@caller",
	//      LabelData:   "@data",
	//      LabelHost:   "@hostname",
	//      LabelLevel:  "@loglevel",
	//      LabelMsg:    "@message",
	//      LabelTime:   "@timestamp",
	// 	}}
	FieldMap FieldMap

	// TimestampFormat allows a custom timestamp format to be used.
	TimestampFormat string

	// Flag noting whether the logger's output is to a terminal
	isTerminal bool

	sync.Once
}

func (f *JSONFormatter) init(entry *Entry) {
	if entry.Logger != nil {
		f.isTerminal = checkIfTerminal(entry.Logger.Out)
	}
}

// Format renders a single log entry
func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	var err error
	var serialized []byte
	isTTY := (f.ForceTTY || f.isTerminal) && !f.DisableTTY

	prefixFieldClashes(entry.Data, f.FieldMap)
	f.Do(func() { f.init(entry) })

	data := getData(entry, f.FieldMap, f.EscapeHTML, isTTY)

	if f.DisableTimestamp {
		data.Timestamp = ""
	} else if "" != f.TimestampFormat {
		data.Timestamp = entry.Time.Format(f.TimestampFormat)
	} else {
		data.Timestamp = entry.Time.Format(defaultTimestampFormat)
	}
	if f.DisableHostname {
		data.Hostname = ""
	}
	if f.DisableCaller {
		data.Caller = ""
	}
	if !f.EnableTrace {
		data.Trace = []string{}
	}

	if isTTY {
		var logLine *bytes.Buffer
		if entry.Buffer != nil {
			logLine = entry.Buffer
		} else {
			logLine = &bytes.Buffer{}
		}
		err = jsonTermTemplate.Execute(logLine, data)
		if err != nil {
			return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
		}
		serialized = logLine.Bytes()

	} else {
		// Relabel data for JSON output
		jsonData := map[string]interface{}{}
		if !f.DisableCaller || f.EnableTrace {
			jsonData[f.FieldMap.resolve(LabelCaller)] = data.Caller
		}
		if f.EnableTrace {
			jsonData[f.FieldMap.resolve(LabelTrace)] = data.Trace
		}
		if !f.DisableHostname {
			jsonData[f.FieldMap.resolve(LabelHost)] = data.Hostname
		}
		if !f.DisableLevel {
			jsonData[f.FieldMap.resolve(LabelLevel)] = data.Level
		}
		if !f.DisableMessage {
			jsonData[f.FieldMap.resolve(LabelMsg)] = data.Message
		}
		if !f.DisableTimestamp {
			if "" != f.TimestampFormat {
				jsonData[f.FieldMap.resolve(LabelTime)] = entry.Time.Format(f.TimestampFormat)
			} else {
				jsonData[f.FieldMap.resolve(LabelTime)] = entry.Time.Format(defaultTimestampFormat)
			}
		}

		//
		jsonData[f.FieldMap.resolve(LabelData)] = data.Data

		buf := new(bytes.Buffer)
		encoder := json.NewEncoder(buf)
		encoder.SetEscapeHTML(f.EscapeHTML)

		err = encoder.Encode(jsonData)
		serialized = []byte(strings.Trim(buf.String(), "\n"))
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
