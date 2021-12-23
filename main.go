package main

import (
	"fmt"
	"net"
	"socket-router-table/cmd"
	"socket-router-table/route"
	"socket-router-table/socks"
)

func main() {
	conf := GetConf("./conf.yaml")

	table := route.New(conf.Name, fmt.Sprintf("%s:%d", conf.LocalIp, conf.ChannelPort), conf.Ping, conf.Gateway)

	channel := cmd.NewChannel()
	err := channel.StartChannel(table, fmt.Sprintf(":%d", conf.ChannelPort))
	if err != nil {
		fmt.Printf("open channel error: %s \n", err.Error())
		return
	}

	onConnected := func(client net.Conn, address string, port int) {
		channel.Connect(table, client, fmt.Sprintf("%s:%d", address, port))
	}
	socks5 := socks.New(socks.ServerOptions{Authorization: false, Users: []socks.AuthorizationUser{}}, onConnected)
	err = socks5.Listen(fmt.Sprintf(":%d", conf.Socks5Port))
	if err != nil {
		fmt.Printf("open socks5 server error: %s \n", err.Error())
		return
	}
}
