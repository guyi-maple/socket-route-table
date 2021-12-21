package route

import (
	"time"
)

type ChannelRoute struct {
	Name   string   `json:"name"`
	Cidr   []string `json:"cidr"`
	Target string   `json:"target"`
}

type Route struct {
	cidr   string
	name   string
	target string
	time   int
}

type Table struct {
	Ping    int
	Gateway string
	Routes  map[string]Route
}

func New(ping int, gateway string, routes map[string]Route) Table {
	table := Table{
		Ping:    ping,
		Gateway: gateway,
		Routes:  routes,
	}
	go checkNode(table)
	return table
}

func checkNode(table Table) {
	c := time.Tick(time.Duration(table.Ping) * time.Second)
	go func() {
		now := time.Now().Second()
		for {
			<-c
			removeKeys := make(map[string]bool)
			for key, route := range table.Routes {
				if now-route.time >= table.Ping {
					removeKeys[key] = true
				}
			}
			for key, _ := range removeKeys {
				delete(table.Routes, key)
			}
		}
	}()
}

func (channel ChannelRoute) converter() []Route {
	now := time.Now().Second()
	routes := make([]Route, len(channel.Cidr))
	for index, cidr := range channel.Cidr {
		routes[index] = Route{
			cidr:   cidr,
			name:   channel.Name,
			target: channel.Target,
			time:   now,
		}
	}

	return routes
}

func (table Table) Add(routes []Route) {
	for _, route := range routes {
		table.Routes[route.cidr] = route
	}
}

func (table Table) Find(address string) string {
	for _, route := range table.Routes {
		if isBelong(address, route.cidr) {
			return route.target
		}
	}
	return ""
}
