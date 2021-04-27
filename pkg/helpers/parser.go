package helpers

import (
	"dns-proxy/pkg/domain/proxy"
	"errors"
	"log"

	"golang.org/x/net/dns/dnsmessage"
)

type msgParser struct {
}

func NewMsgParser() proxy.MsgParser {
	return &msgParser{}
}

func (mp *msgParser) ParseUPDMsg(m proxy.Msg) (*dnsmessage.Message, proxy.UnsolvedMsg, error) {
	var dnsm dnsmessage.Message
	err := dnsm.Unpack(m[:])
	if err != nil {
		log.Printf("Unable to parse UDP Message: %v \n", err)
		return nil, nil, errors.New("Unable to unpack request, invalid message.")
	}

	var tcpBytes []byte
	tcpBytes = make([]byte, 2)
	tcpBytes[0] = 0
	tcpBytes[1] = byte(len(m))
	m = append(tcpBytes, m...)

	return &dnsm, m, nil
}

func (mp *msgParser) ParseTCPMsg(m proxy.Msg) (*dnsmessage.Message, error) {
	var dnsm dnsmessage.Message
	err := dnsm.Unpack(m[2:])
	if err != nil {
		log.Printf("Unable to parse TCP Message: %v \n", err)
		return nil, errors.New("Unable to unpack request, invalid message.")
	}
	return &dnsm, nil
}

func (mp *msgParser) PackTCP(dnsm *dnsmessage.Message) (proxy.SolvedMsg, error) {
	m, err := dnsm.Pack()
	if err != nil {
		log.Printf("Unable to pack TCP Response: %v \n", err)
		return nil, errors.New("Unable to pack response, invalid message.")
	}
	var tcpBytes []byte
	tcpBytes = make([]byte, 2)
	tcpBytes[0] = 0
	tcpBytes[1] = byte(len(m))
	m = append(tcpBytes, m...)

	return m[:len(m)-135], nil
}

func (mp *msgParser) PackUDP(dnsm *dnsmessage.Message) (proxy.SolvedMsg, error) {
	m, err := dnsm.Pack()
	if err != nil {
		log.Printf("Unable to pack TCP Response: %v \n", err)
		return nil, errors.New("Unable to pack response, invalid message.")
	}

	return m, nil
}
