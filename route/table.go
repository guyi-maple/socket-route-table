package route

import (
	"encoding/json"
	"fmt"
	"net"
	"socket-router-table/cons"
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
	Name    string
	Ip      string
	Ping    int
	Gateway string
	Routes  map[string]Route
}

func New(name string, ip string, ping int, gateway string, routes map[string]Route) Table {
	table := Table{
		Name:    name,
		Ip:      ip,
		Ping:    ping,
		Gateway: gateway,
		Routes:  routes,
	}
	if table.Gateway != "" {
		go reportRoute(table)
	}
	go checkNode(table)
	return table
}

func reportRoute(table Table) {
	c := time.Tick(time.Duration(table.Ping) * time.Second)
	go func() {
		for {
			<-c
			ciders := make([]string, len(table.Routes))
			for cidr := range table.Routes {
				ciders = append(ciders, cidr)
			}
			routes := ChannelRoute{
				Name:   table.Name,
				Cidr:   ciders,
				Target: table.Ip,
			}
			conn, err := net.Dial("tcp", table.Gateway)
			if err != nil {
				fmt.Printf("report routes error: %s \n", err.Error())
				return
			}
			conn.Write([]byte{byte(cons.PING)})
			bytes, _ := json.Marshal(routes)
			conn.Write(bytes)
			conn.Close()
		}
	}()
}

func checkNode(table Table) {
	c := time.Tick(time.Duration(table.Ping) * time.Second)
	go func() {
		for {
			<-c
			now := time.Now().Second()
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

func (channel ChannelRoute) Converter() []Route {
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
