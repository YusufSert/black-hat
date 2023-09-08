package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

// LoggerEntry is the structure passed to the tamplate.

type LoggerEntry struct {
	StartTime string
	Status    int
	Duration  time.Duration
	Hostname  string
	Method    string
	Path      string
	Request   *http.Request
}

// LoggerDefaultFormat is the format logged used by the default Logger instance
var LoggerDefaultFormat = "{{.StartTime}} | {{.Status}} | \t {{.Duration}} | {{.Hostname}} | {{.Method}} {{.Path}}"

// LoggerDefaultDateFormat is the format used for date by the default Logger instance.
var LoggerDefaultDateFormat = time.RFC3339

// ALogger interface

type ALogger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

// Logger is a middleware  handler that logs the request as it goes in and the response as it goes ot
type Logger struct {
	// ALogger implements just enough log.Logger interface to be compatible with other implementations
	ALogger
	dateFormat string
	template   *template.Template
}

// NewLogger returns a new Logger instance
func NewLogger() *Logger {
	logger := &Logger{ALogger: log.New(os.Stdout, "[gabira] ", 0), dateFormat: LoggerDefaultDateFormat}
	logger.SetFormat(LoggerDefaultFormat)
	return logger
}

func (l *Logger) SetFormat(format string) {
	l.template = template.Must(template.New("gabira_parser").Parse(format))
}

func (l *Logger) SetDateFormat(format string) {
	l.dateFormat = format
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	start := time.Now()
	logcuk := LoggerEntry{
		StartTime: start.Format(l.dateFormat),
		Duration:  time.Since(start),
		Status:    200,
		Hostname:  r.Host,
		Method:    r.Method,
		Path:      r.URL.Path,
		Request:   r,
	}

	buff := &bytes.Buffer{}
	err := l.template.Execute(buff, logcuk)
	if err != nil {
		// do nothing, act like everything ok :)
	}
	l.Println(buff.String())
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", NewLogger()))
}
