package socks

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func handleCmd(client net.Conn) bool {
	buf := make([]byte, 256)

	// 读取连接信息
	_, err := io.ReadFull(client, buf[:2])
	if err != nil {
		fmt.Printf("socks5 read cmd error: %s \n", err.Error())
		return false
	}
	ver, methods := int(buf[0]), int(buf[1])
	if ver != 5 {
		fmt.Printf("socks5 not support protocol ver %b \n", ver)
		return false
	}
	n, err := io.ReadFull(client, buf[:methods])
	if n != methods {
		fmt.Printf("socks5 reading methods error: %s \n", err.Error())
		return false
	}

	return true
}

func handleAuth(client net.Conn, options ServerOptions) bool {
	// 不需要加密
	n, err := client.Write([]byte{0x05, 0x00})
	if n != 2 || err != nil {
		fmt.Printf("socks5 write rsp err: %s \n", err.Error())
		return false
	}
	if err != nil {
		fmt.Printf("socks5 check password error: %s \n" + err.Error())
		return false
	}

	return true
}

func handleConnect(client net.Conn) (bool, string, int) {
	buf := make([]byte, 256)

	n, err := io.ReadFull(client, buf[:4])
	if n != 4 {
		fmt.Printf("socks5 read header error: %s \n", err.Error())
		return false, "", 0
	}

	ver, cmd, _, atyp := int(buf[0]), buf[1], buf[2], buf[3]
	if ver != 5 || cmd != 1 {
		fmt.Printf("socks5 invalid ver/cmd")
		return false, "", 0
	}

	addr := ""
	switch atyp {
	case 1:
		n, err = io.ReadFull(client, buf[:4])
		if n != 4 {
			fmt.Printf("socks5 invalid IPv4: %s \n", err.Error())
			return false, "", 0
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])

	case 3:
		n, err = io.ReadFull(client, buf[:1])
		if n != 1 {
			fmt.Printf("socks5 invalid hostname: %s \n", err.Error())
			return false, "", 0
		}
		addrLen := int(buf[0])
		n, err = io.ReadFull(client, buf[:addrLen])
		if n != addrLen {
			fmt.Printf("socks5 invalid hostname: %s \n", err.Error())
			return false, "", 0
		}
		addr = string(buf[:addrLen])

	case 4:
		fmt.Printf("socks5 IPv6: no supported yet")
		return false, "", 0
	default:
		fmt.Printf("socks5 invalid atyp")
		return false, "", 0
	}

	n, err = io.ReadFull(client, buf[:2])
	if n != 2 {
		fmt.Printf("read port error: %s \n", err.Error())
		return false, "", 0
	}
	port := int(binary.BigEndian.Uint16(buf[:2]))

	n, err = client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		fmt.Printf("socks5 write resp error: %s \n", err.Error())
		return false, "", 0
	}

	return true, addr, port
}
