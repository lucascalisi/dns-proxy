package blocker

import (
	"dns-proxy/pkg/domain/proxy"
	"time"

	"golang.org/x/net/dns/dnsmessage"
)

type Blocker interface {
	IsBlocked(domain string) bool
}

type blocker struct {
	Sources []string
	Updater Updater
}

func NewBlocker(refresh time.Duration) proxy.Blocker {
	return &blocker{
		Sources: []string{},
		Updater: &updater{refresh},
	}
}

func (b *blocker) IsBlocked(domain string) bool {
	var list map[string]bool
	list = make(map[string]bool)
	list["lucascontre.site"] = true
	list["tunnel.us.ngrok.com"] = true
	list["ngrok.io"] = true
	list["lanacion.com"] = true
	list["addons-pa.clients6.google.com"] = true

	return list[domain[:len(domain)-1]]
}

func (b *blocker) MockBlockedQuery(dnsm *dnsmessage.Message) *dnsmessage.Message {
	dnsm.Header.RecursionDesired = false
	dnsm.Header.Response = true
	dnsm.Header.RCode = dnsmessage.RCodeRefused
	dnsm.Additionals = []dnsmessage.Resource{}
	return dnsm
}
