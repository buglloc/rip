package cfg

import "net"

const (
	// RIP version
	Version = "0.1.0"
)

var (
	// Address to listen on, ":dns" if empty.
	Addr string
	// list of acceptable zones names
	Zones []string
	// default IPv4 address
	IPv4 net.IP
	// default IPv6 address
	IPv6 net.IP
	// upstream DNS server for proxying
	Upstream = "1.1.1.1:53"
	// Enable "strict" mode
	UseDefault bool
	AllowProxy bool
	TTL        uint32 = 0
)
