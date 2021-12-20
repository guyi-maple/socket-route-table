package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var num int8 = 1
	fmt.Printf("%d", unsafe.Sizeof(num))
}
