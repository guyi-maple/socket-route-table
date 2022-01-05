package channel

import (
	"fmt"
	"net"
	"socket-router-table/util"
)

type DirectChannel struct {
	name    string
	gateway string
	local   string
}

// NewDirectChannel 创建临时通道
// NewDirectChannel gateway 网关地址
// NewDirectChannel name 节点名称
func NewDirectChannel(gateway string, name string, localIp string, port int) Channel {
	return DirectChannel{
		gateway: gateway,
		name:    name,
		local:   fmt.Sprintf("%s:%d", localIp, port),
	}
}

func (channel DirectChannel) Ping() {
	conn := util.Connect(channel.gateway)
	defer conn.Close()
	if conn != nil {
		SendPing(conn, channel.name)
	}
}

func (channel DirectChannel) UpdateRoute(cidr []string) {
	conn := util.Connect(channel.gateway)
	defer conn.Close()
	if conn != nil {
		UpdateRoute(conn, cidr, channel.name, channel.local)
	}
}

func (channel DirectChannel) Forward(address string, routeAddress string, current net.Conn) {
	conn := util.Connect(routeAddress)
	defer conn.Close()
	if conn != nil {
		args := fmt.Sprintf("none|%s", address)
		Write(conn, []byte{byte(ForwardCmd), byte(len(args))})
		Write(conn, []byte(args))
		go util.Forward(current, conn)
		go util.Forward(conn, current)
	}
}

func (channel DirectChannel) ForwardGateway(address string, current net.Conn) {
	conn := util.Connect(channel.gateway)
	defer conn.Close()
	if conn != nil {
		args := fmt.Sprintf("none|%s", address)
		Write(conn, []byte{byte(ForwardCmd), byte(len(args))})
		Write(conn, []byte(args))
		go util.Forward(current, conn)
		go util.Forward(conn, current)
	}
}

func (channel DirectChannel) Direct(address string, current net.Conn) {
	dest := util.Connect(address)
	if dest != nil {
		go util.Forward(current, dest)
		go util.Forward(dest, current)
	}
}
