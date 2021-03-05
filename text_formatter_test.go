package log

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatting(t *testing.T) {
	defer newStd()
	tf := &TextFormatter{
		DisableHostname: true,
		DisableTTY:      true,
	}

	testCases := []struct {
		value    string
		expected string
	}{
		{`foo`, "time=\"0001-01-01T00:00:00.000Z\" level=\"fatal\" msg=\"\" data.test=\"foo\" caller=\"text_formatter_test.go:28 github.com/bdlm/log/v2.TestFormatting\"\n"},
	}

	for _, tc := range testCases {
		b, _ := tf.Format(WithField("test", tc.value))

		if string(b) != tc.expected {
			t.Errorf(
				"formatting expected for %q (result was %q instead of %q)",
				tc.value,
				string(b),
				tc.expected,
			)
		}
	}
}

func TestEscaping(t *testing.T) {
	defer newStd()
	tf := &TextFormatter{DisableTTY: true}

	testCases := []struct {
		value    string
		expected string
	}{
		{`ba"r`, `ba\\\"r`},
		{`ba'r`, `ba'r`},
	}

	for _, tc := range testCases {
		b, _ := tf.Format(WithField("test", tc.value))
		if !bytes.Contains(b, []byte(tc.expected)) {
			t.Errorf("escaping expected for %q (result was %q instead of %q)", tc.value, string(b), tc.expected)
		}
	}
}

func TestEscaping_Interface(t *testing.T) {
	defer newStd()
	tf := &TextFormatter{DisableTTY: true}

	ts := time.Now()

	testCases := []struct {
		value    interface{}
		expected string
	}{
		{ts.Format(defaultTimestampFormat), ts.Format(defaultTimestampFormat)},
	}

	for _, tc := range testCases {
		b, _ := tf.Format(WithField("test", tc.value))
		if !bytes.Contains(b, []byte(tc.expected)) {
			t.Errorf("escaping expected for %q (result was %q instead of %q)", tc.value, string(b), tc.expected)
		}
	}
}

func TestEscaping_Error(t *testing.T) {
	defer newStd()
	tf := &TextFormatter{DisableTTY: true, DisableHostname: true}

	testCases := []struct {
		value    interface{}
		expected string
	}{
		{errors.New("error: something went wrong"), "time=\"0001-01-01T00:00:00.000Z\" level=\"fatal\" msg=\"\" error=\"error: something went wrong\" caller=\"text_formatter_test.go:94 github.com/bdlm/log/v2.TestEscaping_Error\"\n"},
	}

	for _, tc := range testCases {
		b, _ := tf.Format(WithError(tc.value.(error)))
		if !bytes.Contains(b, []byte(tc.expected)) {
			t.Errorf("escaping expected for %q (result was %q instead of %q)", tc.value, string(b), tc.expected)
		}
	}
}

func TestTimestampFormat(t *testing.T) {
	defer newStd()
	checkTimeStr := func(format string) {
		customFormatter := &TextFormatter{DisableTTY: true, TimestampFormat: format}
		customStr, _ := customFormatter.Format(WithField("test", "test"))
		timeStart := bytes.Index(customStr, ([]byte)("time=\""))
		timeEnd := bytes.Index(customStr, ([]byte)("level="))
		timeStr := customStr[timeStart+5+len("\"") : timeEnd-1-len("\"")]
		if format == "" {
			format = time.RFC3339
		}
		_, e := time.Parse(format, (string)(timeStr))
		if e != nil {
			t.Errorf(`time string '%s' did not match provided time format '%s': %s`, timeStr, format, e)
		}
	}

	checkTimeStr("2006-01-02T15:04:05.000000000Z07:00")
	checkTimeStr("Mon Jan _2 15:04:05 2006")
	checkTimeStr("")
}

//func TestDisableLevelTruncation(t *testing.T) {
//	defer newStd()
//	entry := &Entry{
//		Time:    time.Now(),
//		Message: "testing",
//	}
//	keys := []string{}
//	timestampFormat := "Mon Jan 2 15:04:05 -0700 MST 2006"
//	checkDisableTruncation := func(disabled bool, level Level) {
//		tf := &TextFormatter{DisableLevelTruncation: disabled}
//		var b bytes.Buffer
//		entry.Level = level
//		tf.printColored(&b, entry, keys, timestampFormat)
//		logLine := (&b).String()
//		if disabled {
//			expected := strings.ToUpper(level.String())
//			if !strings.Contains(logLine, expected) {
//				t.Errorf("level string expected to be %s when truncation disabled", expected)
//			}
//		} else {
//			expected := strings.ToUpper(level.String())
//			if len(level.String()) > 4 {
//				if strings.Contains(logLine, expected) {
//					t.Errorf("level string %s expected to be truncated to %s when truncation is enabled", expected, expected[0:4])
//				}
//			} else {
//				if !strings.Contains(logLine, expected) {
//					t.Errorf("level string expected to be %s when truncation is enabled and level string is below truncation threshold", expected)
//				}
//			}
//		}
//	}
//
//	checkDisableTruncation(true, DebugLevel)
//	checkDisableTruncation(true, InfoLevel)
//	checkDisableTruncation(false, ErrorLevel)
//	checkDisableTruncation(false, InfoLevel)
//}

func TestDisableTimestampWithColoredOutput(t *testing.T) {
	defer newStd()
	tf := &TextFormatter{DisableTimestamp: true, ForceTTY: true}

	b, _ := tf.Format(WithField("test", "test"))
	if strings.Contains(string(b), "[0000]") {
		t.Error("timestamp not expected when DisableTimestamp is true")
	}
}

func TestTextFormatterFieldMap(t *testing.T) {
	defer newStd()

	formatter := &TextFormatter{
		DisableTTY:      true,
		DisableHostname: true,
		DisableCaller:   true,
		FieldMap: FieldMap{
			LabelCaller: "caller-label",
			LabelData:   "data-label",
			LabelHost:   "host-label",
			LabelLevel:  "level-label",
			LabelMsg:    "msg-label",
			LabelTime:   "time-field-label",
		},
	}

	entry := &Entry{
		Message: "oh hi",
		Level:   WarnLevel,
		Time:    time.Date(1981, time.February, 24, 4, 28, 3, 100, time.UTC),
		Data: Fields{
			"field1":           "f1",
			"msg-label":        "messagefield",
			"level-label":      "levelfield",
			"time-field-label": "timefield",
		},
	}

	b, err := formatter.Format(entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	assert.Equal(
		t,
		`time-field-label="1981-02-24T04:28:03.000Z" `+
			`level-label="warn" msg-label="oh hi" `+
			`data-label.field1="f1" `+
			`data-label.level-label="levelfield" `+
			`data-label.msg-label="messagefield" `+
			`data-label.time-field-label="timefield"`+"\n",
		string(b),
		"Formatted output doesn't respect FieldMap")
}
