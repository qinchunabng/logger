package logger

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"
)

type LogLevel uint16

const (
	UNKNOWN LogLevel = iota
	DEBUG
	TRACE
	INFO
	WARN
	ERROR
	FATAL
)

// 日志结构体
type ConsoleLogger struct {
	Level LogLevel
}

func parseLogLevel(s string) (LogLevel, error) {
	switch strings.ToLower(s) {
	case "debug":
		return DEBUG, nil
	case "trace":
		return TRACE, nil
	case "info":
		return INFO, nil
	case "warn":
		return WARN, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	default:
		return UNKNOWN, errors.New("无效的日志级别")
	}
}

// Logger构造函数
func NewLog(level string) ConsoleLogger {
	l, err := parseLogLevel(level)
	if err != nil {
		panic(err)
	}
	return ConsoleLogger{l}
}

func (log *ConsoleLogger) logPrint(format, level string, args ...interface{}) {
	l, err := parseLogLevel(level)
	if err != nil {
		panic(err)
	}
	if log.Level > l {
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05.000")
	funcName, fileName, lineNo := getInfo(3)
	s := fmt.Sprintf(format, args...)
	fmt.Printf("[%s] [%s] [%s:%s:%d] %s\n", now, level, fileName, funcName, lineNo, s)
}

func (log *ConsoleLogger) Trace(format string, args ...interface{}) {
	log.logPrint(format, "TRACE", args...)
}

func (log *ConsoleLogger) Debug(format string, args ...interface{}) {
	log.logPrint(format, "DEBUG", args...)
}

func (log *ConsoleLogger) Info(format string, args ...interface{}) {
	log.logPrint(format, "INFO", args...)
}

func (log *ConsoleLogger) Warn(format string, args ...interface{}) {
	log.logPrint(format, "WARN", args...)
}

func (log *ConsoleLogger) Error(format string, args ...interface{}) {
	log.logPrint(format, "ERROR", args...)
}

func (log *ConsoleLogger) Fatal(format string, args ...interface{}) {
	log.logPrint(format, "FATAL", args...)
}

func getInfo(skip int) (funcName, fileName string, lineNo int) {
	pc, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		fmt.Printf("runtime.Caller() failed\n")
	}
	funcName = runtime.FuncForPC(pc).Name()
	funcName = strings.Split(funcName, ".")[1]
	fileName = path.Base(file)
	return
}
