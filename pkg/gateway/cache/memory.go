package cache

import (
	"crypto/sha256"
	"dns-proxy/pkg/domain/proxy"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/net/dns/dnsmessage"
)

type memCache struct {
	ttl     time.Duration
	mx      sync.RWMutex
	entries map[string]cacheValue
}

type cacheValue struct {
	msg        *dnsmessage.Message
	expiration time.Time
}

func NewMemoryCache(ttl time.Duration) proxy.Cache {
	c := memCache{
		ttl:     ttl,
		entries: map[string]cacheValue{},
	}
	return &c
}

func (mc *memCache) AutoPurge() {
	for now := range time.Tick(time.Second) {
		for key, cValue := range mc.entries {
			if cValue.expiration.Before(now) {
				log.Println(cValue.msg.Questions[0].Name.String())
				mc.mx.Lock()
				log.Printf("Clearing entry: %v \n", key)
				delete(mc.entries, key)
				mc.mx.Unlock()
			}
		}
	}
}

func (mc *memCache) Get(dnsm *dnsmessage.Message) (*dnsmessage.Message, error) {
	mc.mx.RLock()
	defer mc.mx.RUnlock()
	if cValue, ok := mc.entries[mc.hashKey(dnsm.Questions)]; ok {
		return cValue.msg, nil
	}
	return nil, nil
}

func (mc *memCache) Store(dnsm *dnsmessage.Message) error {
	mc.mx.Lock()
	mc.entries[mc.hashKey(dnsm.Questions)] = cacheValue{dnsm, time.Now().Add(mc.ttl)}
	mc.mx.Unlock()
	return nil
}

func (mc *memCache) hashKey(questions []dnsmessage.Question) string {

	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", questions)))
	return fmt.Sprintf("%x", h.Sum(nil))

}
