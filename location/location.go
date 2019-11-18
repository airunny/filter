package location

import (
	"errors"
	"net"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

/**
 * https://dev.maxmind.com/geoip/geoip2/geolite2
 */
var (
	locationDB *geoip2.Reader
	once       sync.Once
)

func NewLocationWithDBFile(file string) (err error) {
	if file == "" {
		err = errors.New("empty file")
		return
	}

	once.Do(func() {
		locationDB, err = geoip2.Open(file)
		if err != nil {
			return
		}
	})

	return
}

func GetLocation(ipStr string) (country, province, city string, err error) {
	if locationDB == nil {
		err = errors.New("no initialization location")
		return
	}

	ip := net.ParseIP(ipStr)
	record, err := locationDB.City(ip)
	if err != nil {
		return
	}

	country = record.Country.Names["zh-CN"]
	if len(record.Subdivisions) > 0 {
		province = record.Subdivisions[0].Names["zh-CN"]
	}
	city = record.City.Names["zh-CN"]
	return
}

func Close() error {
	if locationDB == nil {
		return nil
	}
	return locationDB.Close()
}
