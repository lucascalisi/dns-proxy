package blocker

import (
	"dns-proxy/pkg/domain/proxy"
	"log"
	"time"

	"golang.org/x/net/dns/dnsmessage"
)

type blocker struct {
	Sources []string
	List    map[string]bool
	Updater proxy.ListUpdater
}

func NewBlocker(refresh time.Duration, sources []string) proxy.Blocker {
	updater := NewUpdater(refresh, sources)
	return &blocker{
		Sources: sources,
		Updater: updater,
		List:    make(map[string]bool),
	}
}

func (b *blocker) IsBlocked(domain string) bool {
	return b.List[domain[:len(domain)-1]]
}

func (b *blocker) MockBlockedQuery(dnsm *dnsmessage.Message) *dnsmessage.Message {
	dnsm.Header.RecursionDesired = false
	dnsm.Header.Response = true
	dnsm.Header.RCode = dnsmessage.RCodeRefused
	dnsm.Additionals = []dnsmessage.Resource{}
	return dnsm
}

func (b *blocker) Update() {
	for _ = range time.Tick(time.Second) {
		list, errors := b.Updater.UpdateAll()
		if list != nil {
			b.List = list
			log.Printf("Block List [\033[1;33mUpdated\033[0m] -> : %d", len(b.List))
			log.Printf("Block List [\033[1;33mErrors\033[0m] -> : %d", errors)
		}
	}
}
