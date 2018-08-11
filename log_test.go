package log

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func LogAndAssertJSON(t *testing.T, log func(*Logger), assertions func(data logData)) {
	var buffer bytes.Buffer
	var data logData

	logger := New()
	logger.SetLevel(DebugLevel)
	logger.Out = &buffer
	logger.Formatter = new(JSONFormatter)

	log(logger)

	err := json.Unmarshal(buffer.Bytes(), &data)
	assert.Nil(t, err)

	assertions(data)
}

func LogAndAssertText(t *testing.T, log func(*Logger), assertions func(fields map[string]string)) {
	var buffer bytes.Buffer

	logger := New()
	logger.Out = &buffer
	logger.Formatter = &TextFormatter{
		DisableColors: true,
	}

	log(logger)

	fields := make(map[string]string)

	re := regexp.MustCompile(`[a-zA-Z0-9\\.]+=\"(\\"|\\" |[\d\w\s!@#$%^&*()_+\-=\[\]{};':\\|,.<>\/?])*(" |"\n|)`)
	for _, kv := range re.FindAllString(buffer.String(), -1) {
		if !strings.Contains(kv, "=") {
			continue
		}
		kvArr := strings.Split(kv, "=")
		key := strings.TrimSpace(kvArr[0])
		val := strings.TrimSpace(kvArr[1])
		if '"' == kvArr[1][0] && "" != string(val) {
			var err error
			val, err = strconv.Unquote(val)
			assert.NoError(t, err)
		}
		fields[key] = val
	}
	assertions(fields)
}

func TestPrint(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Print("test")
	}, func(data logData) {
		assert.Equal(t, "test", data.Message)
		assert.Equal(t, "info", data.Level)
	})
}

func TestInfo(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Info("test")
	}, func(data logData) {
		assert.Equal(t, "test", data.Message)
		assert.Equal(t, "info", data.Level)
	})
}

func TestWarn(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Warn("test")
	}, func(data logData) {
		assert.Equal(t, "test", data.Message)
		assert.Equal(t, "warn", data.Level)
	})
}

func TestInfolnShouldAddSpacesBetweenStrings(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Infoln("test", "test")
	}, func(data logData) {
		assert.Equal(t, "test test", data.Message)
	})
}

func TestInfolnShouldAddSpacesBetweenStringAndNonstring(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Infoln("test", 10)
	}, func(data logData) {
		assert.Equal(t, "test 10", data.Message)
	})
}

func TestInfolnShouldAddSpacesBetweenTwoNonStrings(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Infoln(10, 10)
	}, func(data logData) {
		assert.Equal(t, "10 10", data.Message)
	})
}

func TestInfoShouldAddSpacesBetweenTwoNonStrings(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Infoln(10, 10)
	}, func(data logData) {
		assert.Equal(t, "10 10", data.Message)
	})
}

func TestInfoShouldNotAddSpacesBetweenStringAndNonstring(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Info("test", 10)
	}, func(data logData) {
		assert.Equal(t, "test10", data.Message)
	})
}

func TestInfoShouldNotAddSpacesBetweenStrings(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.Info("test", "test")
	}, func(data logData) {
		assert.Equal(t, "testtest", data.Message)
	})
}

func TestWithFieldsShouldAllowAssignments(t *testing.T) {
	var buffer bytes.Buffer
	var data logData

	logger := New()
	logger.Out = &buffer
	logger.Formatter = new(JSONFormatter)

	localLog := logger.WithFields(Fields{
		"key1": "value1",
	})

	localLog.WithField("key2", "value2").Info("test")
	err := json.Unmarshal(buffer.Bytes(), &data)
	assert.Nil(t, err)

	assert.Equal(t, "value2", data.Data["key2"])
	assert.Equal(t, "value1", data.Data["key1"])

	buffer = bytes.Buffer{}
	data = logData{}
	localLog.Info("test")
	err = json.Unmarshal(buffer.Bytes(), &data)
	assert.Nil(t, err)

	_, ok := data.Data["key2"]
	assert.Equal(t, false, ok)
	assert.Equal(t, "value1", data.Data["key1"])
}

func TestUserSuppliedFieldDoesNotOverwriteDefaults(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.WithField("msg", "hello").Info("test")
	}, func(data logData) {
		assert.Equal(t, "test", data.Message)
	})
}

func TestUserSuppliedMsgFieldHasPrefix(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.WithField("msg", "hello").Info("test")
	}, func(data logData) {
		assert.Equal(t, "test", data.Message)
		assert.Equal(t, "hello", data.Data["msg"])
	})
}

func TestUserSuppliedTimeFieldHasPrefix(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.WithField("time", "hello").Info("test")
	}, func(data logData) {
		assert.Equal(t, "hello", data.Data["time"])
	})
}

