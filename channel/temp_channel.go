package channel

import (
	"net"
	"socket-router-table/util"
)

type tempChannel struct {
	gateway string
}

// NewTempChannel 创建临时通道
// NewTempChannel gateway 网关地址
func NewTempChannel(gateway string) Channel {
	return tempChannel{
		gateway: gateway,
	}
}

func (channel tempChannel) Ping() {
	conn := util.Connect(channel.gateway)
	defer conn.Close()
	if conn != nil {
		SendPing(conn)
	}
}

func (channel tempChannel) UpdateRoute(cidr []string) {

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
