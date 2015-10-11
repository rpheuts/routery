package providers

import (
	"github.com/rpheuts/routery/config"
	"github.com/rpheuts/routery/router"
)

type ProviderConfig struct {
	Enabled bool
	Docker *config.DockerConfig
}

type Provider interface {
	Initialize(config *ProviderConfig) error
	Provide(routeRequestChannel chan *router.RouteRequest) error
}