func TestUserSuppliedLevelFieldHasPrefix(t *testing.T) {
	LogAndAssertJSON(t, func(log *Logger) {
		log.WithField("level", 1).Info("test")
	}, func(data logData) {
		assert.Equal(t, "info", data.Level)
		assert.Equal(t, 1.0, data.Data["level"]) // JSON has floats only
	})
}

func TestDefaultFieldsAreNotPrefixed(t *testing.T) {
	LogAndAssertText(
		t,
		func(log *Logger) {
			ll := log.WithField("herp", "derp")
			ll.Info("hello")
			ll.Info("bye")
		},
		func(fields map[string]string) {
			fieldMap := &FieldMap{}
			for _, fieldName := range []string{
				fieldMap.resolve(LabelData) + ".level",
				fieldMap.resolve(LabelData) + ".time",
				fieldMap.resolve(LabelData) + ".msg",
			} {
				if _, ok := fields[fieldName]; ok {
					t.Fatalf("should not have prefixed %q: %v", fieldName, fields)
				}
			}
		},
	)
}

func TestWithTimeShouldOverrideTime(t *testing.T) {
	now := time.Now().Add(24 * time.Hour)

	LogAndAssertJSON(t, func(log *Logger) {
		log.WithTime(now).Info("foobar")
	}, func(data logData) {
		assert.Equal(t, now.Format(defaultTimestampFormat), data.Timestamp)
	})
}

func TestWithTimeShouldNotOverrideFields(t *testing.T) {
	now := time.Now().Add(24 * time.Hour)

	LogAndAssertJSON(t, func(log *Logger) {
		log.WithField("herp", "derp").WithTime(now).Info("blah")
	}, func(data logData) {
		assert.Equal(t, now.Format(defaultTimestampFormat), data.Timestamp)
		assert.Equal(t, "derp", data.Data["herp"])
	})
}

func TestWithFieldShouldNotOverrideTime(t *testing.T) {
	now := time.Now().Add(24 * time.Hour)

	LogAndAssertJSON(t, func(log *Logger) {
		log.WithTime(now).WithField("herp", "derp").Info("blah")
	}, func(data logData) {
		assert.Equal(t, now.Format(defaultTimestampFormat), data.Timestamp)
		assert.Equal(t, "derp", data.Data["herp"])
	})
}

func TestTimeOverrideMultipleLogs(t *testing.T) {
	var buffer bytes.Buffer
	var firstFields, secondFields Fields

	logger := New()
	logger.Out = &buffer
	formatter := new(JSONFormatter)
	formatter.TimestampFormat = time.StampMilli
	logger.Formatter = formatter

	llog := logger.WithField("herp", "derp")
	llog.Info("foo")

	err := json.Unmarshal(buffer.Bytes(), &firstFields)
	assert.NoError(t, err, "should have decoded first message")

	buffer.Reset()

	time.Sleep(10 * time.Millisecond)
	llog.Info("bar")

	err = json.Unmarshal(buffer.Bytes(), &secondFields)
	assert.NoError(t, err, "should have decoded second message")

	assert.NotEqual(t, firstFields["time"], secondFields["time"], "timestamps should not be equal")
}

func TestDoubleLoggingDoesntPrefixPreviousFields(t *testing.T) {

	var buffer bytes.Buffer
	var data logData

	logger := New()
	logger.Out = &buffer
	logger.Formatter = new(JSONFormatter)

	llog := logger.WithField("context", "eating raw fish")

	llog.Info("looks delicious")

	err := json.Unmarshal(buffer.Bytes(), &data)
	assert.NoError(t, err, "should have decoded first message")
	assert.Equal(t, "looks delicious", data.Message)
	assert.Equal(t, "eating raw fish", data.Data["context"])

	buffer.Reset()

	llog.Warn("omg it is!")

	err = json.Unmarshal(buffer.Bytes(), &data)
	assert.NoError(t, err, "should have decoded second message")
	assert.Equal(t, "omg it is!", data.Message)
	assert.Equal(t, "eating raw fish", data.Data["context"])
	assert.Nil(t, data.Data["msg"], "should not have prefixed previous `msg` entry")

}

func TestConvertLevelToString(t *testing.T) {
	assert.Equal(t, "debug", DebugLevel.String())
	assert.Equal(t, "info", InfoLevel.String())
	assert.Equal(t, "warn", WarnLevel.String())
	assert.Equal(t, "error", ErrorLevel.String())
	assert.Equal(t, "fatal", FatalLevel.String())
	assert.Equal(t, "panic", PanicLevel.String())
}

