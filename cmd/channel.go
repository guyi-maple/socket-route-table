package cmd

import (
	"io"
	"net"
)

func StartChannel(address string) error {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	for {
		client, _ := listen.Accept()
		go func() {
			defer client.Close()
			cmd, args := readCmd(client)
		}()
	}
}

func readCmd(client net.Conn) (int8, string) {
	buf := make([]byte, 4)
	_, _ = client.Read(buf[:1])
	cmd := int8(buf[0])

	_, _ = client.Read(buf[:4])
	length := int32(buf[0])
	buf = make([]byte, length)
	_, _ = io.ReadFull(client, buf[:length])
	args := string(buf)

	return cmd, args
}
