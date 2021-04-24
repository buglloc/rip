package cfg

import (
	"net"
	"time"
)

const (
	// RIP version
	Version = "0.1.1"
)

var (
	// Addr is address to listen on, ":dns" if empty.
	Addr string
	// Zones - list of acceptable zones names
	Zones []string
	// IPv4 is default IPv4 address
	IPv4 net.IP
	// IPv6 is default IPv6 address
	IPv6 net.IP
	// Upstream DNS server for proxying
	Upstream = "1.1.1.1:53"
	// UseDefault enables "strict" mode
	UseDefault bool
	AllowProxy bool
	TTL        uint32        = 0
	StickyTTL  time.Duration = 30
)
