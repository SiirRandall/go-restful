// Package tui includes functions to display and handle text based User Interface tasks
package tui

import (
	"fmt" // Standard library used for formatted I/O operations

	"github.com/rivo/tview" // External library used for managing terminal-based user interfaces
)

// AddKeyValueFieldsToForm takes two form parameter and field index.
// It adds key-value fields to the form and updates URL with the parameters when they are modified.
func AddKeyValueFieldsToForm(paramsForm, urlForm *tview.Form, index int, logView *tview.TextView) {
	// Formatting the label of the Key input field
	keyLabel := fmt.Sprintf("┌Key %d:", index)
	valueLabel := "└Value:"

	// Adding a new Key input field and a Value input field to the Params form
	paramsForm.
		// Add key field
		AddInputField(keyLabel, "", 50, func(textToCheck string, lastChar rune) bool {
			// Updating URL whenever there is a change in params
			UpdateURLWithParams(urlForm, paramsForm, logView)
			return true
		}, nil).
		// Add value field
		AddInputField(valueLabel, "", 50, func(textToCheck string, lastChar rune) bool {
			// Updating URL whenever there is a change in params
			UpdateURLWithParams(urlForm, paramsForm, logView)
			return true
		}, nil)
}
