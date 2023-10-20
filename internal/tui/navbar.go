package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

// NavBar displays the navigation bar at the bottom of the screen.
type NavBar struct {
	*tview.TextView
}

func NewNavBar() *NavBar {
	navBar := &NavBar{TextView: tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false),
	}
	return navBar
}

// AddPage shows the names of the pages in the navBar.
func (navBar *NavBar) AddPage(idx int, name string) {
	fmt.Fprintf(navBar, `%d ["%s"][yellow]%s[white][""]  `, idx, name, name)
}

// SelectPage highlight the page.
func (navBar *NavBar) SelectPage(name string) {
	navBar.Highlight(name).ScrollToHighlight()
}
