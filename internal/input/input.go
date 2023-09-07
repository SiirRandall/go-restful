package input // Package 'input' handles mouse and keyboard captures

import (
	"fmt"

	"github.com/gdamore/tcell/v2" // For terminal cells
	"github.com/rivo/tview"       // For TUI views

	"github.com/SiirRandall/go-restful/internal/tui" // Internal import of tui
)

// Global boolean variable that keeps track if a logView is visible or not
var isVisible = false

// TextViewMouseCapture handels mouse captures on the TextView
func TextViewMouseCapture(
	logView *tview.TextView, // The logView where events will be logged to
	textView *tui.ScrollTextView, // The textView on which the click event happened
	headersForm *tview.Form, // The form which contains header fields
) {
	textView.SetMouseCapture( // Sets a handler which gets executed when a mouse event occurs on the textView
		func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
			// Check if it's a left button click
			if buttons := event.Buttons(); buttons&tcell.Button1 != 0 {
				x, y := event.Position() // Get position of the click in X, Y format
				// Log the captured mouse click to logView
				tui.LogMessage(
					logView,
					fmt.Sprintf(
						"from input package Mouse clicked in textView at x: %d, y: %d",
						x,
						y,
					),
				)
				hl := headersForm.GetFormItemCount() // Get the count of items in headersForm
				tui.LogMessage(logView, fmt.Sprintf("Items in headersForm %d", hl))
				headersfocusitem, headersfocusbutton := headersForm.GetFocusedItemIndex() // Get index of focused item and button in headersForm
				// Log the focused item and button from headersForm
				tui.LogMessage(
					logView,
					fmt.Sprintf(
						"From input package headersForm focus item %d focus button %d",
						headersfocusitem,
						headersfocusbutton,
					),
				)
			}
			return action, event // Return the unchanged action and event so they can continue being processed
		},
	)
}

// TextViewKBCapture handels input captures (keypresses) on the TextView
func TextViewKBCapture(
	app *tview.Application, // TUI application instance
	textView *tui.ScrollTextView, // The textView on which the keypress event occurred
	logView *tview.TextView, // The logView to log events
	detailsForm *tview.Form, // The form which contains detail fields
	grid *tview.Flex, // The grid layout container
) {
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey { // Sets a handler function that gets executed whenever a key press occurs.
		// Log the captured key to logView
		tui.LogMessage(
			logView,
			fmt.Sprintf(
				"Key captured in textView: %v, Modifiers: %v",
				event.Key(),
				event.Modifiers(),
			),
		)

		if event.Rune() == 'l' { // If "l" was pressed
			if isVisible {
				grid.RemoveItem(logView) // Remove logView from grid
			} else {
				grid.AddItem(logView, 12, 1, false) // Add logView to grid
			}
			isVisible = !isVisible // Toggle isVisible flag
			return event
		}
		// Check if Tab is pressed
		if event.Key() == tcell.KeyTab {
			// If Shift modifier is present
			if event.Modifiers() == tcell.ModShift {
				// Set focus on Fetch JSON button
				item := detailsForm.GetFormItem(1)
				app.SetFocus(item)
			} else {
				// Set focus on URL input
				item := detailsForm.GetFormItem(0)
				app.SetFocus(item)
			}
			return nil // Event fully handled, no further processing needed
		}

		return event // Return the unchanged event so it can continue being processed
	})
}

// UrlInputCapture handels input captures (keypresses) on the UrlInputField
func UrlInputCapture(
	logView *tview.TextView, // The logView to log events
	urlForm *tview.Form, // The form which contains the URL field
	detailsForm *tview.Form, // The form which contains detail fields
	textView *tui.ScrollTextView, // The textView to display logs
	detailsView *tview.TextView, // The detailsView to display response details
	methodbox *tview.Form, // The methodbox to select API request method
) {
	urlField := urlForm.GetFormItem(0).(*tview.InputField)                 // Get the URL input field from urlForm
	urlField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey { // Sets a handler function that gets executed whenever a key press occurs.
		if event.Key() == tcell.KeyEnter { // When Enter key is pressed...
			tui.SendAction(urlForm, detailsForm, logView, textView, detailsView, methodbox) // Send the HTTP request using the entered data
			_, method := methodbox.GetFormItem(0).(*tview.DropDown).GetCurrentOption()      // Get the current selected option from the drop-down in methodbox
			tui.LogMessage(logView, fmt.Sprintf("Main Using method: %s", method))           // Log the method used
		}
		return event // Return the unchanged event so it can continue being processed
	})
}
