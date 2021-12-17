package socks

import (
	"fmt"
	"net"
	"socket-router-table/util"
)

type AuthorizationUser struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type ServerOptions struct {
	Authorization bool                `yaml:"authorization"` //是否需要认证
	Users         []AuthorizationUser `yaml:"users"`
}

type Server struct {
	Options     ServerOptions
	OnConnected func(net.Conn, string, int)
}

func New(options ServerOptions, onConnected func(net.Conn, string, int)) Server {
	return Server{Options: options, OnConnected: onConnected}
}

func (server Server) Listen(localAddr string) error {
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		return err
	}
	fmt.Printf("listen socks5 server in %s \n", localAddr)
	for {
		client, err := listener.Accept()
		if err != nil {
			fmt.Printf("socks5  cleint accept failed: %s \n", err.Error())
			continue
		}
		go onAccept(client, server)
	}
}

func onAccept(client net.Conn, server Server) {
	if !handleCmd(client) {
		_ = client.Close()
		return
	}
	if !handleAuth(client, server.Options) {
		_ = client.Close()
		return
	}
	boo, addr, port := handleConnect(client)
	if !boo {
		_ = client.Close()
		return
	}

	if server.OnConnected != nil {
		go server.OnConnected(client, addr, port)
	} else {
		dist, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
		if err != nil {
			fmt.Printf("socks5 connect remote error: %s \n", err.Error())
			_ = client.Close()
			return
		}
		go util.Forward(client, dist)
		go util.Forward(dist, client)
	}
}
