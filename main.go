package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	tcell "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/valyala/fasthttp"
)

func main() {
	app := tview.NewApplication()

	// Initialize the textView.
	textView := tview.NewTextView()
	textView.SetDynamicColors(true)
	textView.SetBorder(true)
	textView.SetTitle("JSON Viewer")

	// Initialize the logView.

	logView := tview.NewTextView()
	logView.SetDynamicColors(true)
	logView.SetScrollable(true)
	logView.SetBorder(true)
	logView.SetTitle("Logs")

	detailsView := tview.NewTextView()
	detailsView.SetDynamicColors(true)
	detailsView.SetScrollable(true)
	detailsView.SetBorder(true)
	detailsView.SetTitle("Details")

	var form *tview.Form

	app.EnableMouse(true)

	// Initialize the form after logView, so we can use logView in the button's function.
	form = tview.NewForm().
		AddInputField("URL", "https://jsonplaceholder.typicode.com/posts", 100, nil, nil).
		AddButton("Fetch JSON", func() {
			url := form.GetFormItem(0).(*tview.InputField).GetText()
			statusCode, body, err := fasthttp.Get(nil, url)
			if err != nil {
				logMessage(logView, fmt.Sprintf("Error fetching URL: %v", err))
				return
			}
			if statusCode != fasthttp.StatusOK {
				textView.SetText(fmt.Sprintf("API responded with status code: %d", statusCode))
				return
			}

			var jsonData interface{}
			err = json.Unmarshal(body, &jsonData)
			if err != nil {
				textView.SetText(fmt.Sprintf("Error parsing JSON: %v", err))
				return
			}
			structure := visualizeJSONStructure(jsonData, "")
			textView.SetText(structure)
		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	textView.SetMouseCapture(
		func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
			// Check if it's a left button click
			if buttons := event.Buttons(); buttons&tcell.Button1 != 0 {
				x, y := event.Position()
				logMessage(logView, fmt.Sprintf("Mouse clicked in textView at x: %d, y: %d", x, y))
			}
			return action, event
		},
	)

	// Input capture for textView
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Log the captured key to logView
		logMessage(
			logView,
			fmt.Sprintf(
				"Key captured in textView: %v, Modifiers: %v",
				event.Key(),
				event.Modifiers(),
			),
		)

		// Check if Tab is pressed

		if event.Key() == tcell.KeyTab {
			// If Shift modifier is present
			if event.Modifiers() == tcell.ModShift {
				// Set focus on Fetch JSON button
				item := form.GetFormItem(1)
				app.SetFocus(item)
			} else {
				// Set focus on URL input
				item := form.GetFormItem(0)
				app.SetFocus(item)
			}
			return nil
		}
		return event
	})

	logMessage(logView, "Log window initialized!")

	grid := tview.NewGrid().
		SetRows(3, 0, 6). // No change here; same rows
		SetColumns(0, 0). // Split middle pane into two equal columns
		SetBorders(true).
		AddItem(form, 0, 0, 1, 2, 0, 0, true). // form spans both columns
		AddItem(detailsView, 1, 0, 1, 1, 0, 0, false).
		// detailsView is now in the left column of the middle row
		AddItem(textView, 1, 1, 1, 1, 0, 0, false).
		// textView is now in the right column of the middle row
		AddItem(logView, 2, 0, 1, 2, 0, 0, false) // logView spans both columns

	if err := app.SetRoot(grid, true).Run(); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}

func visualizeJSONStructure(data interface{}, indent string) string {
	switch v := data.(type) {
	case map[string]interface{}:
		b := &bytes.Buffer{}
		fmt.Fprintf(b, "{\n")
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		for _, key := range keys {
			valueStr := visualizeJSONStructure(v[key], indent+"  ")
			fmt.Fprintf(b, "%s  [blue]%s[white]: %s,\n", indent, key, strings.Trim(valueStr, "\n"))
		}
		fmt.Fprintf(b, "%s}", indent)
		return b.String()
	case []interface{}:
		b := &bytes.Buffer{}
		fmt.Fprintf(b, "[\n")
		for _, item := range v {
			itemStr := visualizeJSONStructure(item, indent+"  ")
			fmt.Fprintf(b, "%s  %s,\n", indent, strings.Trim(itemStr, "\n"))
		}
		fmt.Fprintf(b, "%s]", indent)
		return b.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func logMessage(logView *tview.TextView, msg string) {
	fmt.Fprintf(logView, "[::b]%s[::-]: %s\n", getTimeStamp(), msg)
	logView.ScrollToEnd() // Automatically scroll to the latest log
}

func getTimeStamp() string {
	return time.Now().Format("15:04:05")
}
