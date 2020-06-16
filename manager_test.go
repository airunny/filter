package filter

import (
	"context"
	"testing"

	"github.com/liyanbing/filter/cache"
	"github.com/liyanbing/filter/variables"
	"github.com/stretchr/testify/assert"
)

var filterConf = `{
    "filters": {
        "1": {
            "filter_data": [
                [
                    "order_num",
                    "=",
                    10
                ],
                [
                    [
                        "__set_exp_id",
                        "=",
                        "uid"
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

func TestNewFilter(t *testing.T) {
	variables.RegisterVariableFunc("order_num", func(ctx context.Context, data interface{}, cache *cache.Cache) interface{} {
		return 10
	})
	filterManager, err := NewFilter(context.Background(), filterConf, nil)
	assert.Equal(t, nil, err)
	ret, err := filterManager.Execute(context.Background(), nil)
	assert.Equal(t, nil, err)
	t.Logf("result: %#v", ret)
}