func TestParseLevel(t *testing.T) {
	l, err := ParseLevel("panic")
	assert.Nil(t, err)
	assert.Equal(t, PanicLevel, l)

	l, err = ParseLevel("PANIC")
	assert.Nil(t, err)
	assert.Equal(t, PanicLevel, l)

	l, err = ParseLevel("fatal")
	assert.Nil(t, err)
	assert.Equal(t, FatalLevel, l)

	l, err = ParseLevel("FATAL")
	assert.Nil(t, err)
	assert.Equal(t, FatalLevel, l)

	l, err = ParseLevel("error")
	assert.Nil(t, err)
	assert.Equal(t, ErrorLevel, l)

	l, err = ParseLevel("ERROR")
	assert.Nil(t, err)
	assert.Equal(t, ErrorLevel, l)

	l, err = ParseLevel("warn")
	assert.Nil(t, err)
	assert.Equal(t, WarnLevel, l)

	l, err = ParseLevel("WARN")
	assert.Nil(t, err)
	assert.Equal(t, WarnLevel, l)

	l, err = ParseLevel("warning")
	assert.Nil(t, err)
	assert.Equal(t, WarnLevel, l)

	l, err = ParseLevel("WARNING")
	assert.Nil(t, err)
	assert.Equal(t, WarnLevel, l)

	l, err = ParseLevel("info")
	assert.Nil(t, err)
	assert.Equal(t, InfoLevel, l)

	l, err = ParseLevel("INFO")
	assert.Nil(t, err)
	assert.Equal(t, InfoLevel, l)

	l, err = ParseLevel("debug")
	assert.Nil(t, err)
	assert.Equal(t, DebugLevel, l)

	l, err = ParseLevel("DEBUG")
	assert.Nil(t, err)
	assert.Equal(t, DebugLevel, l)

	_, err = ParseLevel("invalid")
	assert.Equal(t, "not a valid log Level: \"invalid\"", err.Error())
}

func TestGetSetLevelRace(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				SetLevel(InfoLevel)
			} else {
				GetLevel()
			}
		}(i)

	}
	wg.Wait()
}

func TestLoggingRace(t *testing.T) {
	logger := New()

	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			logger.Info("info")
			wg.Done()
		}()
	}
	wg.Wait()
}

// Compile test
func TestLogInterface(t *testing.T) {
	var buffer bytes.Buffer
	fn := func(l FieldLogger) {
		b := l.WithField("key", "value")
		b.Debug("Test")
	}
	// test logger
	logger := New()
	logger.Out = &buffer
	fn(logger)

	// test Entry
	e := logger.WithField("another", "value")
	fn(e)
}

// Implements io.Writer using channels for synchronization, so we can wait on
// the Entry.Writer goroutine to write in a non-racey way. This does assume that
// there is a single call to Logger.Out for each message.
type channelWriter chan []byte

func (cw channelWriter) Write(p []byte) (int, error) {
	cw <- p
	return len(p), nil
}

func TestEntryWriter(t *testing.T) {
	cw := channelWriter(make(chan []byte, 1))
	log := New()
	log.Out = cw
	log.Formatter = new(JSONFormatter)
	log.WithField("foo", "bar").WriterLevel(WarnLevel).Write([]byte("hello\n"))

	bs := <-cw
	var data logData
	err := json.Unmarshal(bs, &data)
	assert.Nil(t, err)
	assert.Equal(t, "bar", data.Data["foo"])
	assert.Equal(t, "warn", data.Level)
}

func TestLogSecrets(t *testing.T) {
	AddSecret("my-secret-text")
	AddSecret("secret2")
	AddSecret("secret3")

	LogAndAssertJSON(t, func(log *Logger) {
		log.Debugln("my secret text is 'my-secret-text'. and I know secret2 and secret3")
	}, func(data logData) {
		assert.Equal(t, "my secret text is '**************'. and I know ******* and *******", data.Message)
	})
	LogAndAssertJSON(t, func(log *Logger) {
		log.Infoln("my secret text is 'my-secret-text'. and I know secret2 and secret3")
	}, func(data logData) {
		assert.Equal(t, "my secret text is '**************'. and I know ******* and *******", data.Message)
	})
	LogAndAssertJSON(t, func(log *Logger) {
		log.Warnln("my secret text is 'my-secret-text'. and I know secret2 and secret3")
	}, func(data logData) {
		assert.Equal(t, "my secret text is '**************'. and I know ******* and *******", data.Message)
	})
	LogAndAssertJSON(t, func(log *Logger) {
		log.Errorln("my secret text is 'my-secret-text'. and I know secret2 and secret3")
	}, func(data logData) {
		assert.Equal(t, "my secret text is '**************'. and I know ******* and *******", data.Message)
	})
}
