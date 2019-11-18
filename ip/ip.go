package ip

import (
	"fmt"
	"net"
)

type Range struct {
	start int
	end   int
}

func (r Range) String() string {
	return fmt.Sprintf("%v-%v", IntToIP(r.start), IntToIP(r.end))
}

func ToInt(address net.IP) (intP int) {
	var ip1, ip2, ip3, ip4 int
	_, _ = fmt.Sscanf(address.To4().String(), "%d.%d.%d.%d", &ip1, &ip2, &ip3, &ip4)
	return ip1<<24 + ip2<<16 + ip3<<8 + ip4
}

func IntToIP(intP int) (address net.IP) {
	var ip1, ip2, ip3, ip4 int
	ip1 = intP & 0xFF000000
	ip2 = intP & 0x00FF0000
	ip3 = intP & 0x0000FF00
	ip4 = intP & 0x000000FF

	return net.ParseIP(fmt.Sprintf("%d.%d.%d.%d", ip1, ip2, ip3, ip4))
}

func InRange(rs []Range, ipAddress net.IP) bool {
	intIP := ToInt(ipAddress)
	for _, r := range rs {
		if intIP >= r.start && intIP <= r.end {
			return true
		}
	}

	return false
}

func Ranges(ipCIDR ...string) ([]Range, error) {
	ranges := make([]Range, 0, len(ipCIDR))

	for _, cidr := range ipCIDR {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}

		ranges = append(ranges, Range{
			start: ToInt(ipNet.IP),
			end:   ToInt(BytesOR(ipNet.IP, BytesNOT(ipNet.Mask))),
		})
	}

	return ranges, nil
}

func BytesOR(a, b []byte) (c []byte) {
	g, l := a, b
	if len(g) < len(l) {
		g, l = l, g
	}

	lg := len(g)
	ll := len(l)
	gsl := lg - ll

	c = make([]byte, lg)
	for i, _ := range g {
		if i < gsl {
			c[i] = g[i]
			continue
		}
		c[i] = l[i-gsl] | g[i]
	}
	return
}

func BytesNOT(a []byte) (b []byte) {
	b = make([]byte, len(a))
	for i, _ := range a {
		b[i] = ^a[i]
	}
	return
}
