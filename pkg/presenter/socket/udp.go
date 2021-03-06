package socket

import (
	"log"
	"net"
	"strconv"

	"dns-proxy/pkg/domain/proxy"
)

func StarUDPtServer(proxy proxy.Service, port int, host string) {
	portStr := strconv.Itoa(port)
	log.Println("Listening UDP DNS Proxy on PORT " + portStr)
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: port})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	udpHandler(conn, proxy)
}

func udpHandler(conn *net.UDPConn, p proxy.Service) {
	for {
		unsolvedMsg := make([]byte, 2048)
		n, addr, err := conn.ReadFromUDP(unsolvedMsg)
		if err != nil {
			log.Println("Failed to read from connection.")
		}

		solvedMsg, proxyErr := p.Solve(unsolvedMsg[:n], "udp")
		if proxyErr != nil {
			log.Printf("Error solving message: %v \n", proxyErr)
		}

		_, err = conn.WriteToUDP(solvedMsg, addr)
		if err != nil {
			log.Println(err)
		}
	}
}
