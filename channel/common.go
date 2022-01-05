package channel

import (
	"encoding/json"
	"fmt"
	"net"
)

type CMD int8

type RouteInfo struct {
	cidr  []string
	name  string
	local string
}

const (
	PingCmd        CMD = 1
	UpdateRouteCmd CMD = 2
	ForwardCmd     CMD = 3
)

func Write(conn net.Conn, bytes []byte) bool {
	_, err := conn.Write(bytes)
	if err != nil {
		fmt.Printf("write ping cmd error: %s \n", err.Error())
		return false
	}
	return true
}

func SendPing(conn net.Conn, name string) {
	if Write(conn, []byte{byte(PingCmd), byte(len(name))}) {
		Write(conn, []byte(name))
	}
}

func UpdateRoute(conn net.Conn, cidr []string, name string, local string) {
	info := RouteInfo{
		name:  name,
		cidr:  cidr,
		local: local,
	}
	args, err := json.Marshal(info)
	if err != nil {
		fmt.Printf("marshal route info error: %s \n", err.Error())
		return
	}

	if Write(conn, []byte{byte(UpdateRouteCmd), byte(len(args))}) {
		Write(conn, args)
	}
}
