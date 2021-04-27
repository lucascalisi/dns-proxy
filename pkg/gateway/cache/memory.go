package cache

import (
	"crypto/sha256"
	"dns-proxy/pkg/domain/proxy"
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/dns/dnsmessage"
)

type memCache struct {
	mutex   sync.Mutex
	ttl     time.Duration
	entries map[string]cachedEntry
}

type cachedEntry struct {
	present bool
	msg     *dnsmessage.Message
	time    time.Time
}

func NewMemoryCache(ttl time.Duration) proxy.Cache {
	return &memCache{ttl: ttl, entries: map[string]cachedEntry{}}
}

func (mc *memCache) Get(dnsm *dnsmessage.Message) (*dnsmessage.Message, error) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	result := mc.entries[AsSha256(dnsm.Questions)]
	if result.present && time.Now().Sub(result.time) < mc.ttl {
		dnsm.Answers = result.msg.Answers
		dnsm.Authorities = result.msg.Authorities
		dnsm.Additionals = result.msg.Additionals

		return dnsm, nil
	}

	return nil, nil
}

func (mc *memCache) Store(dnsm *dnsmessage.Message, sm proxy.SolvedMsg) error {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	mc.entries[AsSha256(dnsm.Questions)] = cachedEntry{true, dnsm, time.Now()}
	return nil
}

func AsSha256(o interface{}) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", o)))

	return fmt.Sprintf("%x", h.Sum(nil))
}
