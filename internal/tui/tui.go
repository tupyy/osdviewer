package tui

import (
	"context"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tupyy/osdviewer/internal/entity"
	"github.com/tupyy/osdviewer/internal/service"
)

const (
	keyOne = rune('1')
)

var (
	enviroments = map[int]service.Environment{1: service.Integration, 2: service.Stage, 3: service.Production}
)

type result[T any, E error] struct {
	Result T
	Err    E
}

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
	done chan chan interface{}
}

func New(app *tview.Application, fmReader service.FleetManagerReader) *Tui {
	t := Tui{
		app:      app,
		pages:    tview.NewPages(),
		fmReader: fmReader,
		views:    make(map[string]*View),
		done:     make(chan chan interface{}),
	}

	for i := 1; i < 4; i++ {
		v := t.addPage(enviroments[i].String())
		t.views[enviroments[i].String()] = v
	}

	t.pages.AddPage("help", newHelpView(), true, true)

	return &t
}

// Start starts a go routin which draws app every 0.5s.
// In this way, we avoid to pass app pointer to every primitive which needs to be redrawn
func (t *Tui) Start() {
	// read clusters from all env
	for _, e := range enviroments {
		result := <-t.getClusters(context.TODO(), e)
		if result.Err == nil {
			if view, ok := t.views[e.String()]; ok {
				view.SetData(result.Result)
			}
		}
	}

	go func(done chan chan interface{}) {
		for {
			select {
			case <-time.After(500 * time.Millisecond):
				t.app.Draw()
			case d := <-done:
				d <- struct{}{}
				return
			}
		}
	}(t.done)

	go func(done chan chan interface{}) {
		for {
			select {
			case <-time.After(1 * time.Second):
				page := t.currentPage()
				// get the current page
				var e service.Environment
				switch page {
				case "integration":
					e = service.Integration
				case "stage":
					e = service.Stage
				case "production":
					e = service.Production
				default:
					break
				}

				// wait until ocm reads the clusters
				result := <-t.getClusters(context.TODO(), e)
				if result.Err == nil {
					if view, ok := t.views[page]; ok {
						view.SetData(result.Result)
					}
				}
			case d := <-done:
				d <- struct{}{}
				return
			}
		}
	}(t.done)
}

func (t *Tui) Stop() {
	d := make(chan interface{})
	t.done <- d
	<-d
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
			t.showPage(enviroments[idx].String())
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

	t.showPage(enviroments[t.currentPageIdx].String())
}

// Show the previous page. If the current page is the first one than show the last page.
func (t *Tui) previousPage() {
	t.currentPageIdx -= 1
	if t.currentPageIdx < 1 {
		t.currentPageIdx = len(t.views)
	}

	t.showPage(enviroments[t.currentPageIdx].String())
}

func (t *Tui) showPage(name string) {
	if t.navBar == nil {
		t.navBar = t.createNavBar()
		t.rootFlex.AddItem(t.navBar, 1, 1, true)
	}

	view := t.views[name]
	view.Clear()
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
		navBar.AddPage(i, enviroments[i].String())
	}
	return navBar
}

func (t *Tui) currentPage() string {
	currentPageName, _ := t.pages.GetFrontPage()
	return currentPageName
}

func (t *Tui) getClusters(ctx context.Context, e service.Environment) chan result[[]entity.Cluster, error] {
	resultCh := make(chan result[[]entity.Cluster, error])

	go func() {
		clusters, err := t.fmReader.GetClusters(ctx, e)
		resultCh <- result[[]entity.Cluster, error]{clusters, err}
	}()

	return resultCh
}
