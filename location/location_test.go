package location

import (
	"testing"
)

func TestGetLocation(t *testing.T) {
	err := NewLocationWithDBFile("/Users/Leo/Desktop/GeoLite2-City/GeoLite2-City.mmdb")
	if err != nil {
		t.Errorf("NewLocationWithDBFile Error: %v", err)
		return
	}
	defer Close()
	country, province, city, err := GetLocation("47.107.69.99")
	if err != nil {
		t.Errorf("GetLocation Error: %v", err)
		return
	}
	t.Logf("%v_%v_%v", country, province, city)
}
