package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/onsi/gomega"
	v1 "github.com/openshift-online/ocm-sdk-go/osdfleetmgmt/v1"
	"github.com/tupyy/osdviewer/internal/entity"
)

type FleetManagerReaderMock struct {
	GetClusterFunc func(ctx context.Context, env Environment) ([]entity.Cluster, error)
}

func (r *FleetManagerReaderMock) GetClusters(ctx context.Context, env Environment) ([]entity.Cluster, error) {
	return r.GetClusterFunc(ctx, env)
}

func TestCache(t *testing.T) {
	g := gomega.NewWithT(t)
	counter := 0
	mock := &FleetManagerReaderMock{
		GetClusterFunc: func(ctx context.Context, env Environment) ([]entity.Cluster, error) {
			counter += 1
			return []entity.Cluster{createServiceCluster(fmt.Sprintf("id-%d", counter))}, nil
		},
	}

	// hit the cache
	cacheFleetManager := NewFleetManagerCache(mock, 1*time.Second)

	for i := 1; i < 4; i++ {
		clusters, err := cacheFleetManager.GetClusters(context.TODO(), Integration)
		g.Expect(len(clusters)).To(gomega.Equal(1))
		g.Expect(err).To(gomega.BeNil())
		cluster := clusters[0]
		g.Expect(cluster.ID()).To(gomega.Equal(fmt.Sprintf("id-%d", i)))

		<-time.After(3 * time.Second)
	}

}

func TestCacheHitCache(t *testing.T) {
	g := gomega.NewWithT(t)
	counter := 0
	mock := &FleetManagerReaderMock{
		GetClusterFunc: func(ctx context.Context, env Environment) ([]entity.Cluster, error) {
			counter += 1
			return []entity.Cluster{createServiceCluster(fmt.Sprintf("id-%d", counter))}, nil
		},
	}

	// hit the cache
	cacheFleetManager := NewFleetManagerCache(mock, 30*time.Second)

	for i := 1; i < 4; i++ {
		clusters, err := cacheFleetManager.GetClusters(context.TODO(), Integration)
		g.Expect(len(clusters)).To(gomega.Equal(1))
		g.Expect(err).To(gomega.BeNil())
		cluster := clusters[0]
		g.Expect(cluster.ID()).To(gomega.Equal("id-1"))

		<-time.After(1 * time.Second)
	}

}
func createServiceCluster(id string) *entity.ServiceCluster {
	svc, _ := v1.NewServiceCluster().
		ID(id).Name("name").
		Sector("sector").
		Status("state").
		ClusterManagementReference(v1.NewClusterManagementReference().ClusterId("clusterID")).
		Build()

	return entity.NewServiceCluster(svc)
}
