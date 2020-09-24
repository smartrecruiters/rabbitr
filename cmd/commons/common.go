package commons

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"strings"
	"text/tabwriter"
)

// AbortIfTrue aborts the program execution depending on a provided condition and logs provided message
func AbortIfTrue(condition bool, message string) {
	if condition {
		log.Fatal(message)
	}
}

// AbortIfError aborts the program execution depending on a provided error and logs provided messages
func AbortIfError(err error, messages ...interface{}) {
	if err != nil {
		if len(messages) > 0 {
			log.Fatal(messages...)
		} else {
			log.Fatal(err)
		}
	}
}

// AbortIfErrorWithMsg aborts the program execution depending on a provided error and logs provided message
func AbortIfErrorWithMsg(msg string, err error) {
	if err != nil {
		log.Fatalf(msg, err.Error())
	}
}

// PrintIfTrue prints provided message depending on a condition
func PrintIfTrue(condition bool, msg string) {
	if condition {
		log.Println(msg)
	}
}

// PrintToWriterIfErrorWithMsg prints provided message to a writer depending on a provided error
func PrintToWriterIfErrorWithMsg(w *tabwriter.Writer, msg string, err error) {
	if err != nil {
		Fprintf(w, msg+" %s", err.Error())
	}
}

// PrintIfErrorWithMsg prints provided message depending on a provided error
func PrintIfErrorWithMsg(msg string, err error) {
	if err != nil {
		log.Printf(msg, err.Error())
	}
}

// Fprintln prints provided message to a writer
func Fprintln(w io.Writer, a ...interface{}) {
	_, err := fmt.Fprintln(w, a...)
	PrintIfError(err)
}

// Fprintf prints provided message to a writer
func Fprintf(w io.Writer, format string, a ...interface{}) {
	_, err := fmt.Fprintf(w, format, a...)
	PrintIfError(err)
}

// PrintIfError prints error if it was provided
func PrintIfError(err error) {
	if err != nil {
		log.Printf(err.Error())
	}
}

// IsOSX checks if running system is OSX
func IsOSX() bool {
	return strings.Contains(runtime.GOOS, "darwin")
}
