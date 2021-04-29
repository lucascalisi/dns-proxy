package main

import (
	"bufio"
	"dns-proxy/pkg/domain/proxy"
	"dns-proxy/pkg/gateway/blocker"
	"dns-proxy/pkg/gateway/cache"
	"dns-proxy/pkg/gateway/repository"
	"dns-proxy/pkg/gateway/resolver"
	"dns-proxy/pkg/helpers"
	"dns-proxy/pkg/presenter/socket"
	"log"
	"os"
	"time"

	"github.com/sevlyar/go-daemon"
)

func main() {
	if isNotRunningInDockerContainer() {
		cntxt := goDaemon("/tmp/dns-proxy.pid", "/tmp", "/tmp/dns-proxy.log", []string{"dns-proxy"})
		d, err := cntxt.Reborn()
		if err != nil {
			log.Fatal("Unable to run: ", err)
		}

		if d != nil {
			return
		}
		defer cntxt.Release()
	}

	config := GetConfig()
	cache := cache.NewMemoryCache(time.Duration(config.CACHE_TLL) * time.Second)
	go cache.AutoPurge()
	resolver := resolver.NewDNSOverTlsResolver("1.1.1.1", 853, config.RESOLVER_READ_TO)
	parser := helpers.NewMsgParser()

	var sources []string
	sources, _ = source("/Users/lcalisi/dns-proxy/blocker/lists.list")
	blocker := blocker.NewBlocker(time.Duration(10)*time.Minute, sources)
	repository, err := repository.NewRepository("./dns-proxy.db")
	if err != nil {
		log.Fatalf("could not open sqlite database: %v", err)
	}
	go blocker.Update()

	proxy := proxy.NewDNSProxy(resolver, parser, cache, blocker, repository)
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

func isNotRunningInDockerContainer() bool {
	// docker creates a .dockerenv file at the root
	// of the directory tree inside the container.
	// if this file exists then the viewer is running
	// from inside a container so return true

	if _, err := os.Stat("/.dockerenv"); err == nil {
		return false
	}

	return true
}

func source(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open source list")
	}

	scanner := bufio.NewScanner(file)
	var sources []string

	for scanner.Scan() {
		if scanner.Text()[:1] != "#" {
			sources = append(sources, scanner.Text())
		}
	}

	return sources, nil

}
