package tui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/tupyy/osdviewer/internal/entity"
)

type View struct {
	app *tview.Application

	table *tview.Table
}

func NewView(name string, app *tview.Application) *View {
	m := &View{
		app:   app,
		table: tview.NewTable().SetBorders(false),
	}

	return m
}

func (m *View) Layout() tview.Primitive {
	return m.table
}

func (m *View) ShowMenu() {
}

func (m *View) SetData(clusters []entity.Cluster) {
	m.table.SetContent(NewTableContent(clusters))
}

func (m *View) HandleEventKey(key *tcell.EventKey) {
	switch key.Rune() {
	case rune('m'):
		m.ShowMenu()
	}
}
