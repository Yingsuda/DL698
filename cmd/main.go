package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	srcConn, err := net.Dial("tcp", "127.0.0.1:7001")
	if err != nil {
		fmt.Println("net dial src err:", err)
		return
	}

	dstConn, err := net.Dial("tcp", "192.168.50.102:28080")
	if err != nil {
		fmt.Println("net dial dst err:", err)
		return
	}

	forward := func(src, dst net.Conn) {
		defer src.Close()
		defer dst.Close()
		n, err := io.Copy(src, dst)
		if err != nil {
			fmt.Println("copy err:", err)
		}
		fmt.Println("n", n)
	}

	fmt.Println("Starting")
	go forward(srcConn, dstConn)
	go forward(dstConn, srcConn)
	select {}
}
