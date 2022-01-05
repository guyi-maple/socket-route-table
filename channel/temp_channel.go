package channel

import (
	"net"
	"socket-router-table/util"
)

type tempChannel struct {
	name    string
	gateway string
}

// NewTempChannel 创建临时通道
// NewTempChannel gateway 网关地址
// NewTempChannel name 节点名称
func NewTempChannel(gateway string, name string) Channel {
	return tempChannel{
		gateway: gateway,
	}
}

func (channel tempChannel) Ping() {
	conn := util.Connect(channel.gateway)
	defer conn.Close()
	if conn != nil {
		SendPing(conn, channel.name)
	}
}

func (channel tempChannel) UpdateRoute(cidr []string) {
	conn := util.Connect(channel.gateway)
	defer conn.Close()
	if conn != nil {
		UpdateRoute(conn, cidr, channel.name)
	}
}

func (channel tempChannel) Forward(address string, routeAddress string, current net.Conn) {

}

func (channel tempChannel) ForwardGateway(address string, current net.Conn) {

}

func (channel tempChannel) Direct(address string, current net.Conn) {
	dest := util.Connect(address)
	if dest != nil {
		go util.Forward(current, dest)
		go util.Forward(dest, current)
	}
}
