package service

import "time"

type Environment int

const (
	Integration Environment = iota
	Stage
	Production

	IntegrationURl = "https://api.integration.openshift.com"
	StageURL       = "https://api.stage.openshift.com"
	ProdURL        = "https://api.openshift.com"

	defaultCacheTTL = 30 * time.Second
)
