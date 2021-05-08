package iputil

import (
	"math/big"
	"net"
)

func DecodeIp(hex string) (net.IP, bool) {
	ip, ok := big.NewInt(0).SetString(hex, 16)
	if !ok {
		return nil, false
	}
	return ip.Bytes(), true
}

func EncodeIP4(ip net.IP) string {
	ipInt := big.NewInt(0)
	ipInt.SetBytes(ip.To4())
	return ipInt.Text(16)
}

func EncodeIP6(ip net.IP) string {
	ipInt := big.NewInt(0)
	ipInt.SetBytes(ip.To16())
	return ipInt.Text(16)
}
