package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestErrorNotLost(t *testing.T) {
	formatter := &JSONFormatter{}
	entry := WithField("error", errors.New("wild walrus"))

	b, err := formatter.Format(entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	result := logData{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if result.Data["error"].(string) != "wild walrus" {
		t.Fatal("Error field not set")
	}
}

func TestErrorNotLostOnFieldNotNamedError(t *testing.T) {
	formatter := &JSONFormatter{}
	entry := WithField("omg", errors.New("wild walrus"))

	b, err := formatter.Format(entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	result := logData{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if result.Data["omg"] != "wild walrus" {
		t.Fatal("Error field not set")
	}
}

func TestFieldClashWithTime(t *testing.T) {
	formatter := &JSONFormatter{}
	entry := WithField("time", "right now!")

	b, err := formatter.Format(entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	result := logData{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}
	fmt.Printf("\n\n### %v\n\n", result.Data)
	if result.Data["time"] != "right now!" {
		t.Fatal("time not set to original time field")
	}

	if result.Timestamp != "0001-01-01T00:00:00.000Z" {
		t.Fatalf("time field not set to current time, was: '%s', expected: '%s'", result.Data["time"], "0001-01-01T00:00:00.000Z")
	}
}

func TestFieldClashWithMsg(t *testing.T) {
	formatter := &JSONFormatter{}
	entry := WithField("msg", "something")

	b, err := formatter.Format(entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	result := logData{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}
	if result.Data["msg"] != "something" {
		t.Fatal("msg not set to original msg field")
	}
}

func TestFieldClashWithLevel(t *testing.T) {
	formatter := &JSONFormatter{}
	entry := WithField("level", "something")

	b, err := formatter.Format(entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	result := logData{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if result.Data["level"] != "something" {
		t.Fatal("level not set to original level field")
	}
}

func TestFieldClashWithRemappedFields(t *testing.T) {
	formatter := &JSONFormatter{FieldMap: FieldMap{
		LabelTime:   "@timestamp",
		LabelLevel:  "@level",
		LabelMsg:    "@message",
		LabelData:   "@data",
		LabelCaller: "@caller",
	}}
	entry := WithFields(Fields{
		"@timestamp": "@timestamp",
		"@level":     "@level",
		"@message":   "@message",
		"time":       "time",
		"level":      "level",
		"msg":        "msg",
	})

	b, err := formatter.Format(entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	result := logData{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	for _, field := range []string{"timestamp", "level", "msg"} {
		if result.Data[field] == field {
			t.Errorf("Expected field %v to be untouched; got %v", field, result.Data[field])
		}

		remappedKey := fmt.Sprintf(formatter.FieldMap.resolve(LabelData)+".%s", field)
		if remapped, ok := result.Data[remappedKey]; ok {
			t.Errorf("Expected %s to be empty; got %v", remappedKey, remapped)
		}
	}

	for _, field := range []string{"@timestamp", "@level", "@message"} {
		if result.Data[field] == field {
			t.Errorf(
				"Expected field %v to be mapped to an Entry value: %v\n%s\n\n",
				field,
				result.Data,
				string(b),
			)
		}
	}
}

func TestFieldsInNestedDictionary(t *testing.T) {
	formatter := &JSONFormatter{
		DataKey: "args",
	}

	entry := WithFields(Fields{
		"level": "level",
		"test":  "test",
	})
	entry.Level = InfoLevel

	b, err := formatter.Format(entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	result := logData{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	//	args := entry["args"].(map[string]interface{})
	//
	//	for _, field := range []string{"test", "level"} {
	//		if value, present := args[field]; !present || value != field {
	//			t.Errorf("Expected field %v to be present under 'args'; untouched", field)
	//		}
	//	}
	//
	//	for _, field := range []string{"test", formatter.FieldMap.resolve(LabelData)+".level"} {
	//		if _, present := result.Data[field]; present {
	//			t.Errorf("Expected field %v not to be present at top level", field)
	//		}
	//	}
	//
	//	// with nested object, "level" shouldn't clash
	//	if result.Data["level"] != "info" {
	//		t.Errorf("Expected 'level' field to contain 'info'")
	//	}
}

func TestJSONEntryEndsWithNewline(t *testing.T) {
	formatter := &JSONFormatter{}
	entry := WithField("level", "something")

	b, err := formatter.Format(entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	if b[len(b)-1] != '\n' {
		t.Fatal("Expected JSON log entry to end with a newline")
	}
}

func TestJSONMessageKey(t *testing.T) {
	formatter := &JSONFormatter{}
	entry := &Entry{Message: "oh hai"}

	b, err := formatter.Format(entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !(strings.Contains(s, "msg") && strings.Contains(s, "oh hai")) {
		t.Fatal("Expected JSON to format message key")
	}
}

func TestJSONLevelKey(t *testing.T) {
	formatter := &JSONFormatter{}
	entry := WithField("level", "something")

	b, err := formatter.Format(entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, "level") {
		t.Fatal("Expected JSON to format level key")
	}
}

func TestJSONTimeKey(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, "time") {
		t.Fatal("Expected JSON to format time key")
	}
}

func TestJSONDisableTimestamp(t *testing.T) {
	formatter := &JSONFormatter{
		DisableTimestamp: true,
	}

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if strings.Contains(s, LabelTime) {
		t.Errorf("Did not prevent timestamp field '%s': %s", LabelTime, s)
	}
}

func TestJSONEnableTimestamp(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, LabelTime) {
		t.Error("Timestamp not present", s)
	}
}
