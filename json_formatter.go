package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// JSONFormatter formats logs into parsable json
type JSONFormatter struct {
	// DataKey allows users to put all the log entry parameters into a
	// nested dictionary at a given key.
	DataKey string

	// DisableCaller disables caller data.
	DisableCaller bool

	// DisableHostname disables hostname logging.
	DisableHostname bool

	// DisableLevel controls level logging.
	DisableLevel bool

	// DisableMessage controls message logging.
	DisableMessage bool

	// DisableTimestamp controls timestamp logging.
	DisableTimestamp bool

	// DisableTTY disables TTY formatted output.
	DisableTTY bool

	// Enable the full backtrace.
	EnableTrace bool

	// EscapeHTML escapes HTML characters.
	EscapeHTML bool

	// ForceTTY forces TTY output.
	ForceTTY bool

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

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// Flag noting whether the logger's out is to a terminal
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
	prefixFieldClashes(entry.Data, f.FieldMap)
	f.Do(func() { f.init(entry) })

	data := getData(entry, f.FieldMap, f.EscapeHTML)
	jsonData := map[string]interface{}{}

	//
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
	isTTY := (f.ForceTTY || f.isTerminal) && !f.DisableTTY
	var serialized []byte
	var err error

	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(f.EscapeHTML)

	if isTTY {
		encoder.SetIndent("", "    ")
		err = encoder.Encode(jsonData)
		serialized = []byte(strings.Trim(buf.String(), "\n"))
		serialized = append([]byte(data.Color), serialized...)
		serialized = append(serialized, []byte("\033[0m")...)
	} else {
		err = encoder.Encode(jsonData)
		serialized = []byte(strings.Trim(buf.String(), "\n"))
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
