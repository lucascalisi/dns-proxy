package socket

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"tls-dns-proxy/pkg/domain/proxy"
)

func StartTCPServer(proxy proxy.Service, port int, host string) {
	portStr := strconv.Itoa(port)
	fmt.Println("Starting TCP DNS Proxy on PORT " + portStr)
	ln, err := net.Listen("tcp", host+":"+portStr)
	if err != nil {
		log.Println("error creating listener")
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("error creating connection")
			panic(err)
		}
		go proxy.Pass(conn)
	}

}
