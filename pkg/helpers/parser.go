package helpers

import (
	"bytes"
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

func (mp *msgParser) PackMessage(dnsm *dnsmessage.Message, msgFormat string) (proxy.SolvedMsg, error) {
	m, err := dnsm.Pack()
	if err != nil {
		log.Printf("Unable to pack %s Response: %v \n", msgFormat, err)
		return nil, errors.New("Unable to pack response, invalid message.")
	}
	m = bytes.Trim(m, "\x00")
	size := len(m)

	if msgFormat == "tcp" {
		var tcpBytes []byte
		tcpBytes = make([]byte, 2)
		tcpBytes[0] = 0
		tcpBytes[1] = byte(len(m))

		m = append(tcpBytes, m...)
		m[1] = byte(size)
	}
	return m, nil
}
