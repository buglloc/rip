package ip_utils

import (
	"math/big"
	"net"

	log "github.com/buglloc/simplelog"
)

const hexDigit = "0123456789abcdef"

func HexToIp(hex string) net.IP {
	ip, ok := big.NewInt(0).SetString(hex, 16)
	if !ok {
		log.Error("failed to parse base16 ip", "ip", hex)
		return nil
	}
	return ip.Bytes()
}

func Ip4ToHex(ip net.IP) string {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(ip.To4())
	return hexString(IPv4Int.Bytes())
}

func Ip6ToHex(ip net.IP) string {
	IPv6Int := big.NewInt(0)
	IPv6Int.SetBytes(ip.To16())
	return hexString(IPv6Int.Bytes())
}

func hexString(b []byte) string {
	s := make([]byte, len(b)*2)
	for i, tn := range b {
		s[i*2], s[i*2+1] = hexDigit[tn>>4], hexDigit[tn&0xf]
	}
	return string(s)
}
