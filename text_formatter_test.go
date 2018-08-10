package log

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatting(t *testing.T) {
	tf := &TextFormatter{DisableColors: true, DisableHostname: true}

	testCases := []struct {
		value    string
		expected string
	}{
		{`foo`, "time=\"0001-01-01T00:00:00.000Z\" level=\"panic\" test=\"foo\" caller=\"text_formatter_test.go:25 github.com/bdlm/log.TestFormatting\"\n"},
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
	tf := &TextFormatter{DisableColors: true}

	testCases := []struct {
		value    string
		expected string
	}{
		{`ba"r`, `ba\"r`},
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
	tf := &TextFormatter{DisableColors: true}

	ts := time.Now()

	testCases := []struct {
		value    interface{}
		expected string
	}{
		{ts, fmt.Sprintf("\"%s\"", ts.String())},
		{errors.New("error: something went wrong"), "\"error: something went wrong\""},
	}

	for _, tc := range testCases {
		b, _ := tf.Format(WithField("test", tc.value))
		if !bytes.Contains(b, []byte(tc.expected)) {
			t.Errorf("escaping expected for %q (result was %q instead of %q)", tc.value, string(b), tc.expected)
		}
	}
}

func TestTimestampFormat(t *testing.T) {
	checkTimeStr := func(format string) {
		customFormatter := &TextFormatter{DisableColors: true, TimestampFormat: format}
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
	tf := &TextFormatter{DisableTimestamp: true, ForceColors: true}

	b, _ := tf.Format(WithField("test", "test"))
	if strings.Contains(string(b), "[0000]") {
		t.Error("timestamp not expected when DisableTimestamp is true")
	}
}

func disabledTestTextFormatterFieldMap(t *testing.T) {

	formatter := &TextFormatter{
		DisableColors: true,
		FieldMap: FieldMap{
			LabelMsg:   "message",
			LabelLevel: "somelevel",
			LabelTime:  "timeywimey",
		},
	}

	entry := &Entry{
		Message: "oh hi",
		Level:   WarnLevel,
		Time:    time.Date(1981, time.February, 24, 4, 28, 3, 100, time.UTC),
		Data: Fields{
			"field1":     "f1",
			"message":    "messagefield",
			"somelevel":  "levelfield",
			"timeywimey": "timeywimeyfield",
		},
	}

	b, err := formatter.Format(entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	assert.Equal(t,
		`timeywimey="1981-02-24T04:28:03Z" `+
			`somelevel=warning `+
			`message="oh hi" `+
			`field1=f1 `+
			`fields.message=messagefield `+
			`fields.somelevel=levelfield `+
			`fields.timeywimey=timeywimeyfield`+"\n",
		string(b),
		"Formatted output doesn't respect FieldMap")
}

// TODO add tests for sorting etc., this requires a parser for the text
// formatter output.
