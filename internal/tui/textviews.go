// Package tui holds components for terminal-based user interfaces
package tui

import (
	"github.com/rivo/tview" // Importing the library that provides tools for building rich terminal applications
)

// InitJsonViewer initializes a new scrollable text view with dynamic colors,
// a border, a title and disables text wrapping. It's designed to display JSON data.
func InitJsonViewer() *ScrollTextView {
	textView := NewScrollTextView()  // Creating new scrollable text view instance
	textView.SetDynamicColors(true)  // Enabling dynamic colors to display different info levels
	textView.SetBorder(true)         // Adding border around the text view
	textView.SetTitle("JSON Viewer") // Setting title for text view
	textView.SetWrap(false)          // Disabling text wrapping for better JSON formatting

	return textView // Returns the configured text view
}

// InitLogView initializes a new scrollable text view with dynamic colors,
// a border, and a title. It's designed to display application logs.
func InitLogView() *tview.TextView {
	logView := tview.NewTextView() // Creating new text view instance
	logView.SetDynamicColors(true) // Enabling dynamic colors to display various log levels
	logView.SetScrollable(true)    // Making the text view scrollable
	logView.SetBorder(true)        // Adding border around the text view
	logView.SetTitle("Logs")       // Setting title for text view

	return logView // Returns the configured text view
}

// InitDetailsView initializes a new scrollable text view with dynamic colors,
// no border, and a title. It's designed to provide more detailed information on screen.
func InitDetailsView() *tview.TextView {
	detailsView := tview.NewTextView() // Creating new text view instance
	detailsView.SetDynamicColors(true) // Enabling dynamic colors to display details in visually differentiated manner
	detailsView.SetScrollable(true)    // Making the text view scrollable
	detailsView.SetBorder(false)       // No border for this view
	detailsView.SetTitle("Details")    // Setting title for text view

	return detailsView // Returns the configured text view
}
