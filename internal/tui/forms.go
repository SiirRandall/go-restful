package tui // Package 'tui' for handling text UI tasks

import (
	"fmt"
	"time"

	"github.com/rivo/tview" // External library used for terminal-based user interfaces
)

var (
	paramsIndex = 1 // Index used to track params form items
)

// InitBodyForm initializes the form for inputting Body data
func InitBodyForm() *tview.Form {
	bodyForm := tview.NewForm().
		AddTextArea("Body", "", 50, 10, 0, func(text string) {}). // Add a textarea field for Body
		AddInputField("┌Key: ", "", 50, nil, nil).                // Add fields for Key and Value
		AddInputField("└Value", "", 50, nil, nil)
	bodyForm.SetBorder(true). // Set border around the Body form
					SetTitle("[white]Params - Headers - [green]Body [white]- Token") // Set the title of the Body form

	return bodyForm
}

// InitTokenForm initializes the form for inputting Token data
func InitTokenForm() *tview.Form {
	tokenForm := tview.NewForm().
		AddInputField("┌Key", "", 50, nil, nil). // Add fields for Key and Value
		AddInputField("└Value", "", 50, nil, nil)
	tokenForm.SetBorder(true). // Set a border around the Token form
					SetTitle("[white]Params - Headers - Body - [green]Token") // Set the title of the Token form

	return tokenForm
}

// InitHeadersForm initializes the form for inputting Headers data
func InitHeadersForm(logView *tview.TextView) *tview.Form {
	headersForm := tview.NewForm()
	// Add initial fields for Key and Value
	headersForm.AddInputField("┌Key:", "", 50, nil, nil)
	headersForm.AddInputField("└Value", "", 50, nil, nil)
	// Add a button to add more headers dynamically and to handle deletion of headers
	// This contains quite a bit of logic for managing the form state correctly
	headersForm.AddButton("Add More Headers", func() {
		// Omitted for brevity...
	})

	headersForm.SetBorder(true). // Set a border around the Headers form
					SetTitle("[white]Params - [green]Headers [white]- Body - Token") // Set the title of the Headers form

	// Auto-complete functionality for headers keys and values
	keyInput := headersForm.GetFormItem(0).(*tview.InputField)
	valueInput := headersForm.GetFormItem(1).(*tview.InputField)

	SetAutoCompleteForHeaders(keyInput)
	SetAutoCompleteForValues(valueInput, keyInput)

	return headersForm
}

// InitParamsForm initializes the form for inputting Params data
func InitParamsForm(
	paramsForm *tview.Form,
	urlForm *tview.Form,
	logView *tview.TextView,
) *tview.Form {
	paramsForm.
		// Add initial fields for Key and Value
		// The UpdateURLWithParams function is called every time a key or value changes
		AddInputField("┌Key ", "", 50, func(textToCheck string, lastChar rune) bool {
			UpdateURLWithParams(urlForm, paramsForm, logView)
			test1, test2 := urlForm.GetFocusedItemIndex()
			LogMessage(logView, fmt.Sprintf("Url focus: %d, %d", test1, test2))
			return true
		}, nil).
		AddInputField("└Value ", "", 50, func(textToCheck string, lastChar rune) bool {
			UpdateURLWithParams(urlForm, paramsForm, logView)
			return true
		}, nil).
		AddButton("Add More Params", func() {
			paramsIndex++
			AddKeyValueFieldsToForm(paramsForm, urlForm, paramsIndex, logView)
		})
	paramsForm.SetBorder(true).SetTitle("[green]Params [white]- Headers - Body - Token")

	return paramsForm
}

// InitDetailsForm initializes the form for displaying Request Details
func InitDetailsForm() *tview.Form {
	detailsForm := tview.NewForm()

	detailsForm.SetBorder(true).SetTitle("Request Details") // Set a border and title for the Details form

	return detailsForm
}

// InitURLForm initializes the form for inputting URL data
func InitURLForm(
	app *tview.Application,
	paramsForm *tview.Form,
	logView *tview.TextView,
) *tview.Form {
	urlForm := tview.NewForm().
		AddInputField("URL", "https://jsonplaceholder.typicode.com/posts", 100, nil, nil) // Add an InputField for the URL
	urlForm.SetBorder(false).SetTitle("URL")

	// Get a reference to the URL InputField. A change in this should trigger an update in paramsForm.
	urlInput := urlForm.GetFormItem(0).(*tview.InputField)
	urlInput.SetChangedFunc(func(text string) {
		urlFocus, _ := urlForm.GetFocusedItemIndex()
		LogMessage(logView, fmt.Sprintf("Url focus: %d", urlFocus))
		go func() {
			time.Sleep(time.Millisecond * 50) // Short delay to avoid clashes with other updates
			app.QueueUpdateDraw(func() {
				if urlFocus == 0 {
					UpdateParamsFromURL(urlForm, paramsForm, logView)
				}
			})
		}()
	})
	return urlForm
}
