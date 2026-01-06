package cmd

import "fmt"

type Logger struct {
	verbose bool
}

func NewLogger(verbose bool) *Logger {
	return &Logger{verbose: verbose}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.verbose {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}

func (l *Logger) Error(format string, args ...interface{}) {
	fmt.Printf("[ERROR] "+format+"\n", args...)
}
