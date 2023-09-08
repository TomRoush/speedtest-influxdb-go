package log

import (
	"fmt"
	"os"
	"strings"
)

var logLevel = ErrorLevel

func SetLevel(level Level) {
	logLevel = level
}

func Debug(format string, args ...any) {
	print(DebugLevel, format, args...)
}

func Info(format string, args ...any) {
	print(InfoLevel, format, args...)
}

func Error(format string, args ...any) {
	print(ErrorLevel, format, args...)
}

func print(level Level, format string, args ...any) {
	if logLevel <= level {
		msg := fmt.Sprintf(format, args...)
		fmt.Fprintf(getOutputStream(level), "[%s] %s\n", strings.ToUpper(level.ToString()), msg)
	}
}

func getOutputStream(level Level) *os.File {
	if (level >= WarnLevel) {
		return os.Stderr
	} else {
		return os.Stdout
	}
}
