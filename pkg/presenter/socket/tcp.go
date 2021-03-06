package socket

import (
	"dns-proxy/pkg/domain/proxy"
	"errors"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

func StartTCPServer(proxy proxy.Service, port int, host string, direct bool, maxPoolConnection int) {
	portStr := strconv.Itoa(port)
	log.Println("Listening TCP DNS Proxy on PORT " + portStr)
	ln, err := net.Listen("tcp", host+":"+portStr)
	if err != nil {
		log.Println("Error creating listener")
		panic(err)
	}
	var conns uint64
	for {
		if conns <= uint64(maxPoolConnection-1) {
			// Holds inil a new connection is set
			conn, err := ln.Accept()
			atomic.AddUint64(&conns, 1)
			if err != nil {
				log.Println("Error creating connection")
				panic(err)
			}
			go tcpHandler(&conn, proxy, &conns, direct)
		}
	}
}

func tcpHandler(conn *net.Conn, p proxy.Service, conns *uint64, direct bool) error {
	defer (*conn).Close()
	if direct {
		err := p.Direct(conn)
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		var unsolvedMsg [2024]byte
		n, err := (*conn).Read(unsolvedMsg[:])
		if err != nil {
			log.Println("Failed to read from connection.")
			return errors.New("Failed to read from connection.")
		}
		solvedMsg, err := p.Solve(unsolvedMsg[:n], "tcp")
		if err != nil {
			log.Printf("Error solving message: %v \n", err)
			return err
		}

		(*conn).Write(solvedMsg)
	}
	atomic.AddUint64(conns, ^uint64(0))
	return nil
}
