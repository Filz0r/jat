package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type timestampWriter struct {
	w      io.Writer
	prefix string
}

func (tw timestampWriter) Write(p []byte) (n int, err error) {
	ts := time.Now().UTC().Format(time.RFC3339)
	line := fmt.Sprintf("%s%s%s %s%s", yellow, ts, reset, tw.prefix, p)
	return tw.w.Write([]byte(line))
}

func newServerLogger() *log.Logger {
	return log.New(
		timestampWriter{w: log.Writer(), prefix: blue + "[jat-server]: " + reset},
		"",
		0,
	)
}

const (
	reset  = "\x1b[0m"
	red    = "\x1b[31m"
	green  = "\x1b[32m"
	yellow = "\x1b[33m"
	blue   = "\x1b[34m"
	gray   = "\x1b[90m"
	purple = "\x1b[35m"
)

func statusColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return blue
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{ResponseWriter: w, status: http.StatusOK}
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Write(p []byte) (int, error) {
	n, err := r.ResponseWriter.Write(p)
	r.bytes += n
	return n, err
}
