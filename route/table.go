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
	ping    int
	gateway string
	routes  map[string]Route
}

func New(ping int, gateway string, routes map[string]Route) Table {
	table := Table{
		ping:    ping,
		gateway: gateway,
		routes:  routes,
	}
	go checkNode(table)
	return table
}

func checkNode(table Table) {
	c := time.Tick(time.Duration(table.ping) * time.Second)
	go func() {
		now := time.Now().Second()
		for {
			<-c
			removeKeys := make(map[string]bool)
			for key, route := range table.routes {
				if now-route.time >= table.ping {
					removeKeys[key] = true
				}
			}
			for key, _ := range removeKeys {
				delete(table.routes, key)
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

func (table Table) add(routes []Route) {
	for _, route := range routes {
		table.routes[route.cidr] = route
	}
}

func (table Table) find(address string) string {
	for _, route := range table.routes {
		if isBelong(address, route.cidr) {
			return route.target
		}
	}
	return ""
}
