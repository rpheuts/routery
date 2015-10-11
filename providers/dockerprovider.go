package providers

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/rpheuts/routery/router"
	"log"
	"fmt"
	"strings"
)

type DockerProvider struct {
	config *ProviderConfig
	client *docker.Client
	routeRequestChannel chan *router.RouteRequest
}

func (provider *DockerProvider) Initialize(config *ProviderConfig) error {
	provider.config = config
	var err error

	endpoint := fmt.Sprintf("tcp://%v:%v", provider.config.Docker.IP, provider.config.Docker.Port)

	if provider.config.Docker.SSL {
		provider.client, err = docker.NewTLSClient(endpoint,
			provider.config.Docker.Cert,
			provider.config.Docker.Key,
			provider.config.Docker.CA);
	} else {
		provider.client, err = docker.NewClient(endpoint)
	}

	if  err != nil {
		log.Printf("Failed to create a client for docker, error: %s\n", err)
		return err
	}

	return nil
}

func (provider *DockerProvider) Provide(routeRequestChannel chan *router.RouteRequest) error {
	if !provider.config.Enabled {
		return nil
	}

	// Ensure endpoint is reachable
	err := provider.client.Ping()
	if err != nil {
		log.Printf("Docker connection error %v\n", err)
		return err
	}

	// Subscribe to the Docker events for this endpoint
	dockerEvents := make(chan *docker.APIEvents)
	provider.client.AddEventListener(dockerEvents)
	provider.routeRequestChannel = routeRequestChannel

	// Generate routes for the existing containers
	go provider.generateExistingRouteRequests()

	// Start watching Docker events
	go provider.watchDockerEvents(dockerEvents)

	return nil
}

func (provider *DockerProvider) watchDockerEvents(dockerEvents chan *docker.APIEvents) {
	for {
		event := <-dockerEvents
		if event == nil {
			log.Println("Recived nil response from Docker host")
		}
		if event.Status == "start" || event.Status == "die" {
			provider.generateRouteRequests(event)
		}
	}
}

func (provider *DockerProvider) generateRouteRequests(event *docker.APIEvents) {
	// If it's being removed we don't care about the details
	if event.Status == "die" {
		provider.routeRequestChannel <- &router.RouteRequest{event.ID, "", "", "", true}
		return;
	}

	// Get the container by ID and generate specific route-request
	if container, err := provider.client.InspectContainer(event.ID); err == nil {
		// Loop through the ports and request routes
		for dockerPort, mapping := range container.NetworkSettings.Ports {
			port := strings.Split(string(dockerPort), "/")[0]
			name := strings.Split(container.Name, "/")[1]

			// If there is a port 80 we just use the container name, otherwise append the port number to the name
			if (port != "80") {
				name = fmt.Sprintf("%v-%v", name, port)
			}

			request := router.RouteRequest{event.ID, name, provider.config.Docker.IP, mapping[0].HostPort, false}
			provider.routeRequestChannel <- &request
			log.Printf("Route request generated. %v\n", request)
		}
	} else {
		log.Printf("Error getting container details: %v\n", err)
	}

	return;
}

func (provider *DockerProvider) generateExistingRouteRequests() {

	if containers, err := provider.client.ListContainers(docker.ListContainersOptions{All: false}); err == nil {
		for _, container := range containers {
			event := docker.APIEvents{"start", container.ID, "", container.Created}
			provider.generateRouteRequests(&event)

			log.Printf("Generating route request for existing container. %v", event)
		}
	} else {
		log.Printf("Failed to get existing containers: %v", err)
	}
}
