package log

import (
	"encoding/json"
	"fmt"
)

type FieldLabel string

// FieldMap allows customization of the key names for default fields.
type FieldMap map[FieldLabel]string

var fieldMap = FieldMap{}

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

// JSONFormatter formats logs into parsable json
type JSONFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool

	// Disable hostname logging.
	DisableHostname bool

	// DataKey allows users to put all the log entry parameters into a
	// nested dictionary at a given key.
	DataKey string

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	//  formatter := &JSONFormatter{FieldMap: FieldMap{
	//      LabelTime:  "@timestamp",
	//      LabelLevel: "@level",
	//      LabelMsg:   "@message",
	//  }}
	FieldMap FieldMap
}

// Format renders a single log entry
func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	//data := make(Fields, len(entry.Data)+3)
	//for k, v := range entry.Data {
	//	switch v := v.(type) {
	//	case error:
	//		// Otherwise errors are ignored by `encoding/json`
	//		// https://github.com/sirupsen/logrus/issues/137
	//		data[k] = v.Error()
	//	default:
	//		data[k] = v
	//	}
	//}
	//
	//if f.DataKey != "" {
	//	newData := make(Fields, 4)
	//	newData[f.DataKey] = data
	//	data = newData
	//}

	prefixFieldClashes(entry.Data, f.FieldMap)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	data := getData(entry)
	jsonData := map[string]interface{}{}

	if !f.DisableTimestamp && "" != f.TimestampFormat {
		jsonData[fieldMap.resolve(LabelTime)] = entry.Time.Format(f.TimestampFormat)
	}
	if !f.DisableHostname {
		jsonData[fieldMap.resolve(LabelTime)] = data.Hostname
	}
	jsonData[fieldMap.resolve(LabelCaller)] = data.Caller
	jsonData[fieldMap.resolve(LabelData)] = data.Data
	jsonData[fieldMap.resolve(LabelLevel)] = data.Level
	jsonData[fieldMap.resolve(LabelMsg)] = data.Message

	//	if f.DisableTimestamp {
	//		data[f.FieldMap.resolve(LabelTime)] = ""
	//	} else {
	//		data[f.FieldMap.resolve(LabelTime)] = entry.Time.Format(timestampFormat)
	//	}
	//
	//	data[f.FieldMap.resolve(LabelMsg)] = entry.Message
	//	data[f.FieldMap.resolve(LabelLevel)] = entry.Level.String()

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
