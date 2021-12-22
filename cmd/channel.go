package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"socket-router-table/route"
	"socket-router-table/util"
)

type Channel struct {
}

func NewChannel() Channel {
	return Channel{}
}

func (channel Channel) StartChannel(table route.Table, address string) error {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	fmt.Printf("open channel in %s \n", address)
	go func() {
		for {
			client, _ := listen.Accept()
			go func() {
				cmd, args := readCmd(client)
				fmt.Printf("cmd: %d, args: %s  \n", cmd, args)
				handleCmd(channel, table, cmd, args, client)
			}()
		}
	}()
	return nil
}

func (channel Channel) Connect(table route.Table, conn net.Conn, target string) {
	address := table.Find(target)
	if address != "" {
		// 到下级路由
		child, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Printf("connect child route error: %s \n", err.Error())
			conn.Close()
			return
		}
		go sendConnect(CONNECT, conn, child, target)
	} else {
		if table.Gateway != "" {
			// 到上级网关
			gateway, err := net.Dial("tcp", table.Gateway)
			if err != nil {
				fmt.Printf("connect gateway error: %s \n", err.Error())
				conn.Close()
				return
			}
			go sendConnect(CONNECT, conn, gateway, target)
		} else {
			// 直连
			go func() {
				dest, err := net.Dial("tcp", target)
				if err != nil {
					fmt.Printf("connect error: %s \n", err.Error())
					conn.Close()
					return
				}
				util.Forward(conn, dest)
			}()
		}
	}
}

func readCmd(client net.Conn) (CommandType, string) {
	buf := make([]byte, 4)
	_, _ = client.Read(buf[:1])
	cmd := int8(buf[0])

	_, _ = client.Read(buf[:4])
	length := int32(buf[0])
	buf = make([]byte, length)
	_, _ = io.ReadFull(client, buf[:length])
	args := string(buf)

	return CommandType(cmd), args
}

func sendConnect(cmd CommandType, conn, target net.Conn, address string) {
	target.Write([]byte{byte(cmd)})
	target.Write([]byte(address))
	util.Forward(conn, target)
}

func handleCmd(channel Channel, table route.Table, cmd CommandType, args string, conn net.Conn) {
	switch cmd {
	case CONNECT:
		channel.Connect(table, conn, args)
		break
	case PING:
		channelRoute := route.ChannelRoute{}
		err := json.Unmarshal([]byte(args), &channelRoute)
		if err != nil {
			fmt.Printf("ping cmd error: %s \n", err.Error())
		}
		table.Add(channelRoute.Converter())
		break
	default:
		conn.Close()
	}

}
