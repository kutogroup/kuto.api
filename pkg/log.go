package pkg

import (
	"io"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

//KutoLogger 日志结构体
type KutoLogger struct {
	logger *logrus.Logger
	debug  bool
}

//NewLogger 新建日志对象
func NewLogger(out io.Writer, debugMode bool) *KutoLogger {
	logger := &KutoLogger{
		logger: logrus.New(),
		debug:  debugMode,
	}

	logger.logger.SetOutput(out)
	if debugMode {
		logger.logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.logger.SetLevel(logrus.WarnLevel)
	}

	return logger
}

//I 输出info日志
func (logger *KutoLogger) I(msg string, args ...interface{}) {
	if !logger.debug {
		//非debug模式，不执行
		return
	}

	_, file, line, _ := runtime.Caller(1)
	file = path.Base(file)
	logger.logger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Infof(msg, args...)
}

//E 输出error日志
func (logger *KutoLogger) E(msg string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	file = path.Base(file)
	logger.logger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Errorf(msg, args...)
}
