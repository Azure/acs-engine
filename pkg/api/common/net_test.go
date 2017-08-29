package common

import (
	"net"
	"testing"
)

type test struct {
	cidr     string
	expected string
}

func Test_CidrFirstIP(t *testing.T) {
	scenarios := []test{
		{
			cidr:     "10.0.0.0/16",
			expected: "10.0.0.1",
		},
		{
			cidr:     "10.16.32.32/27",
			expected: "10.16.32.33",
		},
	}

	for _, scenario := range scenarios {
		if first, _ := CidrStringFirstIP(scenario.cidr); first.String() != scenario.expected {
			t.Errorf("expected first ip of subnet %v to be %v but was %v", scenario.cidr, scenario.expected, first)
		}
	}
}

func Test_IP4BroadcastAddress(t *testing.T) {
	scenarios := []test{
		{
			cidr:     "10.0.0.0/16",
			expected: "10.0.255.255",
		},
		{
			cidr:     "10.16.32.32/27",
			expected: "10.16.32.63",
		},
	}

	for _, scenario := range scenarios {
		_, cidr, _ := net.ParseCIDR(scenario.cidr)
		if broadcast := IP4BroadcastAddress(cidr); broadcast.String() != scenario.expected {
			t.Errorf("expected broadcast ip of subnet %v to be %v but was %v", scenario.cidr, scenario.expected, broadcast)
		}
	}
}
