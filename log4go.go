package log4go

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	Version = "0.0.1"
)

type LEVEL = int

const (
	DEBUG LEVEL = iota
	INFO
	WARN
	ERROR
)

type genericLogger interface {
	Printf(format string, v ...interface{})
}

type Logger struct {
	Logger genericLogger
	level  LEVEL
	name   string
}

func NewLogger(lg genericLogger, level LEVEL, name string) (logger *Logger) {
	logger = &Logger{
		Logger: lg,
		level:  level,
		name:   name,
	}

	return
}

func (self *Logger) format(f string, levelName string) string {
	return fmt.Sprintf("[%s][%s][%s]%s\n", self.name,
		time.Now().Format("2006-01-02 15:04:05"),
		levelName, f)
}

func (self *Logger) Debug(format string, v ...interface{}) {
	if self.level <= DEBUG {
		self.Logger.Printf(self.format(format, "DEBUG"), v...)
	}
}

func (self *Logger) Info(format string, v ...interface{}) {
	if self.level <= INFO {
		self.Logger.Printf(self.format(format, "INFO"), v...)
	}
}

func (self *Logger) Warn(format string, v ...interface{}) {
	if self.level <= WARN {
		self.Logger.Printf(self.format(format, "WARN"), v...)
	}
}

func (self *Logger) Error(format string, v ...interface{}) {
	if self.level <= ERROR {
		self.Logger.Printf(self.format(format, "ERROR"), v...)
	}
}

var Std = NewLogger(log.New(os.Stdout, "xt", log.LstdFlags|log.Lshortfile), DEBUG, "xt")

func Debug(format string, v ...interface{}) {
	Std.Debug(format, v...)
}
func Info(format string, v ...interface{}) {
	Std.Info(format, v...)
}
func Warn(format string, v ...interface{}) {
	Std.Warn(format, v...)
}
func Error(format string, v ...interface{}) {
	Std.Error(format, v...)
}
