package frontends

import "github.com/rpheuts/routery/router"

type FrontendConfig struct {
	Enabled bool
	Hostname string
	Port int
}

type Frontend interface {
	Initialize(config *FrontendConfig) error
	Route(routeRequestChannel chan *router.RouteRequest) error
}