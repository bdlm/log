package log

import (
	"bytes"
	"strings"
	"sync"
	"text/template"
)

var (
	termTemplate = template.Must(template.New("tty").Parse(
		"{{.Color}}{{printf \"%5s\" .Level}}\033[0m:{{if .Hostname}} \033[38;5;39m{{.Hostname}}\033[0m{{end}}{{if .Timestamp}} \033[38;5;3m{{.Timestamp}}\033[0m{{end}} \033[1;97m\"{{printf \"%.125s\" .Message}}\"\033[0m\n    {{.Color}}⇢\033[0m {{range $k, $v := .Data}} \033[38;5;159m{{$k}}\033[0m=\"\033[38;5;180m{{$v}}\033[0m\"{{end}}{{if .Caller}}{{$length := len .Data }}{{if ne $length  0}}\n    {{.Color}}⇢\033[0m {{end}} \033[38;5;28m{{.Caller}}\033[0m{{end}}\n",
	))
	textTemplate = template.Must(template.New("log").Parse(
		`{{if .Timestamp}} {{.LabelTime}}="{{.Timestamp}}"{{end}} {{.LabelLevel}}="{{.Level}}"{{if .Message}} {{.LabelMsg}}="{{.Message}}"{{end}}{{$labelData := .LabelData}}{{range $k, $v := .Data}} {{if $labelData}}{{$labelData}}.{{end}}{{$k}}="{{$v}}"{{end}}{{if .Caller}} {{.LabelCaller}}="{{.Caller}}"{{end}}{{if .Hostname}} {{.LabelHost}}="{{.Hostname}}"{{end}}`,
	))
)

// TextFormatter formats logs into text.
type TextFormatter struct {
	// Disable caller data.
	DisableCaller bool

	// Disable hostname logging.
	DisableHostname bool

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool

	// Force disabling colors.
	DisableTTY bool

	// FieldMap allows users to customize the names of keys for default fields.
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

	// Set to true to bypass checking for a TTY before outputting colors.
	ForceTTY bool

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// Whether the logger's out is to a terminal
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
	prefixFieldClashes(entry.Data, f.FieldMap)

	var logLine *bytes.Buffer
	if entry.Buffer != nil {
		logLine = entry.Buffer
	} else {
		logLine = &bytes.Buffer{}
	}

	f.Do(func() { f.init(entry) })

	data := getData(entry, f.FieldMap)
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

	isTTY := (f.ForceTTY || f.isTerminal) && !f.DisableTTY
	if isTTY {
		err = termTemplate.Execute(logLine, data)
	} else {
		err = textTemplate.Execute(logLine, data)
	}
	if nil != err {
		return nil, err
	}

	return append([]byte(strings.Trim(logLine.String(), " \n")), '\n'), nil
}
