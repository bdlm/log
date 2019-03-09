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
	SetCallerLevel(0)
}
