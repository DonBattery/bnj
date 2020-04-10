package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/fatih/color"
)

func getDate() string {
	return time.Now().Format("15:04:05")
}

// Debug [DEBUG 17:47:02] Debug message
// Does nothing if debug is not set in Viper
func Debug(msg string) {
	if !viper.GetBool("debug") { // This may be slow, but does not requires any global var
		return
	}
	out := color.New(color.FgHiCyan, color.Bold).Sprint("DEBUG ", getDate())
	fmt.Printf("[%s] %s\n", out, msg)
}

// Debugf formatted Debug message
func Debugf(format string, a ...interface{}) {
	Debug(fmt.Sprintf(format, a...))
}

// Info [INFO 17:47:02] Info message
func Info(msg string) {
	out := color.New(color.FgHiBlue, color.Bold).Sprint("INFO  ", getDate())
	fmt.Printf("[%s] %s\n", out, msg)
}

// Infof formatted Info message
func Infof(format string, a ...interface{}) {
	Info(fmt.Sprintf(format, a...))
}

// Warn [WARN 17:47:02] Warning message
func Warn(msg string) {
	out := color.New(color.FgHiYellow, color.Bold).Sprint("WARN  ", getDate())
	fmt.Printf("[%s] %s\n", out, msg)
}

// Warnf formatted Warn message
func Warnf(format string, a ...interface{}) {
	Warn(fmt.Sprintf(format, a...))
}

// Error [ERROR 17:47:02] Error message
func Error(msg string) {
	out := color.New(color.FgHiRed, color.Bold).Sprint("ERROR ", getDate())
	fmt.Printf("[%s] %s\n", out, msg)
}

// Errorf formatted Error message
func Errorf(format string, a ...interface{}) {
	Error(fmt.Sprintf(format, a...))
}

// Fatal [FATAL 17:18:02] Fatal error message
// exits the program with status code 1
func Fatal(msg string) {
	out := color.New(color.FgHiYellow, color.BgRed, color.Bold).Sprint("FATAL ", getDate())
	fmt.Printf("[%s] %s\n", out, msg)
	os.Exit(1)
}

// Fatalf formatted Fatal error message
func Fatalf(format string, a ...interface{}) {
	Fatal(fmt.Sprintf(format, a...))
}
