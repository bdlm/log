package log

import (
	"bytes"
	"strings"
	"text/template"
)

var (
	stdTemplate = template.Must(template.New("log").Parse(
		"{{if .Timestamp}}{{.Timestamp}}{{end}}{{if .Message}} {{.Message}}{{end}}{{if .Level}} {{.LabelLevel}}=\"{{.Level}}\"{{end}}{{$labelData := .LabelData}}{{range $k, $v := .Data}} {{if $labelData}}{{$labelData}}.{{end}}{{$k}}={{$v}}{{end}}{{if .Caller}} {{.LabelCaller}}=\"{{.Caller}}\"{{end}}{{if .Hostname}} {{.LabelHost}}=\"{{.Hostname}}\"{{end}}{{range $k, $v := .Trace}} trace.{{$k}}=\"{{$v}}\"{{end}}",
	))
)

// StdFormatter formats logs into text.
type StdFormatter struct {
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

	// Enable full backtrace output.
	EnableTrace bool

	// EscapeHTML is a flag that notes whether HTML characters should be
	// escaped.
	EscapeHTML bool

	// FieldMap allows users to customize the names of keys for default
	// fields.
	//
	// For example:
	// 	formatter := &StdFormatter{FieldMap: FieldMap{
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
}

// Format renders a single log entry
func (f *StdFormatter) Format(entry *Entry) ([]byte, error) {
	var err error
	prefixFieldClashes(entry.Data, f.FieldMap)

	var logLine *bytes.Buffer
	if entry.Buffer != nil {
		logLine = entry.Buffer
	} else {
		logLine = &bytes.Buffer{}
	}

	data := getData(entry, f.FieldMap, f.EscapeHTML, false)
	data.LabelCaller = f.FieldMap.resolve(LabelCaller)
	data.LabelHost = f.FieldMap.resolve(LabelHost)
	data.LabelLevel = f.FieldMap.resolve(LabelLevel)
	data.LabelMsg = f.FieldMap.resolve(LabelMsg)
	data.LabelTime = f.FieldMap.resolve(LabelTime)
	data.LabelData = f.FieldMap.resolve(LabelData)

	if f.DisableTimestamp {
		data.Timestamp = ""
	} else if "" != f.TimestampFormat {
		data.Timestamp = entry.Time.Format(f.TimestampFormat)
	} else {
		data.Timestamp = entry.Time.Format("2006/01/02 15:04:05")
	}
	if f.DisableHostname {
		data.Hostname = ""
	}
	if f.DisableCaller && !f.EnableTrace {
		data.Caller = ""
	}
	if !f.EnableTrace {
		data.Trace = []string{}
	}

	for k, v := range data.Data {
		data.Data[k] = escape(v, f.EscapeHTML)
	}

	err = stdTemplate.Execute(logLine, data)
	if nil != err {
		return nil, err
	}

	return append([]byte(strings.Trim(logLine.String(), " \n")), '\n'), nil
}
