package channel

import (
	"net"
	"socket-router-table/route"
)

// Channel 管道接口
type Channel interface {
	// Ping 发送心跳包, 保持命令通道的长连接
	Ping()

	// UpdateRoute 向网关更新路由信息
	// UpdateRoute cidr CIDR网段列表
	UpdateRoute(cidr []string)

	// Forward 向下级路由转发
	// Forward address 目标地址
	// Forward route 下级路由地址
	// Forward current 当前连接
	Forward(address string, route *route.Route, current net.Conn)

	// ForwardGateway 向上级网关转发
	// ForwardGateway address 目标地址
	// ForwardGateway current 当前连接
	ForwardGateway(address string, current net.Conn)

	// Direct 直连
	// Direct address 目标地址
	// Direct current 当前连接
	Direct(address string, current net.Conn)
}
