package channel

import (
	"fmt"
	"net"
)

func ListenGatewayServer(address string, onAccept func(conn net.Conn)) error {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	fmt.Printf("listen gateway server in %s \n", address)
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("gateway server accept error: %s \n", err.Error())
		} else {
			go onAccept(conn)
		}
	}
}
