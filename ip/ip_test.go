package ip

import (
	"fmt"
	"net"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test001(t *testing.T) {
	ip, ipNet, err := net.ParseCIDR("192.0.2.1/24")
	assert.Nil(t, err)
	fmt.Println(ip.String(), ipNet.String(), ipNet.IP.String(), ipNet.Mask.String())
}

func TestToInt(t *testing.T) {
	ip := net.ParseIP("192.168.2.1")
	expected := 192<<24 + 168<<16 + 2<<8 + 1
	assert.Equal(t, expected, ToInt(ip))
}

func TestIntToIP(t *testing.T) {
	cases := []struct {
		IP string
	}{
		{
			IP: "0.0.0.0",
		},
		{
			IP: "192.168.2.1",
		},
		{
			IP: "255.255.255.255",
		},
	}

	for _, v := range cases {
		assert.Equal(t, v.IP, IntToIP(ToInt(net.ParseIP(v.IP))).String())
	}

}

func TestBytesNOT(t *testing.T) {
	cases := []struct {
		B        []byte
		Expected []byte
	}{
		{
			B:        []byte{0, 1, 2, 3, 4, 100, 200, 250, 255},
			Expected: []byte{255 - 0, 255 - 1, 255 - 2, 255 - 3, 255 - 4, 255 - 100, 255 - 200, 255 - 250, 255 - 255},
		},
		{
			B:        []byte{},
			Expected: []byte{},
		},
		{
			B:        []byte{255, 255, 255, 255, 255},
			Expected: []byte{0, 0, 0, 0, 0},
		},
		{
			B:        []byte{0, 0, 0, 0, 0},
			Expected: []byte{255, 255, 255, 255, 255},
		},
	}

	for _, v := range cases {
		got := BytesNOT(v.B)
		if !reflect.DeepEqual(v.Expected, got) {
			t.Errorf("expected: %v,but Got:%v\n", v.Expected, got)
		}
	}
}

func TestBytesOR(t *testing.T) {
	cases := []struct {
		A []byte
		B []byte
	}{
		{
			A: []byte{1, 11, 20, 33, 50, 100, 200, 255, 254},
			B: []byte{0, 255, 0, 100, 0, 255, 199, 233, 255},
		},
		{
			A: []byte{},
			B: []byte{},
		},
		{
			A: []byte{255, 255, 255, 255, 255},
			B: []byte{255, 255, 255, 255, 255},
		},
		{
			A: []byte{0, 0, 0, 0, 0},
			B: []byte{0, 0, 0, 0, 0},
		},
		{
			A: []byte{255, 255, 255, 255, 255},
			B: []byte{0, 0, 0, 0, 0},
		},
	}

	for index, v := range cases {
		got := BytesOR(v.A, v.B)
		for j := 0; j < len(got); j++ {
			assert.Equal(t, v.A[j]|v.B[j], got[j], index)
		}
	}
}
