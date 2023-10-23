package tui

import (
	"github.com/rivo/tview"
)

type ClusterView struct {
}

func NewClusterView(w, h int) tview.Primitive {
	box := tview.NewBox().
		SetBorder(true).
		SetTitle("Centered Box")
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height-5, 1, true).
				AddItem(nil, 0, 1, false), width-5, 1, true).
			AddItem(nil, 0, 1, false)
	}
	return modal(box, w, h)
}
