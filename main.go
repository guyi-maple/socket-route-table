package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"os"
	"socket-router-table/channel"
	"socket-router-table/route"
	"socket-router-table/socks"
)

type Socks5Conf struct {
	Port  int                       `yaml:"port"`
	Bind  string                    `yaml:"bind"`
	Users []socks.AuthorizationUser `yaml:"users"`
}

type ChannelConf struct {
	Bind string `yaml:"bind"`
	Port int    `yaml:"port"`
}

type Configuration struct {
	Name    string      `yaml:"name"`
	Local   string      `yaml:"local"`
	Gateway string      `yaml:"gateway"`
	Frp     bool        `yaml:"frp"`
	Subnet  []string    `yaml:"subnet"`
	Socks5  Socks5Conf  `yaml:"socks5"`
	Channel ChannelConf `yaml:"channel"`
}

func main() {

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("get user home dir error: %s", err.Error())
		return
	}

	path := fmt.Sprintf("%s/.route/main.yaml", home)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("read config file %s error: %s", path, err.Error())
		return
	}

	var conf Configuration
	err = yaml.Unmarshal(bytes, &conf)
	if err != nil {
		fmt.Printf("parse config file error: %s", err.Error())
		return
	}

	start(conf)

}

func start(conf Configuration) {
	table := route.NewTable(conf.Name, conf.Subnet)
	var chain channel.Channel
	if conf.Frp {
		chain = nil
	} else {
		chain = channel.NewDirectChannel(conf.Gateway, conf.Name, conf.Local, conf.Channel.Port)
	}

	socks5 := socks.NewSocksServer(
		socks.ServerOptions{
			Authorization: conf.Socks5.Users != nil && len(conf.Socks5.Users) > 0,
			Users:         conf.Socks5.Users,
		},
		func(conn net.Conn, ip string, port int) {
			onSocksAccept(conf, table, chain, ip, port, conn)
		},
	)

	go func() {
		err := socks5.Listen(fmt.Sprintf("%s:%d", conf.Socks5.Bind, conf.Socks5.Port))
		if err != nil {
			fmt.Printf("start socks5 server error: %s", err.Error())
			return
		}
	}()

	if conf.Channel.Port != 0 {
		go func() {
			address := fmt.Sprintf("%s:%d", conf.Channel.Bind, conf.Channel.Port)
			err := channel.ListenGatewayServer(address, onGatewayAccept)
			if err != nil {
				fmt.Printf("listen gateway server %s error: %s \n", address, err.Error())
				return
			}
		}()
	}

	// 阻塞
	select {}
}

func onGatewayAccept(conn net.Conn) {

}

func onSocksAccept(conf Configuration, table route.Table, chain channel.Channel, ip string, port int, conn net.Conn) {
	address := fmt.Sprintf("%s:%d", ip, port)
	r := table.Find(ip)
	if r == nil {
		if conf.Gateway == "" {
			chain.Direct(address, conn)
		} else {
			chain.ForwardGateway(address, conn)
		}
	} else {
		chain.Forward(address, r.Address, conn)
	}
}
