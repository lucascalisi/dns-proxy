package proxy

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"golang.org/x/net/dns/dnsmessage"
)

type Service interface {
	Solve(um UnsolvedMsg, msjFormat string) (SolvedMsg, error)
	Direct(conn *net.Conn) error
}

type service struct {
	resolver Resolver
	mparser  MsgParser
	cache    Cache
}

func NewDNSProxy(r Resolver, mp MsgParser, c Cache) Service {
	return &service{r, mp, c}
}

func (s *service) Direct(conn *net.Conn) error {
	resolverConn, err := s.resolver.GetTLSConnection()
	if err != nil {
		log.Println("Could not get TLS Resolver connection")
		return err
	}
	defer resolverConn.Close()
	go io.Copy(resolverConn, *conn) // Holdea
	io.Copy(*conn, resolverConn)
	return nil
}

func (s *service) Solve(um UnsolvedMsg, msgFormat string) (SolvedMsg, error) {
	// Parse the UnsolvedMsg
	var parseErr error
	var dnsm *dnsmessage.Message
	if msgFormat == "udp" {
		dnsm, um, parseErr = s.mparser.ParseUPDMsg(um)
	} else if msgFormat == "tcp" {
		dnsm, parseErr = s.mparser.ParseTCPMsg(um)
	} else {
		return nil, errors.New(fmt.Sprintf("Invalid msg format: %s \n", msgFormat))
	}
	if parseErr != nil {
		log.Printf("Error parsing UnsolvedMsg: %v \n", parseErr)
		return nil, parseErr
	}

	// Log
	for _, q := range dnsm.Questions {
		log.Printf("DNS  [\033[1;36m%s\033[0m] -> : \033[1;34m%s\033[0m", msgFormat, q.Name.String())
	}

	// Check if the response is cached
	cm, cacheErr := s.cache.Get(dnsm)
	if cacheErr != nil {
		log.Printf("\033[1;33mCache error:\033[0m : %v", cacheErr)
	}

	// If cache could resolve the query, then try with the resolver
	if cm == nil {
		packed, _ := s.mparser.PackMessage(dnsm, "tcp")
		sm, err := s.resolver.Solve(packed)
		if err != nil {
			log.Printf("Resolution Error: %v \n", err)
			return nil, err
		}

		//parse response
		dnssm, err := s.mparser.ParseTCPMsg(sm)
		if err != nil {
			log.Printf("Could not parse solved message Error: %v \n", err)
		}

		err = s.cache.Store(dnssm, sm)
		if err != nil {
			log.Printf("\033[1;33mCache error:\033[0m : %v", err)
		}
		return s.mparser.PackMessage(dnssm, msgFormat)
	}
	log.Printf("\033[1;33mFound in cache\033[0m")
	return s.mparser.PackMessage(cm, msgFormat)
}
