package tui

import (
	"github.com/rivo/tview"
	"github.com/tupyy/osdviewer/internal/entity"
)

const (
	COUNT_COL = 4
)

type ClusterTableContent struct {
	clusters []entity.Cluster
}

func (tc *ClusterTableContent) GetCell(row, column int) *tview.TableCell {
	cluster := tc.getCluster(row)
	if cluster == nil {
		return tview.NewTableCell("")
	}

	return RenderCell(cluster, column)
}

func (tc *ClusterTableContent) GetRowCount() int {
	count := 0
	for _, s := range tc.clusters {
		sc := s.(*entity.ServiceCluster)
		count += 1 + len(sc.Children)
	}
	return count
}

func (tc *ClusterTableContent) GetColumnCount() int {
	return COUNT_COL
}

// SetCell does not do anything.
func (tc *ClusterTableContent) SetCell(row, column int, cell *tview.TableCell) {
	// nop.
}

// RemoveRow does not do anything.
func (tc *ClusterTableContent) RemoveRow(row int) {
	// nop.
}

// RemoveColumn does not do anything.
func (tc *ClusterTableContent) RemoveColumn(column int) {
	// nop.
}

// InsertRow does not do anything.
func (tc *ClusterTableContent) InsertRow(row int) {
	// nop.
}

// InsertColumn does not do anything.
func (tc *ClusterTableContent) InsertColumn(column int) {
	// nop.
}

// Clear does not do anything.
func (tc *ClusterTableContent) Clear() {
	// nop.
}

func (tc *ClusterTableContent) getCluster(row int) entity.Cluster {
	i := 0
	for _, s := range tc.clusters {
		if i == row {
			return s
		}
		i += 1

		sc := s.(*entity.ServiceCluster)
		for _, mc := range sc.Children {
			if i == row {
				return mc
			}
			i += 1
		}
	}

	return nil
}

func NewTableContent(clusters []entity.Cluster) *ClusterTableContent {
	return &ClusterTableContent{clusters: clusters}
}
