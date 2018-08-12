package log

import (
	"bufio"
	"io"
	"runtime"

	stdLogger "github.com/bdlm/std/logger"
)

// Writer returns an info-level log writer.
func (logger *Logger) Writer() *io.PipeWriter {
	return logger.WriterLevel(InfoLevel)
}

// WriterLevel returns a log writer with a specified leve.
func (logger *Logger) WriterLevel(level stdLogger.Level) *io.PipeWriter {
	return NewEntry(logger).WriterLevel(level)
}

// Writer returns an info-level log writer.
func (entry *Entry) Writer() *io.PipeWriter {
	return entry.WriterLevel(InfoLevel)
}

// WriterLevel returns a log writer with a specified leve.
func (entry *Entry) WriterLevel(level stdLogger.Level) *io.PipeWriter {
	reader, writer := io.Pipe()

	var printFunc func(args ...interface{})

	switch level {
	case DebugLevel:
		printFunc = entry.Debug
	case InfoLevel:
		printFunc = entry.Info
	case WarnLevel:
		printFunc = entry.Warn
	case ErrorLevel:
		printFunc = entry.Error
	case FatalLevel:
		printFunc = entry.Fatal
	case PanicLevel:
		printFunc = entry.Panic
	default:
		printFunc = entry.Print
	}

	go entry.writerScanner(reader, printFunc)
	runtime.SetFinalizer(writer, writerFinalizer)

	return writer
}

func (entry *Entry) writerScanner(reader *io.PipeReader, printFunc func(args ...interface{})) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		printFunc(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		entry.Errorf("Error while reading from Writer: %s", err)
	}
	reader.Close()
}

func writerFinalizer(writer *io.PipeWriter) {
	writer.Close()
}
