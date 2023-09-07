package tui // Package 'tui' for handling text UI tasks

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// SetAutoCompleteForHeaders is a function that sets autocomplete functionality for header input fields
func SetAutoCompleteForHeaders(input *tview.InputField) {
	// By default, do not show all headers
	showAll := false

	// Setting the autocomplete function for the input field
	input.SetAutocompleteFunc(func(currentText string) (entries []string) {
		if showAll { // If all headers are to be shown...
			showAll = false
			return headers // Return all headers
		} else if currentText == "" { // If the current text in the input field is empty...
			return nil // Return no entries
		}

		// For each header, check if the current text in the input field is its prefix.
		for _, header := range headers {
			if strings.HasPrefix(strings.ToLower(header), strings.ToLower(currentText)) {
				entries = append(entries, header) // If it is, add it to the autocomplete entries
			}
		}
		return
	})

	// Input capture to detect the down arrow key
	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// If the key pressed is the down arrow and there's no current text in the input field...
		if event.Key() == tcell.KeyDown && input.GetText() == "" {
			showAll = true
			input.Autocomplete() // Trigger the autocomplete function

			// If headers is not empty, set the first header as the current text.
			if len(headers) > 0 {
				input.SetText(headers[0])
			}

			return nil // Consume the event to prevent default behavior
		}
		return event
	})
}

// SetAutoCompleteForValues is a function that sets autocomplete functionality for value input fields based on the entered header key
func SetAutoCompleteForValues(input *tview.InputField, keyInput *tview.InputField) {
	// Setting the autocomplete function for the input field
	input.SetAutocompleteFunc(func(currentText string) (entries []string) {
		// Get the current value of the key input
		headerKey := keyInput.GetText()

		// Check if we have predefined values for this header
		if possibleValues, exists := headerValuesMap[headerKey]; exists {
			// If predefined values exist, check if the current text is the prefix for any predefined value
			for _, value := range possibleValues {
				if strings.HasPrefix(strings.ToLower(value), strings.ToLower(currentText)) {
					entries = append(entries, value) // If it is, add it to the autocomplete entries
				}
			}
		}
		return
	})
}
