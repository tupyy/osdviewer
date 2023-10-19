package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tupyy/osdviewer/internal/entity"
)

func RenderCell(cluster entity.Cluster, col int) *tview.TableCell {
	if col > COUNT_COL {
		return nil
	}

	fgColor := tcell.ColorGreen
	bgColor := tcell.ColorBlack
	if cluster.Kind() == entity.ServiceClusterType {
		fgColor = tcell.ColorRed
	}

	var cell *tview.TableCell
	switch col {
	case 0:
		id := cluster.ID()
		if cluster.Kind() == entity.ManagementClusterType {
			id = fmt.Sprintf("     %s", id)
		}
		cell = tview.NewTableCell(id).SetExpansion(2)
	case 1:
		cell = tview.NewTableCell(cluster.State()).SetExpansion(2)
		if cluster.State() == "failed" {
			bgColor = tcell.ColorRed
		}
	case 2:
		if cluster.Kind() == entity.ServiceClusterType {
			cell = tview.NewTableCell(cluster.Sector())
		}
	case 3:
		if cluster.Kind() == entity.ServiceClusterType {
			cell = tview.NewTableCell(cluster.Region())
		}
	}

	if cell != nil {
		cell.SetTextColor(fgColor)
		cell.SetBackgroundColor(bgColor)
	}

	return cell
}
