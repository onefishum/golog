package golog

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type LevelType uint8

const (
	LevelDebug LevelType = iota
	LevelInfo
	LevelWarn
	LevelError
)

const (
	levelDebugMsg = "Debug"
	levelInfoMsg  = "Info"
	levelWarnMsg  = "Warn"
	levelErrorMsg = "Error"
)

const (
	defaultDateFormat           = "2006-01-02"
	defaultDateTimeFormat       = "2006-01-02 15:04:05.000"
	defaultLogFormatPrefixPrint = "[%-5s] [%s] : %s -> %s \n"
)

var (
	pool *sync.Pool
)

func init() {
	pool = &sync.Pool{
		New: func() interface{} {
			return &logEntity{}
		},
	}
}

type config struct {
	dateFormat     string
	dateTimeFormat string
}

func NewConfig() *config {
	return &config{
		dateFormat:     defaultDateFormat,
		dateTimeFormat: defaultDateTimeFormat,
	}
}

func (cfg *config) SetDateFormat(dateFormat string) {
	cfg.dateFormat = dateFormat
}

func (cfg *config) SetDateTimeFormat(dateTimeFormat string) {
	cfg.dateTimeFormat = dateTimeFormat
}

type IPrinter interface {
	Print(level LevelType, str string) error
}

type ILogWriter interface {
	Write(*logEntity) error
	Close() error
}

type Logger struct {
	*config
	level   LevelType
	writer  ILogWriter
	printer IPrinter
}

type logEntity struct {
	msg    string
	level  LevelType
	time   time.Time
	caller string
}

func NewLogger(level LevelType, writer ILogWriter) *Logger {
	logger := &Logger{
		level:   level,
		writer:  writer,
		printer: NewPrinter(),
		config:  NewConfig(),
	}

	return logger
}

func (log *Logger) doLog(level LevelType, msg string, args ...interface{}) {
	fMsg := fmt.Sprintf(msg, args...)
	t := time.Now()
	caller := getFuncCaller(3)
	if level >= log.level {
		str := fmt.Sprintf(defaultLogFormatPrefixPrint, getLevelFlagMsg(level), log.getDateTimeStr(t), caller, fMsg)
		if err := log.printer.Print(level, str); err != nil {
			fmt.Print(str)
		}
	}
	if log.writer != nil {
		le := pool.Get().(*logEntity)
		le.msg = fMsg
		le.level = level
		le.time = t
		le.caller = caller
		if err := log.writer.Write(le); err != nil {
			fmt.Println("[Write Log Error] :", err)
		}
	}
}

func (log *Logger) Debug(msg string, args ...interface{}) {
	log.doLog(LevelDebug, msg, args...)
}

func (log *Logger) Info(msg string, args ...interface{}) {
	log.doLog(LevelInfo, msg, args...)
}

func (log *Logger) Warn(msg string, args ...interface{}) {
	log.doLog(LevelWarn, msg, args...)
}

func (log *Logger) Error(msg string, args ...interface{}) {
	log.doLog(LevelError, msg, args...)
}

func (log *Logger) Close() error {
	if log.writer != nil {
		return log.writer.Close()
	}
	return nil
}

func (log *Logger) getDateStr(t time.Time) string {
	return t.Format(log.dateFormat)
}

func (log *Logger) getDateTimeStr(t time.Time) string {
	return t.Format(log.dateTimeFormat)
}

func getLevelFlagMsg(level LevelType) string {
	switch level {
	case LevelDebug:
		return levelDebugMsg
	case LevelInfo:
		return levelInfoMsg
	case LevelWarn:
		return levelWarnMsg
	case LevelError:
		return levelErrorMsg
	default:
		return ""
	}
}

func getFuncCaller(n int) string {
	_, file, line, _ := runtime.Caller(n)
	return file + ":" + strconv.Itoa(line)
}
