package tui // Package 'tui' for handling text UI tasks

import (
	"fmt"

	tview "github.com/rivo/tview" // External library used for terminal-based user interfaces
)

// InitUrlComponents initializes the UI components for URL input and action buttons (Send, Quit)
func InitUrlComponents(
	app *tview.Application,
	urlForm *tview.Form,
	detailsForm *tview.Form,
	logView *tview.TextView,
	textView *ScrollTextView,
	detailsView *tview.TextView,
) (*tview.Form, *tview.Form) {
	// Dropdown for selecting the HTTP method
	methodbox := tview.NewForm().AddDropDown("", []string{"GET", "POST", "PUT", "DELETE"}, 0, nil)

	// Panel of action buttons - Send and Quit
	buttonPanel := tview.NewForm().
		AddButton("Send", func() {
			SendAction(urlForm, detailsForm, logView, textView, detailsView, methodbox) // Define the function to be called when 'Send' is clicked
			_, method := methodbox.GetFormItem(0).(*tview.DropDown).GetCurrentOption()
			LogMessage(logView, fmt.Sprintf("Method: %s", method)) // Log the selected method when 'Send' is clicked
		}).
		AddButton("Quit", func() {
			app.Stop() // Define the function to be called when 'Quit' is clicked
		})

	return methodbox, buttonPanel
}

// InitUrlandButtons initializes the container for method selection, URL input, and action buttons
func InitUrlandButtons(
	methodbox *tview.Form,
	urlForm *tview.Form,
	buttonPanel *tview.Form,
) *tview.Flex {
	// Create a flexible layout container ("Flex") and add the elements to it.
	urlAndButtons := tview.NewFlex().
		AddItem(methodbox, 9, 1, false).  // Add the methodbox to the Flex
		AddItem(urlForm, 0, 4, true).     // Add the urlForm to the Flex
		AddItem(buttonPanel, 0, 1, false) // Add the buttonPanel to the Flex

	return urlAndButtons
}

// InitGrid initializes the main grid that includes urlAndButtons, htmlPages, textView, and detailsForm
func InitGrid(
	urlAndButtons *tview.Flex,
	htmlPages *tview.Pages,
	textView *ScrollTextView,
	detailsForm *tview.Form,
) *tview.Flex {
	// Set up a new Flex grid that arranges its added items in rows
	grid := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(urlAndButtons, 3, 1, true). // Add the urlAndButtons to the Flex grid
		AddItem(tview.NewFlex().
			AddItem(htmlPages, 35, 1, false). // Add another flex that contains htmlPages, textView, detailsForm
			AddItem(textView, 0, 1, false).
			AddItem(detailsForm, 30, 1, false),
			0, 1, false,
		)

	return grid
}
