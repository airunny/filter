package location

import (
	"errors"
	"net"
	"path"
	"runtime"

	"github.com/oschwald/geoip2-golang"
)

// https://dev.maxmind.com/geoip/geoip2/geolite2

var (
	locationDB *geoip2.Reader
	FileDir    = path.Join(path.Dir(getCurrentFilePath()), "deps")
	FileName   = "GeoLite2-City.mmdb"
	FilePath   = path.Join(FileDir, FileName)
)

func init() {
	var err error
	locationDB, err = geoip2.Open(FilePath)
	if err != nil {
		panic(err)
	}
}

func NewLocationWithPath(file string) error {
	if file == "" {
		return errors.New("empty location file path")
	}

	if locationDB != nil {
		_ = locationDB.Close()
	}

	var err error
	locationDB, err = geoip2.Open(file)
	return err
}

func GetReader() (*geoip2.Reader, bool) {
	return locationDB, locationDB != nil
}

func GetLocation(ipStr string, opts ...Option) (country, province, city string, err error) {
	if locationDB == nil {
		err = errors.New("no init location")
		return
	}

	o := &Options{
		Language: "en",
	}

	for _, opt := range opts {
		opt(o)
	}

	ip := net.ParseIP(ipStr)
	record, err := locationDB.City(ip)
	if err != nil {
		return
	}

	country = record.Country.Names[o.Language]
	if len(record.Subdivisions) > 0 {
		province = record.Subdivisions[0].Names[o.Language]
	}
	city = record.City.Names[o.Language]
	return
}

func Close() error {
	if locationDB == nil {
		return nil
	}
	return locationDB.Close()
}

func getCurrentFilePath() string {
	_, filePath, _, _ := runtime.Caller(1)
	return filePath
}
