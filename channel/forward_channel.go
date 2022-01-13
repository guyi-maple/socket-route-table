package channel

import (
	"fmt"
	"net"
	"socket-router-table/util"
)

type ForwardChannel struct {
	name    string
	gateway string
	local   string
	conn    net.Conn
}

// NewForwardChannel 创建转发通道
// NewForwardChannel gateway 网关地址
// NewForwardChannel name 节点名称
func NewForwardChannel(gateway string, name string, localIp string, port int) (Channel, error) {
	conn, err := net.Dial("tcp", gateway)
	if err != nil {
		return nil, err
	}
	var channel ForwardChannel
	channel = ForwardChannel{
		gateway: gateway,
		name:    name,
		local:   fmt.Sprintf("%s:%d", localIp, port),
		conn:    conn,
	}
	return channel, nil
}

func (channel ForwardChannel) Ping() {
	SendPing(channel.conn, channel.name)
}

func (channel ForwardChannel) UpdateRoute(cidr []string) {
	UpdateRoute(channel.conn, cidr, channel.name, channel.local)
}

func (channel ForwardChannel) Forward(address string, routeAddress string, current net.Conn) {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:0", channel.local))
	if err != nil {
		fmt.Printf("listen frp forward error: %s\n", err.Error())
		return
	}
	// 连接到下级路由， 而不是网关连接
	args := fmt.Sprintf("%s|%s", listen.Addr().String(), address)
	Write(channel.conn, []byte{byte(FrpForwardCmd), byte(len(args))})
	Write(channel.conn, []byte(args))
	go func() {
		conn, _ := listen.Accept()
		go util.ForwardAndCallback(current, conn, func() {
			_ = listen.Close()
		})
		go util.Forward(conn, current)
	}()
}

func (channel ForwardChannel) ForwardGateway(address string, current net.Conn) {
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

func (channel ForwardChannel) Direct(address string, current net.Conn) {
	dest := util.Connect(address)
	if dest != nil {
		go util.Forward(current, dest)
		go util.Forward(dest, current)
	}
}
