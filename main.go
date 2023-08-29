package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	tcell "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/valyala/fasthttp"
)

type ScrollTextView struct {
	*tview.TextView
}

func NewScrollTextView() *ScrollTextView {
	tv := tview.NewTextView()
	tv.SetScrollable(true)
	return &ScrollTextView{tv}
}

func (stv *ScrollTextView) Draw(screen tcell.Screen) {
	stv.TextView.Draw(screen)

	x, y, width, height := stv.GetInnerRect()
	totalRows := strings.Count(stv.GetText(true), "\n")

	if totalRows > height {
		scrollPosition, _ := stv.GetScrollOffset()
		percentageScrolled := float64(scrollPosition) / float64(totalRows-height+1)
		scrollbarHeight := max(
			1,
			int(float64(height-2)*(float64(height-2)/float64(totalRows))),
		) // -2 to account for arrows
		scrollbarY := 1 + int(
			percentageScrolled*float64(height-2-scrollbarHeight),
		) // +1 to start below the up arrow

		// Scrollbar position
		scrollbarX := x + width - 1

		// Draw the scrollbar background in mid-gray
		for i := 1; i < height-1; i++ {
			screen.SetContent(
				scrollbarX,
				y+i,
				'▒',
				nil,
				tcell.StyleDefault.Foreground(tcell.ColorGray).Background(tcell.ColorBlack),
			)
		}

		// Draw the scrollbar thumb in bright white
		for i := 0; i < scrollbarHeight; i++ {
			screen.SetContent(
				scrollbarX,
				y+scrollbarY+i,
				'█',
				nil,
				tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack),
			)
		}

		// Draw simple triangle arrows in bright white
		screen.SetContent(
			scrollbarX,
			y,
			'▲',
			nil,
			tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack),
		)
		screen.SetContent(
			scrollbarX,
			y+height-1,
			'▼',
			nil,
			tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack),
		)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ... rest of your code ...

func main() {
	app := tview.NewApplication()

	// Initialize the textView.
	// textView := tview.NewTextView()
	// textView.SetDynamicColors(true)
	// textView.SetBorder(true)
	// textView.SetTitle("JSON Viewer")
	// textView.SetScrollable(true)
	// textView.SetBorderPadding(0, 0, 1, 1)
	// textView.SetScrollBar(true)
	// Initialize the logView.

	textView := NewScrollTextView()
	textView.SetDynamicColors(true)
	textView.SetBorder(true)
	textView.SetTitle("JSON Viewer")

	logView := tview.NewTextView()
	logView.SetDynamicColors(true)
	logView.SetScrollable(true)
	logView.SetBorder(true)
	logView.SetTitle("Logs")

	detailsView := tview.NewTextView()
	detailsView.SetDynamicColors(true)
	detailsView.SetScrollable(true)
	detailsView.SetBorder(false)
	detailsView.SetTitle("Details")

	urlForm := tview.NewForm().
		AddInputField("URL", "https://jsonplaceholder.typicode.com/posts", 100, nil, nil)
	urlForm.SetBorder(false).SetTitle("URL")

	var form *tview.Form
	var detailsForm *tview.Form
	isVisible := false
	//	logVisible := false

	app.EnableMouse(true)

	buttonPanel := tview.NewForm().
		AddButton("Send", func() {
			sendAction(urlForm, detailsForm, logView, textView, detailsView)
		}).
		AddButton("Quit", func() {
			app.Stop()
		})

	methodbox := tview.NewForm().
		AddDropDown("", []string{"GET", "POST", "PUT", "DELETE"}, 0, nil)

	urlAndButtons := tview.NewFlex().
		AddItem(methodbox, 9, 1, false).
		AddItem(urlForm, 0, 4, true).
		AddItem(buttonPanel, 0, 1, false)

	detailsForm = tview.NewForm().
		AddDropDown("Method", []string{"GET", "POST", "PUT", "DELETE"}, 0, nil).
		AddInputField("Headers", "", 50, nil, nil).
		AddInputField("Body", "", 50, nil, nil)
	detailsForm.SetBorder(true).SetTitle("Request Details")
	// Initialize the form after logView, so we can use logView in the button's function.
	form = tview.NewForm()
	form.SetBorder(true).SetTitle("Details")
	// /	AddInputField("URL", "https://jsonplaceholder.typicode.com/posts", 100, nil, nil).

	// Create two new input fields for the `Headers` and `Body`.

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

	urlField := urlForm.GetFormItem(0).(*tview.InputField)
	urlField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			sendAction(urlForm, detailsForm, logView, textView, detailsView)
		}
		return event
	})

	// Input capture for textView

	logMessage(logView, "Log window initialized!")

	// grid := tview.NewGrid().
	// 	SetRows(3, 0, 6).
	// 	// 4 units for top row, dynamic size for middle row, 6 units for bottom
	// 	SetColumns(35, 0, 30). // Adjust sizes based on your preference
	// 	SetBorders(true).
	// 	AddItem(urlForm, 0, 0, 1, 2, 0, 0, true).      // URL form spanning 2 columns
	// 	AddItem(form, 0, 2, 1, 1, 0, 0, false).        // Main form (buttons) on the right
	// 	AddItem(detailsForm, 1, 0, 1, 1, 0, 0, false). // Details form in the left
	// 	AddItem(textView, 1, 1, 1, 2, 0, 0, false)     // Spanning 2 columns for the textView
	// 	//		AddItem(logView, 2, 0, 1, 3, 0, 0, false)      // Spanning 3 columns for the logViewi
	grid := tview.NewFlex().
		SetDirection(tview.FlexRow).
		//		AddItem(urlForm, 3, 1, true).
		AddItem(urlAndButtons, 3, 1, true).
		AddItem(tview.NewFlex().
			AddItem(detailsForm, 35, 1, false).
			AddItem(textView, 0, 1, false).
			AddItem(form, 30, 1, false),
			0, 1, false,
		)
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

		if event.Rune() == 'l' {
			if isVisible {
				grid.RemoveItem(logView) // Remove logView from grid
			} else {
				grid.AddItem(logView, 6, 1, false) // Add logView to grid
			}
			isVisible = !isVisible
			app.Draw() // Refresh the UI
			return event
		}
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

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		//		log.Println("Key pressed:", event.Rune())
		logMessage(logView, "Toggling logView visibility")

		return event
	})
	if err := app.SetRoot(grid, true).Run(); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}
}

