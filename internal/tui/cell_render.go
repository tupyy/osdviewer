package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tupyy/osdviewer/internal/entity"
)

const (
	leftPadding int = iota
	rightPadding
	middlePadding
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
			id = pad(id, 4, leftPadding)
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
		} else {
			cell = tview.NewTableCell("")
		}
	case 3:
		if cluster.Kind() == entity.ServiceClusterType {
			cell = tview.NewTableCell(cluster.Region())
		} else {
			cell = tview.NewTableCell("")
		}
	}

	if cell != nil {
		cell.SetTextColor(fgColor)
		cell.SetBackgroundColor(bgColor)
	}

	return cell
}

func pad(text string, count int, side int) string {
	if count <= 0 {
		return text
	}
	switch side {
	case leftPadding:
		return fmt.Sprintf("%s%s", strings.Repeat(" ", count), text)
	case rightPadding:
		return fmt.Sprintf("%s%s", text, strings.Repeat(" ", count))
	case middlePadding:
		return fmt.Sprintf("%s%s%s", strings.Repeat(" ", count), text, strings.Repeat(" ", count))
	default:
		return text
	}
}
