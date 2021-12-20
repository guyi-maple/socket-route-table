package route

import (
	"strconv"
	"strings"
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

func isBelong(ip string, cidr string) bool {
	ipAddr := strings.Split(ip, `.`)
	if len(ipAddr) < 4 {
		return false
	}
	cidrArr := strings.Split(cidr, `/`)
	if len(cidrArr) < 2 {
		return false
	}
	var tmp = make([]string, 0)
	for key, value := range strings.Split(`255.255.255.0`, `.`) {
		iint, _ := strconv.Atoi(value)
		iint2, _ := strconv.Atoi(ipAddr[key])
		tmp = append(tmp, strconv.Itoa(iint&iint2))
	}
	return strings.Join(tmp, `.`) == cidrArr[0]
}

func (table Table) find(address string) string {
	for _, route := range table.routes {
		if isBelong(address, route.cidr) {
			return route.target
		}
	}
	return ""
}
