package logger

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"discogo/internal/config"
)

const (
	red       = "\033[31m"
	green     = "\033[32m"
	yellow    = "\033[33m"
	cyan      = "\033[36m"
	blue      = "\033[34m"
	lightCyan = "\033[96m"
	reset     = "\033[0m"
)

func Info(message string, showLocation ...bool) {
	showLoc := false
	if len(showLocation) > 0 {
		showLoc = showLocation[0]
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	useColor := config.IsColorEnabled()
	prefix := infoPrefix(now, useColor)
	if showLoc {
		file, line := callerInfo()
		if useColor {
			fmt.Printf("%s [%s%s%s:%s%d%s] %s\n", prefix, blue, file, reset, lightCyan, line, reset, message)
		} else {
			fmt.Printf("%s [%s:%d] %s\n", prefix, file, line, message)
		}
	} else {
		fmt.Printf("%s %s\n", prefix, message)
	}
}

func Error(message string, showLocation ...bool) {
	showLoc := true
	if len(showLocation) > 0 {
		showLoc = showLocation[0]
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	useColor := config.IsColorEnabled()
	prefix := errorPrefix(now, useColor)
	if showLoc {
		file, line := callerInfo()
		if useColor {
			fmt.Fprintf(os.Stderr, "%s [%s%s%s:%s%d%s] %s\n", prefix, blue, file, reset, green, line, reset, message)
		} else {
			fmt.Fprintf(os.Stderr, "%s [%s:%d] %s\n", prefix, file, line, message)
		}
	} else {
		fmt.Fprintf(os.Stderr, "%s %s\n", prefix, message)
	}
}

func infoPrefix(now string, color bool) string {
	if color {
		return fmt.Sprintf("%sINFO:%s %s%s%s", cyan, reset, yellow, now, reset)
	}
	return fmt.Sprintf("INFO: %s", now)
}

func errorPrefix(now string, color bool) string {
	if color {
		return fmt.Sprintf("%sERROR:%s %s%s%s", red, reset, yellow, now, reset)
	}
	return fmt.Sprintf("ERROR: %s", now)
}

func callerInfo() (string, int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "???", 0
	}
	return file, line
}
