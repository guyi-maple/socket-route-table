package util

import (
	"fmt"
	"io"
	"net"
)

func Forward(src, dest net.Conn) {
	defer func() {
		fmt.Printf("io forward error: %s \n", recover())
	}()
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
}
