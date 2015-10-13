package main

import (
	"github.com/rpheuts/routery/config"
	"github.com/rpheuts/routery/frontends"
	"github.com/rpheuts/routery/logger"
	"github.com/rpheuts/routery/providers"
	"github.com/rpheuts/routery/router"
	"log"
)

func main() {
	// Load config file and set logging preferences
	cfg := config.GetConfig("routery.yaml")
	logger.SetLogging(cfg.Logging.File, cfg.Logging.Path)
	log.Println("Config loaded...")

	// Initialize the providers that provide routing requests
	initializeProviders(cfg)

	// DEV: Idle until termination
	terminate := make(chan bool)
	defer close(terminate)
	<-terminate
}

func initializeProviders(cfg *config.RouteryConfig) {
	routeRequestChan := make(chan *router.RouteRequest)

	// Initialize Docker Providers
	for _, dockerConfig := range cfg.Docker {
		p := providers.DockerProvider{}
		p.Initialize(&providers.ProviderConfig{true, &dockerConfig})
		p.Provide(routeRequestChan)

		log.Printf("Registered Docker Provider. IP: %v Port: %v SSL: %v\n", dockerConfig.IP, dockerConfig.Port, dockerConfig.SSL)
	}

	// Initialize Frontends
	var listenerChannels = []chan *router.RouteRequest{}
	for _, frontendConfig := range cfg.Frontend {
		p := frontends.ForwardFrontend{}
		p.Initialize(&frontends.FrontendConfig{true,
			frontendConfig.Hostname,
			frontendConfig.Port,
			frontendConfig.SSL,
			frontendConfig.Cert,
			frontendConfig.Key,
			frontendConfig.CA,
		}, cfg)

		routeRequestListenerChan := make(chan *router.RouteRequest)
		p.Route(routeRequestListenerChan)
		listenerChannels = append(listenerChannels, routeRequestListenerChan)

		log.Printf("Registered Forwarder Frontend. Hostname: %v Port: %v\n", frontendConfig.Hostname, frontendConfig.Port)
	}

	// Start dispatching route requests
	go routeRequestDispatcher(listenerChannels, routeRequestChan)
}

func routeRequestDispatcher(listenerChannels []chan *router.RouteRequest, routeRequestChan chan *router.RouteRequest) {
	for {
		event := <-routeRequestChan

		// Broadcast on all channels
		for _, channel := range listenerChannels {
			channel <- event
		}
	}
}
