package commons

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"
)

var (
	errLog  = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lmsgprefix)
	infoLog = log.New(os.Stdout, "", 0)
)

// AbortIfTrue aborts the program execution depending on a provided condition and logs provided message
func AbortIfTrue(condition bool, message string) {
	if condition {
		errLog.Fatal(message)
	}
}

// AbortIfError aborts the program execution depending on a provided error and logs provided messages
func AbortIfError(err error, messages ...interface{}) {
	if err != nil {
		if len(messages) > 0 {
			errLog.Fatal(messages...)
		} else {
			errLog.Fatal(err)
		}
	}
}

// AbortIfErrorWithMsg aborts the program execution depending on a provided error and logs provided message
func AbortIfErrorWithMsg(msg string, err error) {
	if err != nil {
		errLog.Fatalf(msg, err.Error())
	}
}

// PrintIfTrue prints provided message depending on a condition
func PrintIfTrue(condition bool, msg string) {
	if condition {
		infoLog.Println(msg)
	}
}

// PrintToWriterIfErrorWithMsg prints provided message to a writer depending on a provided error
func PrintToWriterIfErrorWithMsg(w *tabwriter.Writer, msg string, err error) {
	if err != nil {
		Fprintf(w, msg+" %s", err.Error())
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
		errLog.Printf("%s", err.Error())
	}
}

// IsOSX checks if running system is OSX
func IsOSX() bool {
	return strings.Contains(runtime.GOOS, "darwin")
}

// GetStringOrDefault returns value or default value
func GetStringOrDefault(value, defaultValue string) string {
	if value != "" {
		return value
	}
	return defaultValue
}

// GetStringOrDefaultIfValue returns returnValue or default value depending on the valueToCheck
func GetStringOrDefaultIfValue(valueToCheck, returnValue, defaultValue string) string {
	if valueToCheck != "" {
		return returnValue
	}
	return defaultValue
}

// HandleGeneralResponse handles general response from rabbit client.
// Prints response code and error details if able.
func HandleGeneralResponse(messagePrefix string, res *http.Response) {
	if res == nil {
		return
	}
	Fprintf(os.Stdout, messagePrefix+", Response code: %d\t\n", res.StatusCode)
	PrintResponseBodyToWriterIfError(os.Stdout, res)
}

// HandleGeneralResponseWithWriter handles general response from rabbit client.
// Prints response code and error details if able.
// Writes output to the provided writer.
func HandleGeneralResponseWithWriter(w *tabwriter.Writer, res *http.Response) {
	if res == nil {
		return
	}
	Fprintf(w, "Response code: %d\t", res.StatusCode)
	PrintResponseBodyToWriterIfError(w, res)
}

// PrintResponseBodyToWriterIfError prints response body to provided writer
// for responses with code >= 400 to obtain as much error details as possible.
func PrintResponseBodyToWriterIfError(w io.Writer, res *http.Response) {
	if res == nil || res.StatusCode < 400 {
		return
	}
	buf := new(strings.Builder)
	_, err := io.Copy(buf, res.Body)
	if err != nil {
		Fprintf(w, "%s", err.Error())
	}
	Fprintf(w, "%s", buf.String())
}
