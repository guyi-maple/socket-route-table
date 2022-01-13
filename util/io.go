package util

import (
	"fmt"
	"io"
	"net"
)

// Forward IO转发
func Forward(src, dest net.Conn) {
	defer func() {
		fmt.Printf("io forward error: %s \n", recover())
	}()
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
}

// ForwardAndCallback IO转发
func ForwardAndCallback(src, dest net.Conn, callback func()) {
	defer func() {
		fmt.Printf("io forward error: %s \n", recover())
	}()
	defer src.Close()
	defer dest.Close()
	defer callback()
	io.Copy(src, dest)
}

func Connect(address string) net.Conn {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("tcp connect %s error: %s \n", address, err.Error())
		return nil
	}
	return conn
}
