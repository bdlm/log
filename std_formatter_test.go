package log

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStdFormatting(t *testing.T) {
	defer newStd()
	tf := &StdFormatter{
		DisableHostname: true,
	}

	testCases := []struct {
		value    string
		expected string
	}{
		{`foo`, "0001/01/01 00:00:00 level=\"fatal\" data.test=\"foo\" caller=\"std_formatter_test.go:27 github.com/bdlm/log/v2.TestStdFormatting\"\n"},
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

func TestStdEscaping(t *testing.T) {
	defer newStd()
	tf := &StdFormatter{}

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

func TestStdEscaping_Interface(t *testing.T) {
	defer newStd()
	tf := &StdFormatter{}

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

func TestStdEscaping_Error(t *testing.T) {
	defer newStd()
	tf := &StdFormatter{}

	testCases := []struct {
		value    interface{}
		expected string
	}{
		{errors.New("error: something went wrong"), "\"error: something went wrong\""},
	}

	for _, tc := range testCases {
		b, _ := tf.Format(WithError(tc.value.(error)))
		if !bytes.Contains(b, []byte(tc.expected)) {
			t.Errorf("escaping expected for %q (result was %q instead of %q)", tc.value, string(b), tc.expected)
		}
	}
}

func TestStdTimestampFormat(t *testing.T) {
	defer newStd()
	checkTimeStr := func(format string) {
		customFormatter := &StdFormatter{TimestampFormat: format}
		if "" == format {
			customFormatter = &StdFormatter{}
			format = "2006/01/02 15:04:05"
		}

		customStr, _ := customFormatter.Format(WithField("test", "test"))
		timeStr := customStr[:len(format)]

		if "" == format {
			timeStr = customStr[:bytes.Index(customStr, []byte(" "))]
		}
		_, e := time.Parse(format, string(timeStr))
		if e != nil {
			t.Errorf(`failed to parse '%s'. format: '%s', err: %s`, timeStr, format, e)
		}
	}

	checkTimeStr("2006-01-02T15:04:05.000000000Z08:00")
	checkTimeStr("Mon Jan _2 15:04:05 2006")
	checkTimeStr("")
}

func TestStdDisableTimestampWithColoredOutput(t *testing.T) {
	defer newStd()
	tf := &StdFormatter{DisableTimestamp: true}

	b, _ := tf.Format(WithField("test", "test"))
	if strings.Contains(string(b), "[0000]") {
		t.Error("timestamp not expected when DisableTimestamp is true")
	}
}

func TestStdTextFormatterFieldMap(t *testing.T) {
	defer newStd()

	formatter := &StdFormatter{
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
			"msg-label":        "messageData",
			"level-label":      "levelData",
			"time-field-label": "timeData",
		},
	}

	b, err := formatter.Format(entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	assert.Equal(
		t,
		`1981/02/24 04:28:03 oh hi `+
			`level-label="warn" `+
			`data-label.field1="f1" `+
			`data-label.level-label="levelData" `+
			`data-label.msg-label="messageData" `+
			`data-label.time-field-label="timeData"`+"\n",
		string(b),
		"Formatted output doesn't respect FieldMap")
}
