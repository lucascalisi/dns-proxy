package main

import (
	"time"
	"dns-proxy/pkg/domain/proxy"
	"dns-proxy/pkg/gateway/cache"
	"dns-proxy/pkg/gateway/resolver"
	"dns-proxy/pkg/helpers"
	"dns-proxy/pkg/presenter/socket"
	"github.com/sevlyar/go-daemon"
    "log"
)

func main() {
     cntxt := goDaemon("/tmp/dns-proxy.pid", "/tmp", "/tmp/dns-proxy.log", []string{"dns-proxy"})
     d, err := cntxt.Reborn()
     if err != nil {
         log.Fatal("Unable to run: ", err)
     }

     if d != nil {
         return
     }
     defer cntxt.Release()

	config := GetConfig()
	cache := cache.NewMemoryCache(time.Now().Add(time.Duration(config.CACHE_TLL) * time.Second))
	resolver := resolver.NewCloudFlareResolver("1.1.1.1", 853, config.RESOLVER_READ_TO)
	parser := helpers.NewMsgParser()

	proxy := proxy.NewDNSProxy(resolver, parser, cache)
	go socket.StarUDPtServer(proxy, config.UDP_PORT, "0.0.0.0")
	socket.StartTCPServer(proxy, config.TCP_PORT, "0.0.0.0", config.TCP_DIRECT, config.TCP_MAX_CONN_POOL)

}

func goDaemon(pidName string, workDir string, logFile string, args []string) *daemon.Context {
	ctx := &daemon.Context{
		PidFileName: pidName,
		PidFilePerm: 0644,
		LogFileName: logFile,
		LogFilePerm: 0640,
		WorkDir:     workDir,
		Umask:       027,
	}

	return ctx
}
