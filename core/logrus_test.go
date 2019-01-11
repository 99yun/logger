package core

import (
	"bytes"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)


func LogAndAssertText(t *testing.T, log func(*Logger), assertions func(fields map[string]string)) {
	var buffer bytes.Buffer

	logger := New(true)
	logger.Out = &buffer
	logger.Formatter = &TextFormatter{
		DisableColors: true,
	}

	log(logger)

	fields := make(map[string]string)
	for _, kv := range strings.Split(buffer.String(), " ") {
		if !strings.Contains(kv, "=") {
			continue
		}
		kvArr := strings.Split(kv, "=")
		key := strings.TrimSpace(kvArr[0])
		val := kvArr[1]
		if kvArr[1][0] == '"' {
			var err error
			val, err = strconv.Unquote(val)
			assert.NoError(t, err)
		}
		fields[key] = val
	}
	assertions(fields)
}

func TestConvertLevelToString(t *testing.T) {
	assert.Equal(t, "debug", DebugLevel.String())
	assert.Equal(t, "info", InfoLevel.String())
	assert.Equal(t, "warning", WarnLevel.String())
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

	l, err = ParseLevel("invalid")
	assert.Equal(t, "not a valid logrus Level: \"invalid\"", err.Error())
}


func TestLoggingRace(t *testing.T) {
	logger := New(true)

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
func TestLogrusInterface(t *testing.T) {
	var buffer bytes.Buffer
	fn := func(l FieldLogger) {
		b := l.WithField("key", "value")
		b.Debug("Test")
	}
	// test logger
	logger := New(true)
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

func TestLogLevelEnabled(t *testing.T) {
	log := New(true)
	log.SetLevel(PanicLevel)
	assert.Equal(t, true, log.IsLevelEnabled(PanicLevel))
	assert.Equal(t, false, log.IsLevelEnabled(FatalLevel))
	assert.Equal(t, false, log.IsLevelEnabled(ErrorLevel))
	assert.Equal(t, false, log.IsLevelEnabled(WarnLevel))
	assert.Equal(t, false, log.IsLevelEnabled(InfoLevel))
	assert.Equal(t, false, log.IsLevelEnabled(DebugLevel))

	log.SetLevel(FatalLevel)
	assert.Equal(t, true, log.IsLevelEnabled(PanicLevel))
	assert.Equal(t, true, log.IsLevelEnabled(FatalLevel))
	assert.Equal(t, false, log.IsLevelEnabled(ErrorLevel))
	assert.Equal(t, false, log.IsLevelEnabled(WarnLevel))
	assert.Equal(t, false, log.IsLevelEnabled(InfoLevel))
	assert.Equal(t, false, log.IsLevelEnabled(DebugLevel))

	log.SetLevel(ErrorLevel)
	assert.Equal(t, true, log.IsLevelEnabled(PanicLevel))
	assert.Equal(t, true, log.IsLevelEnabled(FatalLevel))
	assert.Equal(t, true, log.IsLevelEnabled(ErrorLevel))
	assert.Equal(t, false, log.IsLevelEnabled(WarnLevel))
	assert.Equal(t, false, log.IsLevelEnabled(InfoLevel))
	assert.Equal(t, false, log.IsLevelEnabled(DebugLevel))

	log.SetLevel(WarnLevel)
	assert.Equal(t, true, log.IsLevelEnabled(PanicLevel))
	assert.Equal(t, true, log.IsLevelEnabled(FatalLevel))
	assert.Equal(t, true, log.IsLevelEnabled(ErrorLevel))
	assert.Equal(t, true, log.IsLevelEnabled(WarnLevel))
	assert.Equal(t, false, log.IsLevelEnabled(InfoLevel))
	assert.Equal(t, false, log.IsLevelEnabled(DebugLevel))

	log.SetLevel(InfoLevel)
	assert.Equal(t, true, log.IsLevelEnabled(PanicLevel))
	assert.Equal(t, true, log.IsLevelEnabled(FatalLevel))
	assert.Equal(t, true, log.IsLevelEnabled(ErrorLevel))
	assert.Equal(t, true, log.IsLevelEnabled(WarnLevel))
	assert.Equal(t, true, log.IsLevelEnabled(InfoLevel))
	assert.Equal(t, false, log.IsLevelEnabled(DebugLevel))

	log.SetLevel(DebugLevel)
	assert.Equal(t, true, log.IsLevelEnabled(PanicLevel))
	assert.Equal(t, true, log.IsLevelEnabled(FatalLevel))
	assert.Equal(t, true, log.IsLevelEnabled(ErrorLevel))
	assert.Equal(t, true, log.IsLevelEnabled(WarnLevel))
	assert.Equal(t, true, log.IsLevelEnabled(InfoLevel))
	assert.Equal(t, true, log.IsLevelEnabled(DebugLevel))
}
