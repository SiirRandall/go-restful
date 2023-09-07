// Package tui includes functions to display and handle text based User Interface tasks
package tui

import (
	"fmt"
	"time"

	"github.com/rivo/tview" // External library used for managing terminal-based user interfaces
)

// LogMessage logs a new message to the provided TextView.
// It uses text formatting to bold the timestamp and automatically scrolls to the end of the log.
func LogMessage(logView *tview.TextView, msg string) {
	// Format and print the message with a timestamp at the beginning
	fmt.Fprintf(logView, "[::b]%s[::-]: %s\n", getTimeStamp(), msg)
	logView.ScrollToEnd() // Scroll the TextView to show the latest log message
}

// getTimeStamp returns the current time as a formatted string in HH:MM:SS format
func getTimeStamp() string {
	return time.Now().Format("15:04:05") // Using Go's standard time formatting layout
}
