package commons

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"strings"
	"text/tabwriter"
)

func AbortIfTrue(condition bool, message string) {
	if condition {
		log.Fatal(message)
	}
}

func AbortIfError(err error, messages ...interface{}) {
	if err != nil {
		if len(messages) > 0 {
			log.Fatal(messages...)
		} else {
			log.Fatal(err)
		}
	}
}

func AbortIfErrorWithMsg(msg string, err error) {
	if err != nil {
		log.Fatalf(msg, err.Error())
	}
}

func PrintIfTrue(condition bool, msg string) {
	if condition {
		log.Println(msg)
	}
}

func PrintToWriterIfErrorWithMsg(w *tabwriter.Writer, msg string, err error) {
	if err != nil {
		Fprintf(w, msg+" %s", err.Error())
	}
}

func PrintIfErrorWithMsg(msg string, err error) {
	if err != nil {
		log.Printf(msg, err.Error())
	}
}

func Fprintln(w io.Writer, a ...interface{}) {
	_, err := fmt.Fprintln(w, a...)
	PrintIfError(err)
}

func Fprintf(w io.Writer, format string, a ...interface{}) {
	_, err := fmt.Fprintf(w, format, a...)
	PrintIfError(err)
}

func PrintIfError(err error) {
	if err != nil {
		log.Printf(err.Error())
	}
}

func IsOSX() bool {
	return strings.Contains(runtime.GOOS, "darwin")
}
