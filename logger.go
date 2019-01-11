package logger

import (
	"errors"
	"github.com/clareyy/logger/core"
	"github.com/clareyy/logger/file"
)

var (
	defaultLogger *core.Logger
	stdLogger     = core.New(true)
	outPutStd     = false
)

// level 0,1,2,3,4,5
func Init(isStdOut bool, filePath string, level uint32) error {
	if defaultLogger != nil {
		return nil
	}

	outPutStd = isStdOut
	if filePath == "" {
		return errors.New("filePath is empty")
	}

	defaultLogger = core.New(false)
	defaultLogger.SetLevel(core.Level(level))
	w, e := file.NewFileWriter(filePath)
	if e != nil {
		return e
	}

	defaultLogger.SetOutput(w)
	return nil
}

func Fatalf(format string, args ...interface{}) {
	if defaultLogger == nil {
		panic("please init the logger first !")
	}
	if outPutStd {
		stdLogger.Errorf(format, args...)
	}

	if defaultLogger != nil {
		defaultLogger.Fatalf(format, args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if defaultLogger == nil {
		panic("please init the logger first !")
	}
	if defaultLogger != nil {
		defaultLogger.Debugf(format, args...)
	}
	if outPutStd {
		stdLogger.Debugf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if defaultLogger == nil {
		panic("please init the logger first !")
	}
	if defaultLogger != nil {
		defaultLogger.Infof(format, args...)
	}
	if outPutStd {
		stdLogger.Infof(format, args...)
	}
}

func Warningf(format string, args ...interface{}) {
	if defaultLogger == nil {
		panic("please init the logger first !")
	}

	if defaultLogger != nil {
		defaultLogger.Warnf(format, args...)
	}
	if outPutStd {
		stdLogger.Warnf(format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if defaultLogger == nil {
		panic("please init the logger first !")
	}
	if defaultLogger != nil {
		defaultLogger.Errorf(format, args...)
	}
	if outPutStd {
		stdLogger.Errorf(format, args...)
	}
}

func Fatal(args ...interface{}) {
	if defaultLogger == nil {
		panic("please init the logger first !")
	}
	if outPutStd {
		stdLogger.Error(args...)
	}
	if defaultLogger != nil {
		defaultLogger.Fatal(args...)
	}
}

func Debug(args ...interface{}) {
	if defaultLogger == nil {
		panic("please init the logger first !")
	}
	if defaultLogger != nil {
		defaultLogger.Debug(args...)
	}
	if outPutStd {
		stdLogger.Debug(args...)
	}
}

func Info(args ...interface{}) {
	if defaultLogger == nil {
		panic("please init the logger first !")
	}
	if defaultLogger != nil {
		defaultLogger.Info(args...)
	}
	if outPutStd {
		stdLogger.Info(args...)
	}
}

func Warning(args ...interface{}) {
	if defaultLogger == nil {
		panic("please init the logger first !")
	}
	if defaultLogger != nil {
		defaultLogger.Warn(args...)
	}
	if outPutStd {
		stdLogger.Warn(args...)
	}
}

func Error(args ...interface{}) {
	if defaultLogger == nil {
		panic("please init the logger first !")
	}
	if defaultLogger != nil {
		defaultLogger.Error(args...)
	}
	if outPutStd {
		stdLogger.Error(args...)
	}
}

// SetFormatter sets the logger formatter.
func SetFormatter(formatter core.Formatter) {
	if defaultLogger == nil {
		panic("please init the logger first !")
	}

	if defaultLogger != nil {
		defaultLogger.SetFormatter(formatter)
	}
}
