package log

import (
	"errors"
	"strings"
	"testing"
)

func TestSetCallerLevel(t *testing.T) {
	formatter := &JSONFormatter{}
	entry := WithField("error", errors.New("wild walrus"))

	data := getData(entry, formatter.FieldMap, formatter.EscapeHTML, false)
	if !strings.Contains(data.Caller, "formatter_test.go:13") {
		t.Fatal("invalid caller: ", data.Caller)
	}

	SetCallerLevel(-1)
	data = getData(entry, formatter.FieldMap, formatter.EscapeHTML, false)
	if strings.Contains(data.Caller, "formatter_test.go:13") {
		t.Fatal("invalid caller: ", data.Caller)
	}

	// Test error fallback
	SetCallerLevel(10)
	data = getData(entry, formatter.FieldMap, formatter.EscapeHTML, false)
	if !strings.Contains(data.Caller, "formatter_test.go:26") {
		t.Fatal("invalid caller: ", data.Caller)
	}
	SetCallerLevel(0)
}

func TestLevelColors(t *testing.T) {
	formatter := &JSONFormatter{ForceTTY: true}
	SetLevel(DebugLevel)
	entry := WithField("debug", errors.New("wild walrus"))
	entry.Level = DebugLevel
	data := getData(entry, formatter.FieldMap, formatter.EscapeHTML, true)
	if DEBUGColor != data.Color.Level {
		t.Fatal("invalid color: ", data.Color.Level)
	}

	entry = WithField("info", errors.New("wild walrus"))
	entry.Level = InfoLevel
	data = getData(entry, formatter.FieldMap, formatter.EscapeHTML, true)
	if DEFAULTColor != data.Color.Level {
		t.Fatal("invalid color: ", data.Color.Level)
	}

	entry = WithField("warn", errors.New("wild walrus"))
	entry.Level = WarnLevel
	data = getData(entry, formatter.FieldMap, formatter.EscapeHTML, true)
	if WARNColor != data.Color.Level {
		t.Fatal("invalid color: ", data.Color.Level)
	}

	entry = WithField("error", errors.New("wild walrus"))
	entry.Level = ErrorLevel
	data = getData(entry, formatter.FieldMap, formatter.EscapeHTML, true)
	if ERRORColor != data.Color.Level {
		t.Fatal("invalid color: ", data.Color.Level)
	}

	entry = WithField("panic", errors.New("wild walrus"))
	entry.Level = PanicLevel
	data = getData(entry, formatter.FieldMap, formatter.EscapeHTML, true)
	if PANICColor != data.Color.Level {
		t.Fatal("invalid color: ", data.Color.Level)
	}

	entry = WithField("fatal", errors.New("wild walrus"))
	entry.Level = FatalLevel
	data = getData(entry, formatter.FieldMap, formatter.EscapeHTML, true)
	if FATALColor != data.Color.Level {
		t.Fatal("invalid color: ", data.Color.Level)
	}
}
