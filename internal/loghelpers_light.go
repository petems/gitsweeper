package internal

import (
	"log"
	"os"
)

var (
	debugEnabled = false
	logger       = log.New(os.Stderr, "", log.LstdFlags)
)

func SetupLightLogger(debug bool) {
	debugEnabled = debug
	if debug {
		logger.SetPrefix("[DEBUG] ")
		logger.Println("--debug setting detected - Info level logs enabled")
	} else {
		logger.SetPrefix("")
	}
}

func LogInfo(msg string) {
	if debugEnabled {
		logger.Println("[INFO]", msg)
	}
}

func LogInfof(format string, args ...interface{}) {
	if debugEnabled {
		logger.Printf("[INFO] "+format, args...)
	}
}

func LogFatal(msg string) {
	logger.Println("[FATAL]", msg)
	os.Exit(1)
}

func LogFatalError(msg string, err error) {
	if err != nil {
		logger.Printf("[FATAL] %s: %v", msg, err)
	} else {
		logger.Println("[FATAL]", msg)
	}
	os.Exit(1)
}
