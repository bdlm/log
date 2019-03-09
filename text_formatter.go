package log

import (
	"bytes"
	"strings"
	"sync"
	"text/template"
)

var (
	termTemplate = template.Must(template.New("tty").Parse(
		"{{$color := .Color}}{{$caller := .Caller}}" +
			// Level
			"{{$color.Level}}{{printf \"%5s\" .Level}}{{$color.Reset}}" +
			// Hostname
			"{{if .Hostname}} {{$color.Hostname}}{{.Hostname}}{{$color.Reset}}{{end}} " +
			// Message
			"{{printf \"%s\" .Message}}" +
			// Data fields
			"{{if .Data}}\n   {{$color.Level}}⇢{{$color.Reset}} {{range $k, $v := .Data}}" +
			" {{$color.DataLabel}}{{$k}}{{$color.Reset}}={{$color.DataValue}}{{$v}}{{$color.Reset}}" +
			"{{end}}{{end}}" +
			// Caller
			"{{if and (.Caller) (not .Trace)}}" +
			"\n   {{$color.Level}}⇢{{$color.Reset}}  {{$color.Caller}}{{.Caller}}{{$color.Reset}}" +
			"{{end}}" +
			// Trace
			"{{range $k, $v := .Trace}}" +
			"\n   {{$color.Level}}⇢{{$color.Reset}}  {{$color.Trace}}#{{$k}} {{$v}}{{$color.Reset}}" +
			"{{end}}" +
			// Timestamp
			"{{if .Timestamp}}\n   {{$color.Level}}⇢{{$color.Reset}}  {{$color.Timestamp}}{{.Timestamp}}{{$color.Reset}}{{end}}\n",
	))
	//\n   {{$color}}⇢\033[0m  {{if eq $v $caller}}\033[38;5;28m{{else}}\033[38;5;240m{{end}}#{{$k}} {{$v}}\033[0m{{end}}
	textTemplate = template.Must(template.New("log").Parse(
		// Timestamp
		"{{if .Timestamp}} {{.LabelTime}}=\"{{.Timestamp}}\"{{end}} " +
			// Level
			"{{.LabelLevel}}=\"{{.Level}}\"" +
			// Message
			"{{if .Message}} {{.LabelMsg}}={{.Message}}{{end}}" +
			// Data fields
			"{{$labelData := .LabelData}}{{range $k, $v := .Data}} {{if $labelData}}{{$labelData}}.{{end}}{{$k}}={{$v}}{{end}}" +
			// Caller
			"{{if .Caller}} {{.LabelCaller}}=\"{{.Caller}}\"{{end}}" +
			// Hostname
			"{{if .Hostname}} {{.LabelHost}}=\"{{.Hostname}}\"{{end}}" +
			// Trace
			"{{range $k, $v := .Trace}} trace.{{$k}}=\"{{$v}}\"{{end}}",
	))
)

// TextFormatter formats logs into text.
type TextFormatter struct {
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

	// Flag noting whether the logger's out is to a terminal
	isTerminal bool

	sync.Once
}

func (f *TextFormatter) init(entry *Entry) {
	if entry.Logger != nil {
		f.isTerminal = checkIfTerminal(entry.Logger.Out)
	}
}

// Format renders a single log entry
func (f *TextFormatter) Format(entry *Entry) ([]byte, error) {
	var err error
	isTTY := (f.ForceTTY || f.isTerminal) && !f.DisableTTY
	prefixFieldClashes(entry.Data, f.FieldMap)

	var logLine *bytes.Buffer
	if entry.Buffer != nil {
		logLine = entry.Buffer
	} else {
		logLine = &bytes.Buffer{}
	}

	f.Do(func() { f.init(entry) })

	data := getData(entry, f.FieldMap, f.EscapeHTML, isTTY)
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

	for k, v := range data.Data {
		data.Data[k] = escape(v, f.EscapeHTML)
	}

	if isTTY {
		err = termTemplate.Execute(logLine, data)
	} else {
		data.Message = escape(data.Message, f.EscapeHTML)
		err = textTemplate.Execute(logLine, data)
	}
	if nil != err {
		return nil, err
	}

	return append([]byte(strings.Trim(logLine.String(), " \n")), '\n'), nil
}
