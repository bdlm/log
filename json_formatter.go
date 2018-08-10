package log

import (
	"encoding/json"
	"fmt"
)

// JSONFormatter formats logs into parsable json
type JSONFormatter struct {
	// DataKey allows users to put all the log entry parameters into a
	// nested dictionary at a given key.
	DataKey string

	// DisableCaller controls caller logging.
	DisableCaller bool

	// DisableHostname controls hostname logging.
	DisableHostname bool

	// DisableLevel controls level logging.
	DisableLevel bool

	// DisableMessage controls message logging.
	DisableMessage bool

	// DisableTimestamp controls timestamp logging.
	DisableTimestamp bool

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	//  formatter := &JSONFormatter{FieldMap: FieldMap{
	//      LabelTime:  "@timestamp",
	//      LabelLevel: "@level",
	//      LabelMsg:   "@message",
	//  }}
	FieldMap FieldMap

	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string
}

// Format renders a single log entry
func (f *JSONFormatter) Format(entry *Entry) ([]byte, error) {
	prefixFieldClashes(entry.Data, f.FieldMap)

	data := getData(entry, f.FieldMap)
	jsonData := map[string]interface{}{}

	//
	if f.DisableTimestamp {
		delete(jsonData, f.FieldMap.resolve(LabelTime))
	} else if "" != f.TimestampFormat {
		jsonData[f.FieldMap.resolve(LabelTime)] = entry.Time.Format(f.TimestampFormat)
	} else {
		jsonData[f.FieldMap.resolve(LabelTime)] = entry.Time.Format(defaultTimestampFormat)
	}

	//
	if f.DisableHostname {
		delete(jsonData, f.FieldMap.resolve(LabelHost))
	} else {
		jsonData[f.FieldMap.resolve(LabelHost)] = data.Hostname
	}

	//
	if f.DisableCaller {
		delete(jsonData, f.FieldMap.resolve(LabelCaller))
	} else {
		jsonData[f.FieldMap.resolve(LabelCaller)] = data.Caller
	}

	//
	if f.DisableLevel {
		delete(jsonData, f.FieldMap.resolve(LabelLevel))
	} else {
		jsonData[f.FieldMap.resolve(LabelLevel)] = data.Level
	}

	//
	if f.DisableMessage {
		delete(jsonData, f.FieldMap.resolve(LabelMsg))
	} else {
		jsonData[f.FieldMap.resolve(LabelMsg)] = data.Message
	}

	//
	jsonData[f.FieldMap.resolve(LabelData)] = data.Data

	serialized, err := json.Marshal(jsonData)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
