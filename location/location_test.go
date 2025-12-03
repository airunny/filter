package location

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLocation(t *testing.T) {
	country, province, city, err := GetLocation("47.107.69.99")
	assert.Nil(t, err)
	assert.Equal(t, "China", country)
	assert.Equal(t, "Guangdong", province)
	assert.Equal(t, "Shenzhen", city)

	err = NewLocationWithPath("")
	assert.True(t, reflect.DeepEqual(err, errors.New("empty location file path")))

	err = NewLocationWithPath(FilePath)
	assert.Nil(t, err)

	country, province, city, err = GetLocation("47.107.69.99")
	assert.Nil(t, err)
	assert.Equal(t, "China", country)
	assert.Equal(t, "Guangdong", province)
	assert.Equal(t, "Shenzhen", city)
}
