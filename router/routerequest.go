package router

type RouteRequest struct {
	Id string
	Hostname string
	Endpoint string
	Port string
	Remove bool
}