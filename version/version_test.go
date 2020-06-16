package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	cases := []struct {
		Compare  string
		Compared string
		Result   int
	}{
		{
			Compare:  "",
			Compared: "",
			Result:   0,
		},
		{
			Compare:  "1",
			Compared: "",
			Result:   1,
		},
		{
			Compare:  "",
			Compared: "1",
			Result:   -1,
		},
		{
			Compare:  "1.1.1",
			Compared: "1.1.1",
			Result:   0,
		},
		{
			Compare:  "1.2.1",
			Compared: "1.1.1",
			Result:   1,
		},
		{
			Compare:  "1.1.1",
			Compared: "1.2.1",
			Result:   -1,
		},
		{
			Compare:  "1.2.1",
			Compared: "1.1.2",
			Result:   1,
		},
		{
			Compare:  "1.1.1.1",
			Compared: "1.2.1",
			Result:   -1,
		},
		{
			Compare:  "1.1.1.0",
			Compared: "1.1.1",
			Result:   1,
		},
		{
			Compare:  "1.1.1",
			Compared: "1.1.1.0",
			Result:   -1,
		},
	}

	for _, v := range cases {
		assert.Equal(t, v.Result, Compare(v.Compare, v.Compared))
	}
}
