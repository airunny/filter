package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	conf1 = `
{
    "m": {
        "1": {
            "filter": [
                [
                    "order_num",
                    ">",
                    10
                ],
                [
                    "freq.uid.daily",
                    "<=",
                    1
                ],
                [
                    "exp.uid.treasure_box",
                    "=",
                    "a"
                ],
                [
                    [
                        "__set_exp_id",
                        "=",
                        "uid"
                    ],
                    [
                        "__freq.uid.daily",
                        "=",
                        2000
                    ],
                    [
                        "image",
                        "=",
                        "oss/1/201905/52d6edea26365777-1496x1782.png"
                    ],
                    [
                        "title",
                        "=",
                        "支付成功"
                    ],
                    [
                        "text",
                        "=",
                        "恭喜完成首次购物"
                    ],
                    [
                        "btn_text",
                        "=",
                        "去免费拿商品"
                    ],
                    [
                        "btn_link",
                        "=",
                        "/main"
                    ],
                    [
                        "image",
                        "=",
                        "oss/1/201812/db594c683c847696-367x330.png"
                    ]
                ]
            ],
            "weight": 1,
            "priority": 1
        }
    }
}`

	conf2 = `{
    "m": {
        "1": {
            "filter_data": [
                [
                    "or",
                    "=>",
                    [
                        [
                            [
                                "order_day",
                                ">",
                                "4"
                            ],
                            [
                                "free_coupon_day",
                                ">",
                                "10"
                            ]
                        ],
                        [
                            [
                                "order_day",
                                ">",
                                "5"
                            ],
                            [
                                "free_coupon_day",
                                ">",
                                "6"
                            ]
                        ]
                    ]
                ],
                [
                    "weight",
                    "=",
                    1000
                ]
            ]
            "weight": 1,
            "priority": 1
        }
    }
}`
)

func TestTotalWeight(t *testing.T) {
	var (
		filters = make([]*Filter, 0, 10)
		total   = int64(0)
	)

	for i := 0; i < 10; i++ {
		total += int64(i)
		filters = append(filters, &Filter{
			weight: int64(i),
		})
	}
	assert.Equal(t, total, filterWeight(filters))
}

func TestPickByWeight(t *testing.T) {
	var (
		filters = make([]*Filter, 0, 10)
		total   = int64(0)
	)

	for i := 1; i <= 10; i++ {
		total += int64(i)
		filters = append(filters, &Filter{
			weight: int64(i),
		})
	}

	pickCache := make(map[int]int)
	totalCount := 100000
	for i := 0; i < totalCount; i++ {
		pickIndex := pickByWeight(filters, total)
		pickIndex++
		if _, ok := pickCache[pickIndex]; ok {
			pickCache[pickIndex]++
		} else {
			pickCache[pickIndex] = 1
		}
	}

	for k, v := range pickCache {
		expected := float64(k) / float64(total)
		got := float64(v) / float64(totalCount)
		assert.Equal(t, true, (expected-got) < 1 && (got-expected) < 1)
	}
}
