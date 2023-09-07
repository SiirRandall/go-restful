// The package tui includes functions to display and handle text based User Interface tasks
package tui

import (
	"strings" // Standard library used for string operations

	"github.com/gdamore/tcell/v2" // External library used for creating terminal applications
	"github.com/rivo/tview"       // External library used for managing terminal-based user interfaces
)

// ScrollTextView extends TextView from rivo/tview package with custom scrollbar.
type ScrollTextView struct {
	*tview.TextView
}

// NewScrollTextView function creates a new instance of ScrollTextView, which is scrollable.
func NewScrollTextView() *ScrollTextView {
	tv := tview.NewTextView()
	// Make Text View scrollable
	tv.SetScrollable(true)
	return &ScrollTextView{tv}
}

// Draw function draws the scrollable text view with a custom scrollbar on the screen
func (stv *ScrollTextView) Draw(screen tcell.Screen) {
	stv.TextView.Draw(screen)

	// Get inner dimensions of the text view
	x, y, width, height := stv.GetInnerRect()

	// Count total number of lines in text view
	totalRows := strings.Count(stv.GetText(true), "\n")

	if totalRows > height {
		// Calculate scrolling position and percentage scrolled
		scrollPosition, _ := stv.GetScrollOffset()
		percentageScrolled := float64(scrollPosition) / float64(totalRows-height+1)

		// Calculate scrollbar properties including its height and Y position on screen
		scrollbarHeight := max(
			1,
			int(float64(height-2)*(float64(height-2)/float64(totalRows))),
		)
		scrollbarY := 1 + int(
			percentageScrolled*float64(height-2-scrollbarHeight),
		)

		// Set X position of scrollbar
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
