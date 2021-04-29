package proxy

import (
	"crypto/tls"

	"golang.org/x/net/dns/dnsmessage"
)

type Msg []byte
type UnsolvedMsg = Msg
type SolvedMsg = Msg

type Resolver interface {
	Solve(um UnsolvedMsg) (Msg, error)
	GetTLSConnection() (*tls.Conn, error)
}
type Cache interface {
	Get(dnsm *dnsmessage.Message) (*dnsmessage.Message, error)
	Store(dnsm *dnsmessage.Message) error
	AutoPurge()
}

type ListUpdater interface {
	Update(source string) error
	UpdateAll() (map[string]bool, int)
}
type Blocker interface {
	IsBlocked(domain string) bool
	MockBlockedQuery(dnsm *dnsmessage.Message) *dnsmessage.Message
	Update()
}
type MsgParser interface {
	ParseUPDMsg(m Msg) (*dnsmessage.Message, UnsolvedMsg, error)
	ParseTCPMsg(m Msg) (*dnsmessage.Message, error)
	PackMessage(dnsm *dnsmessage.Message, msgFormat string) (SolvedMsg, error)
}

type Repository interface {
	SaveQuery(*dnsmessage.Message, string) error
}

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)
