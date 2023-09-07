// Package tui holds components for terminal-based user interfaces
package tui

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/rivo/tview" // Terminal UI library

	httpclient "github.com/SiirRandall/go-restful/internal/httpClient" // HTTP client module
	"github.com/SiirRandall/go-restful/internal/text"                  // Text utilities module
	"github.com/SiirRandall/go-restful/internal/utils"                 // General utility functions module
)

var lastQueryLength int    // Global variable to keep the count of query parameters from the last request
const JSON_VIEW_WIDTH = 50 // Width of the JSON view panel

// UpdateURLWithParams function reads the parameters from a form and appends them to a URL in the form, logging any errors.
func UpdateURLWithParams(urlForm, paramsForm *tview.Form, logView *tview.TextView) {
	// Get the base URL from the form input field
	baseURL := urlForm.GetFormItem(0).(*tview.InputField).GetText()

	// Parse the base URL into a url.URL object
	u, err := url.Parse(baseURL)
	// If an error occurred while parsing, log the error message and terminate the function
	if err != nil {
		LogMessage(logView, fmt.Sprintf("Error parsing URL: %v", err))
		return
	}

	// Get the current query values (if any) from the parsed URL
	queryValues := u.Query()

	// Log the current query parameters for debugging purposes
	LogMessage(logView, fmt.Sprintf("Query values: %v", queryValues))

	// Iterate over all form items (query parameters) in pairs of two (key-value)
	for i := 0; i < paramsForm.GetFormItemCount()-1; i += 2 {

		// Get both key and value field from the pair
		keyField := paramsForm.GetFormItem(i).(*tview.InputField)
		valueField := paramsForm.GetFormItem(i + 1).(*tview.InputField)

		// Read text from the fields
		key := keyField.GetText()
		value := valueField.GetText()

		// Log the extracted key-value pair for debugging purposes
		LogMessage(logView, fmt.Sprintf("Key: %s, Value: %s", key, value))

		// If either the key or the value is empty, skip this iteration
		if key == "" || value == "" {
			continue
		}

		// Otherwise, set the key and value as a query parameter in the URL
		queryValues.Set(key, value)
	}

	// After iterating all parameters, encode them into the URL's RawQuery field
	u.RawQuery = queryValues.Encode()

	// Finally, update the form's URL field with the new URL which includes all query parameters
	urlForm.GetFormItem(0).(*tview.InputField).SetText(u.String())
}

// The UpdateParamsFromURL function decodes parameters from a URL and populates a form with the parameter keys and values.
func UpdateParamsFromURL(urlForm, paramsForm *tview.Form, logView *tview.TextView) {
	// Get the URL from the urlForm input field
	rawURL := urlForm.GetFormItem(0).(*tview.InputField).GetText()

	// Parse the obtained URL
	u, err := url.Parse(rawURL)
	// If an error occurred while parsing, log the error message and terminate the function
	if err != nil {
		LogMessage(logView, fmt.Sprintf("Error parsing URL: %v", err))
		return
	}

	// Extract the query parameters from the parsed URL
	queryValues := u.Query()

	// Get the length of the current query - number of parameters
	currentQueryLength := len(queryValues)

	// Initialize the currentIndex which will be used to index the form items
	currentIndex := 1

	for key, values := range queryValues {
		// Get value for the respective key. Assuming only one value per key
		value := values[0]

		// If the currentIndex is within the last parameters length or if it's the first index and there were no previous params
		if currentIndex <= lastQueryLength || (currentIndex == 1 && lastQueryLength == 0) {

			// Get form fields for key and value
			keyField := paramsForm.GetFormItem((currentIndex - 1) * 2).(*tview.InputField)
			valueField := paramsForm.GetFormItem((currentIndex-1)*2 + 1).(*tview.InputField)

			// Set labels for input fields
			keyField.SetLabel(fmt.Sprintf("┌[%d]Key", currentIndex))
			valueField.SetLabel(fmt.Sprintf("└[%d]Value", currentIndex))

			// Fill in key and value fields with extracted key and value
			keyField.SetText(key)
			valueField.SetText(value)
		} else if currentQueryLength > 1 {

			// If the currentIndex goes beyond the last parameters length, add new key-value pair form fields
			paramsForm.AddInputField(
				fmt.Sprintf("┌[%d]Key", currentIndex),
				key,
				50,
				func(textToCheck string, lastChar rune) bool {
					// On change of these fields, update the URL with new parameters
					UpdateURLWithParams(urlForm, paramsForm, logView)
					return true
				},
				nil,
			)
			paramsForm.AddInputField(fmt.Sprintf("└[%d]Value", currentIndex), value, 50, func(textToCheck string, lastChar rune) bool {
				UpdateURLWithParams(urlForm, paramsForm, logView)
				return true
			}, nil)
		}

		// Increment the currentIndex after completing processing current one
		currentIndex++
	}

	// If there are less parameters now than before, remove the additional empty form fields
	for currentIndex <= lastQueryLength {
		paramsForm.RemoveFormItem((currentIndex - 1) * 2) // Key
		paramsForm.RemoveFormItem((currentIndex - 1) * 2) // Value
		currentIndex++
	}

	lastQueryLength = currentQueryLength
}

// SendAction function sends an HTTP request based on the settings provided in the form fields.
func SendAction(
	urlForm *tview.Form,
	detailsForm *tview.Form,
	logView *tview.TextView,
	textView *ScrollTextView,
	detailsView *tview.TextView,
	methodbox *tview.Form,
) {
	// Fetch URL from the first form item of urlForm
	url := urlForm.GetFormItem(0).(*tview.InputField).GetText()

	// Get the selected HTTP method from the dropdown menu in methodbox form
	_, methodText := methodbox.GetFormItem(0).(*tview.DropDown).GetCurrentOption()

	// Get Headers and RequestBody from detailsForm
	headersInput := "" // need to pull from headers and body page
	requestBody := ""
	// Log selected HTTP method, headers and body to logView
	LogMessage(logView, fmt.Sprintf("Using method: %s", methodText))
	LogMessage(logView, fmt.Sprintf("Headers: %s", headersInput))
	LogMessage(logView, fmt.Sprintf("Body: %s", requestBody))

	// Show a brief summary of details in detailsView
	detailsView.SetText(
		fmt.Sprintf("Method: %s\nHeaders: %s\nBody: %s", methodText, headersInput, requestBody),
	)

	// Parse headers from each line to key-value pairs
	headerLines := strings.Split(headersInput, "\n")
	headers := make(map[string]string)

	for _, line := range headerLines {
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			headers[parts[0]] = parts[1]
		}
	}

	// Send the HTTP request with the populated details
	response := httpclient.SendHttpRequest(httpclient.HttpRequestDetails{
		URL:         url,
		Method:      methodText,
		Headers:     headers,
		RequestBody: requestBody,
	})

	// Check for errors in response. If error exists, log it and return
	if response.Error != nil {
		LogMessage(logView, response.Error.Error())
		return
	}

	// Process the received response
	if response.JsonData != nil {
		// When valid JSON data is received, visualize the JSON structure and set text of textView
		_, _, jsonwidth, _ := textView.GetRect()
		structure := text.VisualizeJSONStructure(response.JsonData, "", jsonwidth-6)
		textView.SetText(structure)

		// Set max lines of textView according to visualized json structure
		limitLines := utils.CountLines(structure)
		textView.SetMaxLines(limitLines + 1)
	} else {
		// If not JSON, then put received plain body into textView
		textView.SetText(string(response.Body))
	}
}
