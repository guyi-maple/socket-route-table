package channel

import (
	"fmt"
	"net"
)

type CMD int8

const (
	PING CMD = 1
)

func SendPing(conn net.Conn) {
	_, err := conn.Write([]byte{byte(PING)})
	if err != nil {
		fmt.Printf("write ping cmd error: %s \n", err.Error())
	}
}
