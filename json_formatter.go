package log

import (
	"encoding/json"
	"fmt"
	"sync"
)

// JSONFormatter formats logs into parsable json
type JSONFormatter struct {
	// Force disabling colors.
	DisableTTY bool

	// Set to true to bypass checking for a TTY before outputting colors.
	ForceTTY bool

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

	// Whether the logger's out is to a terminal
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

	data := getData(entry, f.FieldMap)
	jsonData := map[string]interface{}{}

	//
	if !f.DisableCaller {
		jsonData[f.FieldMap.resolve(LabelCaller)] = data.Caller
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
	if isTTY {
		serialized, err = json.MarshalIndent(jsonData, "", "    ")
		serialized = append([]byte(data.Color), serialized...)
		serialized = append(serialized, []byte("\033[0m")...)
		//serialized = []byte(strings.Replace(string(serialized), `"level": "debug",`, `"`+data.Color+`level`+"\033[0m"+`": "`+data.Color+`debug`+"\033[0m"+`",`, -1))
		//serialized = []byte(strings.Replace(string(serialized), `"level": "info",`, `"`+data.Color+`level`+"\033[0m"+`": "`+data.Color+`info`+"\033[0m"+`",`, -1))
		//serialized = []byte(strings.Replace(string(serialized), `"level": "warn",`, `"`+data.Color+`level`+"\033[0m"+`": "`+data.Color+`warn`+"\033[0m"+`",`, -1))
		//serialized = []byte(strings.Replace(string(serialized), `"level": "error",`, `"`+data.Color+`level`+"\033[0m"+`": "`+data.Color+`error`+"\033[0m"+`",`, -1))
		//serialized = []byte(strings.Replace(string(serialized), `"level": "panic",`, `"`+data.Color+`level`+"\033[0m"+`": "`+data.Color+`panic`+"\033[0m"+`",`, -1))
		//serialized = []byte(strings.Replace(string(serialized), `"level": "fatal",`, `"`+data.Color+`level`+"\033[0m"+`": "`+data.Color+`fatal`+"\033[0m"+`",`, -1))
	} else {
		serialized, err = json.Marshal(jsonData)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
