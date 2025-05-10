package logger

import (
	"log"
	"net/http"
	"os"
)

type Logger struct {
	file *os.File
}

func NewLogger(logFile string) (*Logger, error) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{file: file}, nil
}

func (l *Logger) LogBlockedRequest(r *http.Request, ruleName string) {
	log.SetOutput(l.file)
	log.Printf("BLOCKED: %s %s - %s\n", r.Method, r.URL.String(), ruleName)
}

func (l *Logger) LogError(message string) {
	log.SetOutput(l.file)
	log.Printf("ERROR: %s\n", message)
}
