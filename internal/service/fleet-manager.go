package service

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/openshift-online/ocm-sdk-go"
	v1 "github.com/openshift-online/ocm-sdk-go/osdfleetmgmt/v1"
	"github.com/tupyy/osdviewer/internal/entity"
)

type FleetManager struct {
	token string
}

func NewFleetManager(token string) *FleetManager {
	return &FleetManager{
		token: token,
	}
}

func (fm *FleetManager) GetClusters(ctx context.Context, env Environment) ([]entity.Cluster, error) {
	url := ""
	switch env {
	case Integration:
		url = IntegrationURl
	case Stage:
		url = StageURL
	case Production:
		url = ProdURL
	default:
		return nil, fmt.Errorf("unknown environment")
	}

	client, err := fm.getOCMClient(url, fm.token)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	clusters, err := fm.getServiceClusters(client)
	if err != nil {
		return nil, err
	}

	return clusters, nil
}

func (fm *FleetManager) getServiceClusters(client *sdk.Connection) ([]entity.Cluster, error) {
	serviceClusters, err := client.OSDFleetMgmt().V1().ServiceClusters().List().Send()
	if err != nil {
		return nil, err
	}

	mgmtClusters, err := client.OSDFleetMgmt().V1().ManagementClusters().List().Send()
	if err != nil {
		return nil, err
	}
	sv := make([]entity.Cluster, 0, serviceClusters.Items().Len())
	serviceClusters.Items().Each(func(item *v1.ServiceCluster) bool {
		serviceCluster := entity.NewServiceCluster(item)
		mgmtClusters.Items().Each(func(mc *v1.ManagementCluster) bool {
			p, ok := mc.GetParent()
			if !ok {
				return true
			}
			if strings.Contains(p.Href(), serviceCluster.ID()) {
				serviceCluster.AddManagementCluster(entity.NewManagementCluster(mc, serviceCluster))
			}
			return true
		})

		sv = append(sv, serviceCluster)
		return true
	})

	return sv, nil
}

func (fm *FleetManager) getOCMClient(url string, token string) (*sdk.Connection, error) {
	connection, err := sdk.NewConnectionBuilder().
		URL(url).
		Tokens(token).
		Build()
	if err != nil {
		return nil, err
	}
	return connection, nil
}
