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
	return []byte(fmt.Sprintf("%s %s\n", entry.Time.Format("2006/01/02 15:04:05"), entry.Message)), nil
}
