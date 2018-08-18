package log

import (
	"fmt"
)

// STDFormatter formats logs with the default format used by the standard
// library.
type STDFormatter struct {

	// DisableMessage disables message output.
	DisableMessage bool

	// DisableTimestamp disables timestamp output.
	DisableTimestamp bool

	// TimestampFormat allows a custom timestamp format to be used.
	TimestampFormat string
}

// Format renders a single log entry
func (f *STDFormatter) Format(entry *Entry) ([]byte, error) {
	var msg string
	var ts string
	format := "2006/01/02 15:04:05"

	if !f.DisableMessage {
		msg = entry.Message
	}
	if "" != f.TimestampFormat {
		format = f.TimestampFormat
	}
	if !f.DisableTimestamp {
		ts = entry.Time.Format(format) + " "
	}

	return []byte(fmt.Sprintf("%s%s\n", ts, msg)), nil
}
