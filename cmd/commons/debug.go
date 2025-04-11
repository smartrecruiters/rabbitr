package commons

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

var (
	// IsDebugEnabled is true when the DEBUG env var is not empty.
	IsDebugEnabled = os.Getenv("DEBUG") != ""
	debugLog       = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lmsgprefix)
)

// Debug prints message when debug mode is enabled.
func Debug(message string) {
	printfMsg("%s", message)
}

// DebugIfError prints message when debug mode is enabled and error has occurred.
func DebugIfError(err error) {
	if err != nil {
		printfMsg("%s", err.Error())
	}
}

// Debugf prints message when debug mode is enabled. Substitutes format with provided arguments. Works like fmt.Sprintf.
func Debugf(message string, args ...interface{}) {
	printfMsg(message, args...)
}

// printfMsg prints the message if logging is enabled.
func printfMsg(msg string, v ...interface{}) {
	if IsDebugEnabled {
		debugLog.Print(color.CyanString(fmt.Sprintf(msg, v...)))
	}
}
