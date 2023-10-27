package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tupyy/osdviewer/internal/entity"
)

type state int

const (
	LoadingState state = iota
	ReadyState
	ErrorState
)

type ViewState struct {
	State state
	Err   error
}

type ClusterView struct {
	*tview.Box

	app         *tview.Application
	state       ViewState
	table       *tview.Table
	textView    *tview.TextView
	searchField *tview.InputField
	flex        *tview.Flex
	sourceData  *ClusterTableContent
}

func NewClusterView(name string) *ClusterView {
	v := &ClusterView{
		Box:         tview.NewBox().SetBackgroundColor(tcell.ColorBlack),
		table:       tview.NewTable().SetBorders(false),
		textView:    tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Reading from OCM..."),
		searchField: tview.NewInputField().SetLabel("Search").SetLabelWidth(len("Search") + 1).SetFieldWidth(0).SetFieldBackgroundColor(tcell.ColorRed),
		flex:        tview.NewFlex().SetDirection(tview.FlexRow),
		state:       ViewState{State: LoadingState},
	}

	v.table.SetBorderPadding(1, 0, 2, 2)
	v.searchField.SetBorderPadding(1, 0, 2, 2)
	v.searchField.SetFieldBackgroundColor(tcell.ColorGrey)
	v.searchField.SetAcceptanceFunc(func(textToCheck string, lastChar rune) bool { return true })
	v.flex.AddItem(v.searchField, 3, 0, true).AddItem(v.table, 0, 3, false)

	return v
}

func (v *ClusterView) Draw(screen tcell.Screen) {
	v.Box.Draw(screen)
	v.Box.SetBorder(false)

	x, y, width, height := v.GetInnerRect()
	switch v.state.State {
	case LoadingState:
		v.textView.SetRect(x, height/2, width, height-1)
		v.textView.Draw(screen)
	case ErrorState:
		v.textView.SetText(v.state.Err.Error())
		v.textView.SetRect(x+5, height/2, width-5, height-1)
		v.textView.Draw(screen)
	default:
		v.searchField.SetFieldWidth(width - len(v.searchField.GetLabel()) - 2)
		v.searchField.Draw(screen)
		v.flex.SetRect(x, y, width, height)
		v.flex.Draw(screen)
	}
}

func (v *ClusterView) Model(model any) {
	c, ok := model.(result[[]entity.Cluster, error])
	if ok {
		if c.Err == nil {
			v.sourceData = NewTableContent(c.Result)
			v.sourceData.SetFilter(v.getFilterFunc())
			v.table.SetContent(v.sourceData)
			v.state = ViewState{State: ReadyState}
		} else {
			if v.sourceData != nil {
				v.sourceData = nil
			}
			v.state = ViewState{State: ErrorState, Err: c.Err}
		}
	}
}

// Focus set the focus either on the menu is showMenu is true or on the textView.
func (v *ClusterView) Focus(delegate func(p tview.Primitive)) {
	delegate(v.searchField)
}

func (v *ClusterView) SetFocus() {
	v.app.SetFocus(v.flex)
}

func (v *ClusterView) HandleEventKey(key *tcell.EventKey) {
	if v.searchField == nil {
		return
	}
	i := v.searchField.InputHandler()
	i(key, func(p tview.Primitive) {})
}

func (v *ClusterView) getFilterFunc() func(entity.Cluster) bool {
	return func(c entity.Cluster) bool {
		if v.searchField.GetText() == "" {
			return true
		}
		return strings.Index(c.String(), v.searchField.GetText()) > 0
	}
}
