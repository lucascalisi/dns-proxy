package cache

import (
	"errors"
	"fmt"
	"time"
	"tls-dns-proxy/pkg/domain/proxy"

	"golang.org/x/net/dns/dnsmessage"
)

func main() {
	fmt.Println("vim-go")
}

type memCache struct {
	ttl time.Time
}

func NewMemoryCache(ttl time.Time) proxy.Cache {
	return &memCache{ttl}
}

func (c *memCache) Solve(dnsm dnsmessage.Message) error {
	return errors.New("Not found")
}