func visualizeJSONStructure(data interface{}, indent string) string {
	switch v := data.(type) {
	case map[string]interface{}:
		b := &bytes.Buffer{}
		fmt.Fprintf(b, "{\n")

		// Extract and sort the keys for consistent output.
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys) // Sorting the keys

		for i, key := range keys {
			valueStr := visualizeJSONStructure(v[key], indent+"  ")
			comma := ","
			if i == len(keys)-1 {
				comma = "" // No comma for the last key
			}
			fmt.Fprintf(b, "%s  [blue]%s[white]: %s%s\n", indent, key, strings.Trim(valueStr, "\n"), comma)
		}

		fmt.Fprintf(b, "%s}", indent)
		return b.String()

	case []interface{}:
		b := &bytes.Buffer{}
		fmt.Fprintf(b, "[\n")

		for i, item := range v {
			itemStr := visualizeJSONStructure(item, indent+"  ")
			comma := ","
			if i == len(v)-1 {
				comma = "" // No comma for the last item
			}
			fmt.Fprintf(b, "%s  %s%s\n", indent, strings.Trim(itemStr, "\n"), comma)
		}

		fmt.Fprintf(b, "%s]", indent)
		return b.String()

	default:
		// Apply custom styles based on the type of the value
		switch t := data.(type) {
		case string:
			return fmt.Sprintf("[orange]\"%s\"[white]", t)
		case float64:
			if t == float64(int(t)) {
				return fmt.Sprintf("[red]%d[white]", int(t))
			}
			return fmt.Sprintf("[green]%f[white]", t)
		case int, int32, int64:
			return fmt.Sprintf("[red]%d[white]", t)
		default:
			return fmt.Sprintf("%v", t)
		}
	}
}

func logMessage(logView *tview.TextView, msg string) {
	fmt.Fprintf(logView, "[::b]%s[::-]: %s\n", getTimeStamp(), msg)
	logView.ScrollToEnd() // Automatically scroll to the latest log
}

func getTimeStamp() string {
	return time.Now().Format("15:04:05")
}

func sendAction(
	urlForm *tview.Form,
	detailsForm *tview.Form,
	logView *tview.TextView,
	textView *ScrollTextView,
	detailsView *tview.TextView,
) {
	url := urlForm.GetFormItem(0).(*tview.InputField).GetText()
	_, methodText := detailsForm.GetFormItem(0).(*tview.DropDown).GetCurrentOption()
	headers := detailsForm.GetFormItem(1).(*tview.InputField).GetText()
	requestBody := detailsForm.GetFormItem(2).(*tview.InputField).GetText()

	logMessage(logView, fmt.Sprintf("Using method: %s", methodText))
	logMessage(logView, fmt.Sprintf("Headers: %s", headers))
	logMessage(logView, fmt.Sprintf("Body: %s", requestBody))

	detailsView.SetText(
		fmt.Sprintf("Method: %s\nHeaders: %s\nBody: %s", methodText, headers, requestBody),
	)

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod(methodText)
	// For the sake of simplicity, I won't set headers here. You can parse the headers string and set them using req.Header.Set().
	req.SetBodyString(requestBody)

	headersInput := detailsForm.GetFormItem(1).(*tview.InputField).GetText()
	headerLines := strings.Split(headersInput, "\n")

	for _, line := range headerLines {
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			headerName := parts[0]
			headerValue := parts[1]
			req.Header.Set(headerName, headerValue)
		}
	}

	err := fasthttp.Do(req, resp)
	if err != nil {
		logMessage(
			logView,
			fmt.Sprintf("Error making %s request to URL: %v", methodText, err),
		)
		return
	}

	statusCode, body, err := fasthttp.Get(nil, url)
	if err != nil {
		logMessage(logView, fmt.Sprintf("Error fetching URL: %v", err))
		return
	}
	textView.SetText(string(body)) // Convert the byte slice to a string

	rawJSONResponse := string(body)
	textView.SetText(rawJSONResponse)

	if statusCode != fasthttp.StatusOK && statusCode != fasthttp.StatusCreated {
		// StatusCreated (201) is also a success status for POST requests
		return
	}
	// var jsonData map[string]interface{} // Using a map to get a structured response
	// err = json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	textView.SetText(fmt.Sprintf("Error parsing JSON: %v", err))
	// 	return
	// }
	body = resp.Body()
	var jsonData1 interface{}
	err = json.Unmarshal(body, &jsonData1)
	if err != nil {
		textView.SetText(fmt.Sprintf("Error parsing JSON: %v", err))
		textView.SetText(string(body))
		return
	}

	textView.SetText(string(body))

	//			prettyJSON, _ := json.MarshalIndent(jsonData1, "", "    ")

	structure := visualizeJSONStructure(jsonData1, "")
	textView.SetText(structure)
	// rawBody := resp.Body()                                            // Get the body
	//	logMessage(logView, fmt.Sprintf("Raw Body: %s", string(rawBody))) // Log the body
}
