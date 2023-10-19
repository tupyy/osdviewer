package service

import "time"

type Environment int

func (e Environment) String() string {
	switch e {
	case Integration:
		return "integration"
	case Stage:
		return "stage"
	case Production:
		return "production"
	default:
		return "unknown"
	}
}

const (
	Integration Environment = iota
	Stage
	Production
	Unknown

	IntegrationURl = "https://api.integration.openshift.com"
	StageURL       = "https://api.stage.openshift.com"
	ProdURL        = "https://api.openshift.com"

	defaultCacheTTL = 30 * time.Second
)
