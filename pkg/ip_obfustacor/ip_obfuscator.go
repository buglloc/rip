package ip_obfuscator

import (
	"fmt"
	"net"
)

// Based on https://github.com/OsandaMalith/IPObfuscator

func IPv4(ipStr string) []string {
	ip := net.ParseIP(ipStr).To4()
	result := make([]string, 0)
	result = append(result, fmt.Sprintf(
		"%d",
		(ip[0]<<24)|(ip[1]<<16)|(ip[2]<<8)|ip[3],
	))
	result = append(result, fmt.Sprintf(
		"%#x.%#x.%#x.%#x",
		ip[0], ip[1], ip[2], ip[3],
	))
	result = append(result, fmt.Sprintf(
		"%#X.%#X.%#X.%#X",
		ip[0], ip[1], ip[2], ip[3],
	))
	result = append(result, fmt.Sprintf(
		"%#X.%#x.%#X.%#x",
		ip[0], ip[1], ip[2], ip[3],
	))
	result = append(result, fmt.Sprintf(
		"%#05X.%#04X.%#03X.%#X",
		ip[0], ip[1], ip[2], ip[3],
	))
	result = append(result, fmt.Sprintf(
		"%010o.%010o.%010o.%010o",
		ip[0], ip[1], ip[2], ip[3],
	))
	result = append(result, fmt.Sprintf(
		"%010o.%010o.%010o.%d",
		ip[0], ip[1], ip[2], ip[3],
	))
	result = append(result, fmt.Sprintf(
		"%010o.%010o.%010o.%d",
		ip[0], ip[1], ip[2], ip[3],
	))
	result = append(result, fmt.Sprintf(
		"%010o.%010o.%d.%d",
		ip[0], ip[1], ip[2], ip[3],
	))
	result = append(result, fmt.Sprintf(
		"%010o.%d.%d.%d",
		ip[0], ip[1], ip[2], ip[3],
	))
	result = append(result, fmt.Sprintf(
		"%010o.%#x.%#X.%d",
		ip[0], ip[1], ip[2], ip[3],
	))

	decSuffix := (ip[1] << 16) | (ip[2] << 8) | ip[3]
	result = append(result, fmt.Sprintf(
		"%010o.%d",
		ip[0], decSuffix,
	))
	result = append(result, fmt.Sprintf(
		"%#x.%d",
		ip[0], decSuffix,
	))

	return result
}
