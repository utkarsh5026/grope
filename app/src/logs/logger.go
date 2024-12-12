package logs

import (
	"os"

	"github.com/fatih/color"
)

type LogFunc func(format string, args ...any)

type Logger struct {
	errorPrinter   LogFunc
	fatalPrinter   LogFunc
	warnPrinter    LogFunc
	infoPrinter    LogFunc
	successPrinter LogFunc
}

func NewLogger() *Logger {
	return &Logger{
		errorPrinter:   color.New(color.FgRed, color.Bold).PrintfFunc(),
		infoPrinter:    color.New(color.FgBlue).PrintfFunc(),
		successPrinter: color.New(color.FgGreen).PrintfFunc(),
		warnPrinter:    color.New(color.FgYellow).PrintfFunc(),
		fatalPrinter:   color.New(color.FgRed, color.Bold).PrintfFunc(),
	}
}

func (l *Logger) Error(format string, args ...any) {
	l.errorPrinter(format, args...)
}

func (l *Logger) Info(format string, args ...any) {
	l.infoPrinter(format, args...)
}

func (l *Logger) Success(format string, args ...any) {
	l.successPrinter(format, args...)
}

func (l *Logger) Warn(format string, args ...any) {
	l.warnPrinter(format, args...)
}

func (l *Logger) Fatal(format string, args ...any) {
	l.fatalPrinter(format, args...)
	os.Exit(1)
}

var defaultLogger = NewLogger()

func Warn(format string, args ...any) {
	defaultLogger.Warn(format, args...)
}

func Info(format string, args ...any) {
	defaultLogger.Info(format, args...)
}

func Success(format string, args ...any) {
	defaultLogger.Success(format, args...)
}

func Fatal(format string, args ...any) {
	defaultLogger.Fatal(format, args...)
}

func Error(format string, args ...any) {
	defaultLogger.Error(format, args...)
}
