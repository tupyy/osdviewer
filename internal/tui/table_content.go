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
	filter   func(entity.Cluster) bool
}

func NewTableContent(clusters []entity.Cluster) *ClusterTableContent {
	return &ClusterTableContent{
		clusters: clusters,
		filter:   func(c entity.Cluster) bool { return true },
	}
}

func NewTableContentWithFilter(clusters []entity.Cluster, filter func(entity.Cluster) bool) *ClusterTableContent {
	return &ClusterTableContent{
		clusters: clusters,
		filter:   filter,
	}
}

func (tc *ClusterTableContent) GetCell(row, column int) *tview.TableCell {
	cluster := tc.getCluster(row)
	if cluster == nil {
		return tview.NewTableCell("")
	}

	return RenderCell(cluster, column)
}

func (tc *ClusterTableContent) GetRowCount() int {
	clusters := tc.applyFilter(tc.filter)
	count := 0
	for _, s := range clusters {
		sc := s.(*entity.ServiceCluster)
		count += 1 + len(sc.Children)
	}
	return count
}

func (tc *ClusterTableContent) GetColumnCount() int {
	return COUNT_COL
}

func (tc *ClusterTableContent) getCluster(row int) entity.Cluster {
	clusters := tc.applyFilter(tc.filter)
	i := 0
	for _, s := range clusters {
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

func (tc *ClusterTableContent) SetFilter(filter func(entity.Cluster) bool) {
	tc.filter = filter
}

func (tc *ClusterTableContent) applyFilter(filter func(c entity.Cluster) bool) []entity.Cluster {
	filteredClusters := make([]entity.Cluster, 0, len(tc.clusters))
	for _, c := range tc.clusters {
		if filter(c) {
			filteredClusters = append(filteredClusters, c)
		}
	}
	return filteredClusters
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
