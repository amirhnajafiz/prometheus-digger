package logger

import (
	"fmt"
	"log"
	"strings"
)

// colorize takes a message and returns it with each word colored
// in a different color.
func colorize(msg string) string {
	colors := []string{
		"\033[31m", // Red
		"\033[32m", // Green
		"\033[33m", // Yellow
		"\033[34m", // Blue
		"\033[35m", // Magenta
		"\033[36m", // Cyan
	}
	reset := "\033[0m"

	words := strings.Split(msg, " ")
	for i, word := range words {
		color := colors[i%len(colors)]
		words[i] = fmt.Sprintf("%s%s%s", color, word, reset)
	}
	return strings.Join(words, " ")
}

// Info logs a message at the INFO level.
func Info(msg string) {
	log.Printf("[INFO] %s\n", colorize(msg))
}

// Error logs a message at the ERROR level.
func Error(msg string) {
	log.Printf("[ERROR] %s\n", colorize(msg))
}
