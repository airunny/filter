package ip

import (
	"net"
	"testing"
)

func TestRanges(t *testing.T) {
	cidr := "192.0.2.1/24"
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%s,%v", ip, *ipNet)
}

func TestBytesNOT(t *testing.T) {
	cidr := "192.0.2.1/24"
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%v,%v_%v_%v", ipNet.IP, ipNet.Mask, BytesNOT(ipNet.Mask), BytesOR(ipNet.IP, BytesNOT(ipNet.Mask)))
}
