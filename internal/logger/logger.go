package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	env "github.com/tahakara/discogo/internal/config"
)

const (
	red          = "\033[31m"
	green        = "\033[32m"
	yellow       = "\033[33m"
	cyan         = "\033[36m"
	blue         = "\033[34m"
	magenta      = "\033[35m"
	lightCyan    = "\033[96m"
	lightGreen   = "\033[92m"
	lightYellow  = "\033[93m"
	lightRed     = "\033[91m"
	lightBlue    = "\033[94m"
	lightMagenta = "\033[95m"
	reset        = "\033[0m"
)

// Singleline, anlaşılır ve renkli log formatı:
// [LEVEL][YYYY-MM-DD HH:MM:SS][file:line] message

func Info(message string, elapsedTime time.Duration, showLocation ...bool) {
	log("INFO", message, elapsedTime, showLocation...)
}

func Error(message string, elapsedTime time.Duration, showLocation ...bool) {
	log("ERROR", message, elapsedTime, showLocation...)
}

func Debug(message string, elapsedTime time.Duration, showLocation ...bool) {
	log("DEBUG", message, elapsedTime, showLocation...)
}

func Fatal(message string, elapsedTime time.Duration, showLocation ...bool) {
	log("FATAL", message, elapsedTime, showLocation...)
	os.Exit(1)
}

func Register(message string, elapsedTime time.Duration, showLocation ...bool) {
	log("REGISTER", message, elapsedTime, showLocation...)
}

func DeRegister(message string, elapsedTime time.Duration, showLocation ...bool) {
	log("DEREGISTER", message, elapsedTime, showLocation...)
}

func HealthCheck(message string, elapsedTime time.Duration, showLocation ...bool) {
	log("HEALTHCHECK", message, elapsedTime, showLocation...)
}

func HeartBeat(message string, elapsedTime time.Duration, showLocation ...bool) {
	log("HEARTBEAT", message, elapsedTime, showLocation...)
}

func Discovery(message string, elapsedTime time.Duration, showLocation ...bool) {
	log("DISCOVERY", message, elapsedTime, showLocation...)
}

func log(level, message string, elapsedTime time.Duration, showLocation ...bool) {
	showLoc := false
	if len(showLocation) > 0 {
		showLoc = showLocation[0]
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	useColor := env.IsColorEnabled()

	levelStr := level
	switch level {
	case "INFO":
		if useColor {
			levelStr = fmt.Sprintf("%sINFO%s", cyan, reset)
		}
	case "ERROR":
		if useColor {
			levelStr = fmt.Sprintf("%sERROR%s", red, reset)
		}
	case "DEBUG":
		if useColor {
			levelStr = fmt.Sprintf("%sDEBUG%s", yellow, reset)
		}
	case "FATAL":
		if useColor {
			levelStr = fmt.Sprintf("%sFATAL%s", red, reset)
		}
	case "REGISTER":
		if useColor {
			levelStr = fmt.Sprintf("%sREGISTER%s", green, reset)
		}
	case "DEREGISTER":
		if useColor {
			levelStr = fmt.Sprintf("%sDEREGISTER%s", lightGreen, reset)
		}
	case "HEALTHCHECK":
		if useColor {
			levelStr = fmt.Sprintf("%sHEALTHCHECK%s", lightYellow, reset)
		}
	case "HEARTBEAT":
		if useColor {
			levelStr = fmt.Sprintf("%sHEARTBEAT%s", magenta, reset)
		}
	case "DISCOVERY":
		if useColor {
			levelStr = fmt.Sprintf("%sDISCOVERY%s", lightMagenta, reset)
		}
	}

	location := ""
	if showLoc {
		file, line := callerInfo()
		shortFile := filepath.Base(file)
		if useColor {
			location = fmt.Sprintf("[%s%s%s:%s%d%s]", blue, shortFile, reset, lightCyan, line, reset)
		} else {
			location = fmt.Sprintf("[%s:%d]", shortFile, line)
		}
	}

	// Elapsed time ekle
	elapsed := ""
	if elapsedTime != 0 {
		if useColor {
			elapsed = fmt.Sprintf(" [%s%s ms%s]", green, fmt.Sprintf("%d", elapsedTime.Milliseconds()), reset)
		} else {
			elapsed = fmt.Sprintf(" [%s ms]", elapsedTime)
		}
	}

	// Singleline log format
	logLine := fmt.Sprintf("[%s][%s]%s %s%s", levelStr, now, location, message, elapsed)

	if level == "ERROR" || level == "FATAL" {
		fmt.Fprintln(os.Stderr, logLine)
	} else {
		fmt.Println(logLine)
	}
}

func callerInfo() (string, int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "???", 0
	}
	return file, line
}
