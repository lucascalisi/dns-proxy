package proxy

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"golang.org/x/net/dns/dnsmessage"
)

type Resolver interface {
	Solve(dnsm dnsmessage.Message) error
	GetTLSConnection() (*tls.Conn, error)
}

type Service interface {
	Solve(conn net.Conn) error
	Pass(conn net.Conn) error
}
type Cache interface {
	Solve(dnsm dnsmessage.Message) error
}

type service struct {
	resolver Resolver
	cache    Cache
}

func NewDNSProxy(r Resolver, c Cache) Service {
	return &service{r, c}
}

func unpackMsg(b []byte) (*dnsmessage.Message, error) {
	var m dnsmessage.Message
	err := m.Unpack(b)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("Unable to unpack request, invalid message.")
	}
	return &m, nil
}

func packMsg(dnsm dnsmessage.Message) ([]byte, error) {
	return dnsm.Pack()
}

func (s *service) Pass(conn net.Conn) error {
	fmt.Println("Llegué")
	defer conn.Close()
	go handleRequest(conn)
	return nil
}

func handleRequest(conn net.Conn) err {
    resolverConn, err := s.resolver.GetTLSConnection()
    if err != nil {
        log.Println("Could not get TLS Resolver connection")
        return err
    }
    defer resolverConn.Close()
    
    m, err := unpackMsg(buf[2:])
    if err != nil {
         errMsg := fmt.Errof("Invalid DNS Request: %v", err)
         log.Println(errMsg)
         return errors.New(errMsg)
     }
}

func handleResolver(conn net.Conn) {

}

