package common

import "net"

// CidrStringFirstIp returns the first IP of the provided subnet.
func CidrFirstIp(cidr net.IP) net.IP {
	for j := len(cidr) - 1; j >= 0; j-- {
		cidr[j]++
		if cidr[j] > 0 {
			break
		}
	}
	return cidr
}

// CidrStringFirstIp returns the first IP of the provided subnet string. Returns an error
// if the string cannot be parsed.
func CidrStringFirstIp(ip string) (net.IP, error) {
	cidr, _, err := net.ParseCIDR(ip)
	if err != nil {
		return nil, err
	}
	return CidrFirstIp(cidr), nil
}

// Ip4BroadcastAddress returns the broadcast address for the given IP subnet.
func Ip4BroadcastAddress(n *net.IPNet) net.IP {
	// see https://groups.google.com/d/msg/golang-nuts/IrfXFTUavXE/8YwzIOBwJf0J
	ip4 := n.IP.To4()
	if ip4 == nil {
		return nil
	}
	last := make(net.IP, len(ip4))
	copy(last, ip4)
	for i := range ip4 {
		last[i] |= ^n.Mask[i]
	}
	return last
}
