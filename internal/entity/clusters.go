package entity

import (
	"fmt"
	"strings"

	v1 "github.com/openshift-online/ocm-sdk-go/osdfleetmgmt/v1"
)

const (
	ServiceClusterType int = iota
	ManagementClusterType
)

type Cluster interface {
	ID() string
	Name() string
	ClusterID() string
	Sector() string
	State() string
	Region() string
	Kind() int
	String() string
}

type cluster struct {
	id        string
	clusterID string
	name      string
	sector    string
	state     string
	region    string
}

type ServiceCluster struct {
	cluster
	Children []*ManagementCluster
}

func (sc *ServiceCluster) ID() string {
	return sc.id
}

func (sc *ServiceCluster) ClusterID() string {
	return sc.clusterID
}

func (sc *ServiceCluster) Name() string {
	return sc.name
}

func (sc *ServiceCluster) Sector() string {
	return sc.sector
}

func (sc *ServiceCluster) State() string {
	return sc.state
}

func (sc *ServiceCluster) Region() string {
	return sc.region
}

func (sc *ServiceCluster) Kind() int {
	return ServiceClusterType
}

func (sc *ServiceCluster) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, sc.ID())
	fmt.Fprintf(&sb, sc.ClusterID())
	fmt.Fprintf(&sb, sc.Sector())
	fmt.Fprintf(&sb, sc.State())
	fmt.Fprintf(&sb, sc.Region())
	return sb.String()
}

func (sc *ServiceCluster) AddManagementCluster(mc *ManagementCluster) {
	sc.Children = append(sc.Children, mc)
}

func NewServiceCluster(sc *v1.ServiceCluster) *ServiceCluster {
	return &ServiceCluster{
		cluster: cluster{
			id:        sc.ID(),
			clusterID: sc.ClusterManagementReference().ClusterId(),
			name:      sc.Name(),
			sector:    sc.Sector(),
			state:     sc.Status(),
			region:    sc.Region(),
		},
		Children: make([]*ManagementCluster, 0),
	}
}

type ManagementCluster struct {
	cluster
	Parent *ServiceCluster
}

func NewManagementCluster(mc *v1.ManagementCluster, parent *ServiceCluster) *ManagementCluster {
	return &ManagementCluster{
		cluster: cluster{
			id:        mc.ID(),
			clusterID: mc.ClusterManagementReference().ClusterId(),
			name:      mc.Name(),
			sector:    mc.Sector(),
			state:     mc.Status(),
			region:    mc.Region(),
		},
		Parent: parent,
	}
}

func (mc *ManagementCluster) ID() string {
	return mc.id
}

func (mc *ManagementCluster) ClusterID() string {
	return mc.clusterID
}

func (mc *ManagementCluster) Name() string {
	return mc.name
}

func (mc *ManagementCluster) Sector() string {
	return mc.sector
}

func (mc *ManagementCluster) State() string {
	return mc.state
}

func (mc *ManagementCluster) Region() string {
	return mc.region
}

func (mc *ManagementCluster) Kind() int {
	return ManagementClusterType
}

func (mc *ManagementCluster) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, mc.ID())
	fmt.Fprintf(&sb, mc.ClusterID())
	fmt.Fprintf(&sb, mc.Sector())
	fmt.Fprintf(&sb, mc.State())
	fmt.Fprintf(&sb, mc.Region())
	return sb.String()
}
