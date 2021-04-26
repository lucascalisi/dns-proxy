package main

import (
	"time"
	"dns-proxy/pkg/domain/proxy"
	"dns-proxy/pkg/gateway/cache"
	"dns-proxy/pkg/gateway/resolver"
	"dns-proxy/pkg/presenter/socket"
)

func main() {

	cache := cache.NewMemoryCache(time.Now().Add(5 * time.Second))

	resolver := resolver.NewCloudFlareResolver("1.1.1.1", 853)

	proxy := proxy.NewDNSProxy(resolver, cache)

	socket.StartTCPServer(proxy, 4545, "localhost")
}
