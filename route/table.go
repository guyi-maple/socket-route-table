package route

import (
	"fmt"
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

func checkNode(table Table) {
	c := time.Tick(time.Duration(table.ping) * time.Second)
	go func() {
		for {
			<-c
			fmt.Printf("11111")
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
