package main

import (
	"log"

	"github.com/rivo/tview" // Importing the tview package for terminal-based UI applications

	"github.com/SiirRandall/go-restful/internal/input" // Importing internal packages for handling inputs
	"github.com/SiirRandall/go-restful/internal/tui"   // Importing internal packages for text UI creation
)

func main() {
	// Initializing our app as a new tview Application.
	app := tview.NewApplication()

	// Enabling mouse support in the app.
	app.EnableMouse(true)

	// Initialize a form that will hold the parameters for HTTP requests.
	paramsForm := tview.NewForm()

	// Initialize the log view which displays all the logs in the application.
	logView := tui.InitLogView()

	// Initialize a form that displays details of requests and responses.
	detailsForm := tui.InitDetailsForm()

	// Initialize URL form which is used to get URL from user.
	urlForm := tui.InitURLForm(app, paramsForm, logView)

	// Initialize parameter form which is used to get request parameters from user.
	paramsForm = tui.InitParamsForm(paramsForm, urlForm, logView)

	// Initialize form used to get request headers from user.
	headersForm := tui.InitHeadersForm(logView)

	// Initialize form used to input the body of an HTTP request.
	bodyForm := tui.InitBodyForm()

	// Initialize form used to get token for authentication.
	tokenForm := tui.InitTokenForm()

	// Initialize the pages rendered on HTML.
	htmlPages := tui.InitHTMLPages(paramsForm, headersForm, bodyForm, tokenForm, logView)

	// Initialises the JSON viewer which shows pretty-printed JSON structures.
	textView := tui.InitJsonViewer()

	// Initialize a view to display request and response details.
	detailsView := tui.InitDetailsView()

	// Initialize components related to the URL input & buttons for different actions.
	methodbox, buttonPanel := tui.InitUrlComponents(
		app,
		urlForm,
		detailsForm,
		logView,
		textView,
		detailsView,
	)

	// Integrates and initializes Url input field and associated action buttons.
	urlAndButtons := tui.InitUrlandButtons(methodbox, urlForm, buttonPanel)

	// Initializes the main grid layout with the elements for displaying http request, response and other details.
	grid := tui.InitGrid(urlAndButtons, htmlPages, textView, detailsForm)

	// Captures mouse interaction within the TextView component.
	input.TextViewMouseCapture(logView, textView, headersForm)

	// Captures keyboard and mouse input for the URL form.
	input.UrlInputCapture(logView, urlForm, detailsForm, textView, detailsView, methodbox)

	// Captures keyboard and mouse input for the TextView.
	input.TextViewKBCapture(app, textView, logView, detailsForm, grid)

	// Run the application - setting the root element and make it full screen. Exits if there are errors.
	if err := app.SetRoot(grid, true).Run(); err != nil {
		log.Fatalf(
			"Failed to run application: %v",
			err,
		) // Logs the fatal error and quits the app smoothly.
	}
}
