package log

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"sync"
	"text/template"
)

var (
	termTemplate = template.Must(template.New("tty").Parse(
		"[{{.Color}}{{printf \"%.4s\" .Level}}\033[0m]{{if .Hostname}} \033[38;5;39m{{.Hostname}}\033[0m{{end}}{{if .Timestamp}} \033[38;5;3m{{.Timestamp}}\033[0m{{end}} \033[1;97m\"{{printf \"%.125s\" .Message}}\"\033[0m\n    ⇢ {{range $k, $v := .Data}} \033[38;5;159m{{$k}}\033[0m=\"\033[38;5;180m{{$v}}\033[0m\"{{end}}{{if .Caller}}{{$length := len .Data }}{{if ne $length  0}}\n    ⇢ {{end}} \033[38;5;124m{{.Caller}}\033[0m{{end}}\n",
	))
	textTemplate = template.Must(template.New("log").Parse(
		`{{if .Timestamp}} {{.LabelTime}}="{{.Timestamp}}"{{end}} {{.LabelLevel}}="{{.Level}}"{{if .Message}} {{.LabelMsg}}="{{.Message}}"{{end}}{{$labelData := .LabelData}}{{range $k, $v := .Data}} {{if $labelData}}{{$labelData}}.{{end}}{{$k}}="{{$v}}"{{end}}{{if .Caller}} {{.LabelCaller}}="{{.Caller}}"{{end}}{{if .Hostname}} {{.LabelHost}}="{{.Hostname}}"{{end}}`,
	))
)

// TextFormatter formats logs into text.
type TextFormatter struct {
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Disable caller data.
	DisableCaller bool

	// Force disabling colors.
	DisableColors bool

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool

	// Disable hostname logging.
	DisableHostname bool

	// Enable logging the full timestamp when a TTY is attached instead of just
	// the time passed since beginning of execution.
	FullTimestamp bool

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// Disables the truncation of the level text to 4 characters.
	DisableLevelTruncation bool

	// QuoteEmptyFields will wrap empty fields in quotes if true
	QuoteEmptyFields bool

	// Whether the logger's out is to a terminal
	isTerminal bool

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	// 	formatter := &TextFormatter{FieldMap: FieldMap{
	// 		LabelTime:  "@timestamp",
	// 		LabelLevel: "@level",
	// 		LabelMsg:   "@message",
	// 	}}
	FieldMap FieldMap

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

	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	if !f.DisableSorting {
		sort.Strings(keys)
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

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
	}
	if f.DisableHostname {
		data.Hostname = ""
	}
	if f.DisableCaller {
		data.Caller = ""
	}

	isColorTerm := (f.ForceColors || f.isTerminal) && !f.DisableColors
	if isColorTerm {
		err = termTemplate.Execute(logLine, data)
	} else {
		err = textTemplate.Execute(logLine, data)
	}
	if nil != err {
		return nil, err
	}

	return append([]byte(strings.Trim(logLine.String(), " \n")), '\n'), nil
}

func (f *TextFormatter) needsQuoting(text string) bool {
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *TextFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}
