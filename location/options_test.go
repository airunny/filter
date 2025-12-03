package location

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithLanguage(t *testing.T) {
	cases := []struct {
		Language string
	}{
		{
			Language: "",
		},
		{
			Language: "zh-CN",
		},
		{
			Language: "en-US",
		},
		{
			Language: "ja-JP",
		},
		{
			Language: "en",
		},
	}

	for _, tt := range cases {
		o := Options{}
		WithLanguage(tt.Language)(&o)
		assert.Equal(t, tt.Language, o.Language)
	}
}
