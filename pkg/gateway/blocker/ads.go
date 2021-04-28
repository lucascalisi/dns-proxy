package blocker

import (
	"dns-proxy/pkg/domain/proxy"
	"time"
)

type Blocker interface {
	IsBlocked(domain string) bool
}

type adsBlocker struct {
	Sources []string
	Updater Updater
}

func NewAdsBlocker(refresh time.Duration) proxy.Blocker {
	return &adsBlocker{
		Sources: []string{},
		Updater: &updater{refresh},
	}
}

func (b *adsBlocker) IsBlocked(domain string) bool {
	var list map[string]bool
	list = make(map[string]bool)
	list["lucascontre.site"] = true
	list["tunnel.us.ngrok.com"] = true
	list["ngrok.io"] = true
	list["lanacion.com"] = true
	list["addons-pa.clients6.google.com"] = true

	return list[domain[:len(domain)-1]]
}
