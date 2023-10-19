package tui

import (
	"context"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tupyy/osdviewer/internal/service"
)

const (
	keyOne = rune('1')
)

var (
	enviroments = map[int]string{1: "integration", 2: "stage", 3: "prod"}
)

type Tui struct {
	app *tview.Application

	// Root flex which contains all other primitives
	rootFlex *tview.Flex

	// Navbar shows page names
	navBar *NavBar

	// Pages
	pages *tview.Pages

	views map[string]*View

	currentPageIdx int

	fmReader service.FleetManagerReader

	// close the start method
	done chan interface{}
}

func New(app *tview.Application, fmReader service.FleetManagerReader) *Tui {
	t := Tui{
		app:      app,
		pages:    tview.NewPages(),
		fmReader: fmReader,
		views:    make(map[string]*View),
		done:     make(chan interface{}),
	}

	for i := 1; i < 4; i++ {
		v := t.addPage(enviroments[i])
		t.views[enviroments[i]] = v
	}

	t.pages.AddPage("help", newHelpView(), true, true)

	return &t
}

// Start starts a go routin which draws app every 0.5s.
// In this way, we avoid to pass app pointer to every primitive which needs to be redrawn
func (t *Tui) Start() {
	go func(done chan interface{}) {
		for {
			select {
			case <-time.After(500 * time.Millisecond):
				t.app.Draw()
			case <-done:
				return
			}
		}
	}(t.done)

	go func(done chan interface{}) {
		for {
			select {
			case <-time.After(1 * time.Second):
				// get the current page
				var e service.Environment
				switch t.currentPage() {
				case "integration":
					e = service.Integration
				case "stage":
					e = service.Stage
				case "prod":
					e = service.Production
				default:
					break
				}

				clusters, err := t.fmReader.GetClusters(context.TODO(), e)
				if err == nil {
					t.views[t.currentPage()].SetData(clusters)
				}
			case <-done:
				return
			}
		}
	}(t.done)
}

func (t *Tui) Stop() {
	t.done <- struct{}{}
}

func (t *Tui) HandleEventKey(key *tcell.EventKey) {
	switch key.Key() {
	case tcell.KeyLeft:
		t.previousPage()
	case tcell.KeyEnter:
		if t.currentPage() == "help" {
			t.nextPage()
		}
	case tcell.KeyRight:
		t.nextPage()
	default:
		// if the key is a page number then show the page
		idx := int(key.Rune() - keyOne)
		if idx < len(t.views) && idx >= 0 {
			t.showPage(enviroments[idx])
		}
	}
}

// Layout returns the root flex
func (t *Tui) Layout() tview.Primitive {
	t.rootFlex = tview.NewFlex().SetDirection(tview.FlexRow).AddItem(t.pages, 0, 1, true)
	return t.rootFlex
}

// Show the next page. If the current page is the last page than show the first page.
func (t *Tui) nextPage() {
	t.currentPageIdx += 1
	if t.currentPageIdx > len(t.views) {
		t.currentPageIdx = 1
	}

	t.showPage(enviroments[t.currentPageIdx])
}

// Show the previous page. If the current page is the first one than show the last page.
func (t *Tui) previousPage() {
	t.currentPageIdx -= 1
	if t.currentPageIdx < 1 {
		t.currentPageIdx = len(t.views)
	}

	t.showPage(enviroments[t.currentPageIdx])
}

func (t *Tui) showPage(name string) {
	if t.navBar == nil {
		t.navBar = t.createNavBar()
		t.rootFlex.AddItem(t.navBar, 1, 1, true)
	}

	t.pages.SwitchToPage(name)
	t.navBar.SelectPage(name)
}

func (t *Tui) addPage(name string) *View {
	v := NewView(name, t.app)
	t.pages.AddPage(name, v.Layout(), true, true)
	return v
}

func (t *Tui) createNavBar() *NavBar {
	navBar := NewNavBar()
	for i := 1; i < 4; i++ {
		navBar.AddPage(i, enviroments[i])
	}
	return navBar
}

func (t *Tui) currentPage() string {
	currentPageName, _ := t.pages.GetFrontPage()
	return currentPageName
}
