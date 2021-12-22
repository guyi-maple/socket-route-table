package main

import (
	"fmt"
	"net"
	"socket-router-table/cmd"
	"socket-router-table/route"
	"socket-router-table/socks"
)

func main() {
	table := route.New("main", "192.168.3.69:8888", 300, "", make(map[string]route.Route))

	channel := cmd.NewChannel()
	err := channel.StartChannel(table, ":8888")
	if err != nil {
		fmt.Printf("open channel error: %s \n", err.Error())
		return
	}

	onConnected := func(client net.Conn, address string, port int) {
		channel.Connect(table, client, fmt.Sprintf("%s:%d", address, port))
	}
	socks5 := socks.New(socks.ServerOptions{Authorization: false, Users: []socks.AuthorizationUser{}}, onConnected)
	socks5.Listen(":9999")
}
