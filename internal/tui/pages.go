// Package tui includes functions to display and handle text based User Interface tasks
package tui

import (
	tcell "github.com/gdamore/tcell/v2" // External library used for handling terminal cell views
	tview "github.com/rivo/tview"       // External library used for managing terminal-based user interfaces
)

// InitHTMLPages initializes a Pages structure with several Form pages and an automatic keyboard input
// handler for switching between these pages via the Tab key.
func InitHTMLPages(
	paramsForm, headersForm, bodyForm, tokenForm *tview.Form, // Input forms for different request components
	logView *tview.TextView, // TextView used for logging activities
) *tview.Pages {
	// Create a new Pages object and add several Form pages to it
	htmlPages := tview.NewPages().
		AddPage("Params", paramsForm, true, true).
		AddPage("Headers", headersForm, true, false).
		AddPage("Body", bodyForm, true, false).
		AddPage("Token", tokenForm, true, false)

	// List of page names in the same order they were added to htmlPages
	requestsPageNames := []string{"Params", "Headers", "Body", "Token"}
	currentPageIndex := 0 // Index of the currently shown page

	// Set an input capture function that switches to the next page when the Tab key is pressed
	htmlPages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			// Cycle through the page indices
			currentPageIndex = (currentPageIndex + 1) % len(requestsPageNames)
			// Switch to the next page
			htmlPages.SwitchToPage(requestsPageNames[currentPageIndex])
		}
		// Log a message if htmlPages has the focus
		if htmlPages.HasFocus() {
			LogMessage(logView, "htmlPages has focus")
		}
		return event
	})

	return htmlPages
}
