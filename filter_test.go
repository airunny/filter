package filter

import (
	"context"
	"errors"
	"testing"

	"github.com/liyanbing/filter/cache"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	cases := []struct {
		JsonStr      string
		Reporter     Reporter
		BuildErr     error
		Data         interface{}
		ExpectedData interface{}
		ResultErr    error
	}{
		// err
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": []
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {}),
			BuildErr: errors.New("filter must contain at least two items"),
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				["success","=",1]	
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {}),
			BuildErr: errors.New("filter must contain at least two items"),
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				"1",
				"2"
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {}),
			BuildErr: errors.New("condition item must contains three element"),
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				["success","=",1],
				"2"
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {}),
			BuildErr: errors.New("executor item must contains 3 elements"),
		},
		// success
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				["success","=",1],
				["name","=","李四"]
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {
				assert.Equal(t, []string{"1"}, filterIds)
			}),
			Data: nil,
			ExpectedData: map[string]interface{}{
				"name": "李四",
			},
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				["success","=",1],
				["name","=","李四"]
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {
				assert.Equal(t, []string{"1"}, filterIds)
			}),
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "李四",
			},
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				["success","=",1],
				["timestamp",">",1],
				[
					["name","=","李四"]
				]
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {
				assert.Equal(t, []string{"1"}, filterIds)
			}),
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "李四",
			},
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				["success","=",1],
				["timestamp",">",1],
				[
					["name","=","李四"],
					["age","=",10]
				]
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {
				assert.Equal(t, []string{"1"}, filterIds)
			}),
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "李四",
				"age":  float64(10), // number is float64 of json unmarshal
			},
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				[
					"and",
					"=>",
					[
						["success","=",1],
						["timestamp",">",1]	
					]
				],
				[
					["name","=","李四"],
					["age","=",10]
				]
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {
				assert.Equal(t, []string{"1"}, filterIds)
			}),
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "李四",
				"age":  float64(10), // number is float64 of json unmarshal
			},
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				[
					"or",
					"=>",
					[
						["success","=",1],
						["timestamp",">",1]	
					]
				],
				[
					["name","=","李四"],
					["age","=",10]
				]
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {
				assert.Equal(t, []string{"1"}, filterIds)
			}),
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "李四",
				"age":  float64(10), // number is float64 of json unmarshal
			},
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				[
					"not",
					"=>",
					[
						["success","<",1],
						["timestamp","<",1]	
					]
				],
				[
					["name","=","李四"],
					["age","=",10]
				]
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {
				assert.Equal(t, []string{"1"}, filterIds)
			}),
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "李四",
				"age":  float64(10), // number is float64 of json unmarshal
			},
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				[
					"or",
					"=>",
					[
						[
							"and",
							"=>",
							[
								["success","<",1],
								["timestamp","<",1]			
							]
						],
						["success","=",1]
					]
				],
				[
					["name","=","李四"],
					["age","=",10]
				]
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {
				assert.Equal(t, []string{"1"}, filterIds)
			}),
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "李四",
				"age":  float64(10), // number is float64 of json unmarshal
			},
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				[
					"and",
					"=>",
					[
						[
							"and",
							"=>",
							[
								["success","<",1],
								["timestamp","<",1]			
							]
						],
						["success","=",1]
					]
				],
				[
					["name","=","李四"],
					["age","=",10]
				]
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {
				assert.Equal(t, []string(nil), filterIds)
			}),
			Data:         map[string]interface{}{},
			ExpectedData: map[string]interface{}{},
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				[
					"not",
					"=>",
					[
						[
							"and",
							"=>",
							[
								["success","<",1],
								["timestamp","<",1]			
							]
						],
						["success","=",1]
					]
				],
				[
					["name","=","李四"],
					["age","=",10]
				]
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {
				assert.Equal(t, []string(nil), filterIds)
			}),
			Data:         map[string]interface{}{},
			ExpectedData: map[string]interface{}{},
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				[
					"not",
					"=>",
					[
						[
							"and",
							"=>",
							[
								["success","<",1],
								["timestamp","<",1]			
							]
						],
						["success","<",1]
					]
				],
				[
					["name","=","李四"],
					["age","=",10]
				]
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {
				assert.Equal(t, []string{"1"}, filterIds)
			}),
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "李四",
				"age":  float64(10),
			},
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				[
					"not",
					"=>",
					[
						[
							"and",
							"=>",
							[
								["success","<",1],
								["timestamp","<",1]			
							]
						],
						["success","<",1]
					]
				],
				[
					["name","=","李四"],
					["age","=",10]
				]
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {
				assert.Equal(t, []string{"1"}, filterIds)
			}),
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "李四",
				"age":  float64(10),
			},
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				["success","=",1],
				["timestamp",">",1],
				[
					["name","=","李四"],
					["age","=",10]
				]
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {
				assert.Equal(t, []string{"1"}, filterIds)
			}),
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "李四",
				"age":  float64(10),
			},
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				["success",">",1],
				[
					["name","=","李四"],
					["age","=",10]
				]
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {
				assert.Equal(t, []string(nil), filterIds)
			}),
			Data:         map[string]interface{}{},
			ExpectedData: map[string]interface{}{},
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				["success","=",1],
				["timestamp", "<=", 1],
				[
					["name","=","李四"],
					["age","=",10]
				]
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {
				assert.Equal(t, []string(nil), filterIds)
			}),
			Data:         map[string]interface{}{},
			ExpectedData: map[string]interface{}{},
		},
		// group
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				["success","=",1],
				["timestamp", "<=", 1],
				[	
					["age","=",10]
				]
			]
		},
		{
			"id":"2",
			"weight": 2,
			"priority": 2,
			"filter": [
				["success","=",1],
				["timestamp", ">", 1],
				[	
					["age","=",10],
					["name","=","李四"]
				]
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {
				assert.Equal(t, []string{"2"}, filterIds)
			}),
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "李四",
				"age":  float64(10),
			},
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				["success","=",1],
				["timestamp", ">", 1],
				[	
					["age","=",10]
				]
			]
		},
		{
			"id":"2",
			"weight": 2,
			"priority": 2,
			"filter": [
				["success","=",1],
				["timestamp", ">", 1],
				[	
					["age","=",10],
					["name","=","李四"]
				]
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {
				assert.Equal(t, []string{"1"}, filterIds)
			}),
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"age": float64(10),
			},
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 2,
			"filter": [
				["success","=",1],
				["timestamp", ">", 1],
				[	
					["age","=",10]
				]
			]
		},
		{
			"id":"2",
			"weight": 2,
			"priority": 1,
			"filter": [
				["success","=",1],
				["timestamp", ">", 1],
				[	
					["age","=",10],
					["name","=","李四"]
				]
			]
		}	
	],
	"batch":false
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {
				assert.Equal(t, []string{"2"}, filterIds)
			}),
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"age":  float64(10),
				"name": "李四",
			},
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 1,
			"filter": [
				["success","=",1],
				["timestamp", ">", 1],
				[	
					["age","=",10]
				]
			]
		},
		{
			"id":"2",
			"weight": 2,
			"priority": 2,
			"filter": [
				["success","=",1],
				["timestamp", ">", 1],
				[	
					["age","=",20],
					["name","=","李四"]
				]
			]
		}	
	],
	"batch": true
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {
				assert.Equal(t, []string{"1", "2"}, filterIds)
			}),
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"age":  float64(20),
				"name": "李四",
			},
		},
		{
			JsonStr: `
{
	"filters":[
		{
			"id":"1",
			"weight": 1,
			"priority": 2,
			"filter": [
				["success","=",1],
				["timestamp", ">", 1],
				[	
					["age","=",10]
				]
			]
		},
		{
			"id":"2",
			"weight": 2,
			"priority": 1,
			"filter": [
				["success","=",1],
				["timestamp", ">", 1],
				[	
					["age","=",20],
					["name","=","李四"]
				]
			]
		}	
	],
	"batch": true
}`,
			Reporter: ReportFunc(func(ctx context.Context, data interface{}, filterIds []string) {
				assert.Equal(t, []string{"2", "1"}, filterIds)
			}),
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"age":  float64(10),
				"name": "李四",
			},
		},
	}

	ctx := context.Background()
	for _, tt := range cases {
		filter, err := NewFilter(ctx, tt.JsonStr, tt.Reporter)
		if err != nil {
			assert.Equal(t, tt.BuildErr, err)
		} else {
			assert.NotNil(t, filter)
			result, resultErr := filter.Execute(ctx, tt.Data)
			assert.Equal(t, tt.ResultErr, resultErr)
			assert.Equal(t, tt.ExpectedData, result)
		}
	}
}

func TestBuildBatchFilter(t *testing.T) {
	cases := []struct {
		conf          *Config
		BuildErr      error
		Data          interface{}
		ExpectedData  interface{}
		SuccessNumber int
		FilterIds     []string
		ResultErr     error
	}{
		// err
		{
			conf: &Config{
				Filters: []struct {
					Id       string        `json:"id"`
					Weight   int64         `json:"weight"`
					Priority int64         `json:"priority"`
					Filter   []interface{} `json:"Filter"`
				}{
					{
						Id:       "1",
						Weight:   1,
						Priority: 2,
						Filter:   []interface{}{},
					},
				},
			},
			BuildErr: errors.New("filter must contain at least two items"),
		},
		{
			conf: &Config{
				Filters: []struct {
					Id       string        `json:"id"`
					Weight   int64         `json:"weight"`
					Priority int64         `json:"priority"`
					Filter   []interface{} `json:"Filter"`
				}{
					{
						Id:       "1",
						Weight:   1,
						Priority: 2,
						Filter:   []interface{}{"1"},
					},
				},
			},
			BuildErr: errors.New("filter must contain at least two items"),
		},
		{
			conf: &Config{
				Filters: []struct {
					Id       string        `json:"id"`
					Weight   int64         `json:"weight"`
					Priority int64         `json:"priority"`
					Filter   []interface{} `json:"Filter"`
				}{
					{
						Id:       "1",
						Weight:   1,
						Priority: 2,
						Filter: []interface{}{
							// condition
							"1",
							// executor
							[]interface{}{"name", "=", "李四"},
						},
					},
				},
			},
			BuildErr: errors.New("condition item must contains three element"),
		},
		{
			conf: &Config{
				Filters: []struct {
					Id       string        `json:"id"`
					Weight   int64         `json:"weight"`
					Priority int64         `json:"priority"`
					Filter   []interface{} `json:"Filter"`
				}{
					{
						Id:       "1",
						Weight:   1,
						Priority: 2,
						Filter: []interface{}{
							// condition
							[]interface{}{"success", "=", 1},
							// executor
							"2",
						},
					},
				},
			},
			BuildErr: errors.New("executor item must contains 3 elements"),
		},
		// success
		// single
		{
			conf: &Config{
				Filters: []struct {
					Id       string        `json:"id"`
					Weight   int64         `json:"weight"`
					Priority int64         `json:"priority"`
					Filter   []interface{} `json:"Filter"`
				}{
					{
						Id:       "1",
						Weight:   1,
						Priority: 2,
						Filter: []interface{}{
							// condition
							[]interface{}{"success", "=", 1},
							// executor
							[]interface{}{"name", "=", "李四"},
						},
					},
				},
			},
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "李四",
			},
			SuccessNumber: 1,
			FilterIds:     []string{"1"},
		},
		{
			conf: &Config{
				Filters: []struct {
					Id       string        `json:"id"`
					Weight   int64         `json:"weight"`
					Priority int64         `json:"priority"`
					Filter   []interface{} `json:"Filter"`
				}{
					{
						Id:       "1",
						Weight:   1,
						Priority: 2,
						Filter: []interface{}{
							// condition
							[]interface{}{"success", "=", 1},
							[]interface{}{"timestamp", ">", 1},
							// executor
							[]interface{}{
								[]interface{}{"name", "=", "李四"},
								[]interface{}{"age", "=", 10},
							},
						},
					},
				},
			},
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "李四",
				"age":  10,
			},
			SuccessNumber: 1,
			FilterIds:     []string{"1"},
		},
		{
			conf: &Config{
				Filters: []struct {
					Id       string        `json:"id"`
					Weight   int64         `json:"weight"`
					Priority int64         `json:"priority"`
					Filter   []interface{} `json:"Filter"`
				}{
					{
						Id:       "1",
						Weight:   1,
						Priority: 2,
						Filter: []interface{}{
							// condition
							[]interface{}{"success", ">", 1},
							// executor
							[]interface{}{"name", "=", "李四"},
						},
					},
				},
			},
			Data:          map[string]interface{}{},
			ExpectedData:  map[string]interface{}{},
			SuccessNumber: 0,
		},
		{
			conf: &Config{
				Filters: []struct {
					Id       string        `json:"id"`
					Weight   int64         `json:"weight"`
					Priority int64         `json:"priority"`
					Filter   []interface{} `json:"Filter"`
				}{
					{
						Id:       "1",
						Weight:   1,
						Priority: 2,
						Filter: []interface{}{
							// condition
							[]interface{}{"success", "=", 1},
							[]interface{}{"timestamp", "<=", 1},
							// executor
							[]interface{}{
								[]interface{}{"age", "=", 10},
							},
						},
					},
				},
			},
			Data:          map[string]interface{}{},
			ExpectedData:  map[string]interface{}{},
			SuccessNumber: 0,
		},
		// batch
		{
			conf: &Config{
				Filters: []struct {
					Id       string        `json:"id"`
					Weight   int64         `json:"weight"`
					Priority int64         `json:"priority"`
					Filter   []interface{} `json:"Filter"`
				}{
					{
						Id:       "1",
						Weight:   1,
						Priority: 2,
						Filter: []interface{}{
							// condition
							[]interface{}{"success", "=", 1},
							[]interface{}{"timestamp", "<=", 1},
							// executor
							[]interface{}{
								[]interface{}{"age", "=", 10},
							},
						},
					},
					{
						Id:       "2",
						Weight:   2,
						Priority: 2,
						Filter: []interface{}{
							// condition
							[]interface{}{"success", "=", 1},
							[]interface{}{"timestamp", ">", 1},
							// executor
							[]interface{}{
								[]interface{}{"name", "=", "李四"},
								[]interface{}{"age", "=", 10},
							},
						},
					},
				},
			},
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "李四",
				"age":  10,
			},
			SuccessNumber: 1,
			FilterIds:     []string{"2"},
		},
		{
			conf: &Config{
				Filters: []struct {
					Id       string        `json:"id"`
					Weight   int64         `json:"weight"`
					Priority int64         `json:"priority"`
					Filter   []interface{} `json:"Filter"`
				}{
					{
						Id:       "1",
						Weight:   1,
						Priority: 1,
						Filter: []interface{}{
							// condition
							[]interface{}{"success", "=", 1},
							[]interface{}{"timestamp", ">", 1},
							// executor
							[]interface{}{
								[]interface{}{"age", "=", 10},
							},
						},
					},
					{
						Id:       "2",
						Weight:   2,
						Priority: 2,
						Filter: []interface{}{
							// condition
							[]interface{}{"success", "=", 1},
							[]interface{}{"timestamp", ">", 1},
							// executor
							[]interface{}{
								[]interface{}{"name", "=", "李四"},
								[]interface{}{"age", "=", 20},
							},
						},
					},
				},
			},
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"age": 10,
			},
			SuccessNumber: 1,
			FilterIds:     []string{"1"},
		},
		{
			conf: &Config{
				Filters: []struct {
					Id       string        `json:"id"`
					Weight   int64         `json:"weight"`
					Priority int64         `json:"priority"`
					Filter   []interface{} `json:"Filter"`
				}{
					{
						Id:       "1",
						Weight:   1,
						Priority: 2,
						Filter: []interface{}{
							// condition
							[]interface{}{"success", "=", 1},
							[]interface{}{"timestamp", ">", 1},
							// executor
							[]interface{}{
								[]interface{}{"age", "=", 10},
							},
						},
					},
					{
						Id:       "2",
						Weight:   2,
						Priority: 1,
						Filter: []interface{}{
							// condition
							[]interface{}{"success", "=", 1},
							[]interface{}{"timestamp", ">", 1},
							// executor
							[]interface{}{
								[]interface{}{"name", "=", "李四"},
								[]interface{}{"age", "=", 20},
							},
						},
					},
				},
			},
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "李四",
				"age":  20,
			},
			SuccessNumber: 1,
			FilterIds:     []string{"2"},
		},
		{
			conf: &Config{
				Filters: []struct {
					Id       string        `json:"id"`
					Weight   int64         `json:"weight"`
					Priority int64         `json:"priority"`
					Filter   []interface{} `json:"Filter"`
				}{
					{
						Id:       "1",
						Weight:   1,
						Priority: 1,
						Filter: []interface{}{
							// condition
							[]interface{}{"success", "=", 1},
							[]interface{}{"timestamp", ">", 1},
							// executor
							[]interface{}{
								[]interface{}{"age", "=", 10},
							},
						},
					},
					{
						Id:       "2",
						Weight:   2,
						Priority: 2,
						Filter: []interface{}{
							// condition
							[]interface{}{"success", "=", 1},
							[]interface{}{"timestamp", ">", 1},
							// executor
							[]interface{}{
								[]interface{}{"name", "=", "李四"},
								[]interface{}{"age", "=", 20},
							},
						},
					},
				},
				Batch: true,
			},
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "李四",
				"age":  20,
			},
			SuccessNumber: 2,
			FilterIds:     []string{"1", "2"},
		},
		{
			conf: &Config{
				Filters: []struct {
					Id       string        `json:"id"`
					Weight   int64         `json:"weight"`
					Priority int64         `json:"priority"`
					Filter   []interface{} `json:"Filter"`
				}{
					{
						Id:       "1",
						Weight:   1,
						Priority: 2,
						Filter: []interface{}{
							// condition
							[]interface{}{"success", "=", 1},
							[]interface{}{"timestamp", ">", 1},
							// executor
							[]interface{}{
								[]interface{}{"age", "=", 10},
							},
						},
					},
					{
						Id:       "2",
						Weight:   2,
						Priority: 1,
						Filter: []interface{}{
							// condition
							[]interface{}{"success", "=", 1},
							[]interface{}{"timestamp", ">", 1},
							// executor
							[]interface{}{
								[]interface{}{"name", "=", "李四"},
								[]interface{}{"age", "=", 20},
							},
						},
					},
				},
				Batch: true,
			},
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "李四",
				"age":  10,
			},
			SuccessNumber: 2,
			FilterIds:     []string{"2", "1"},
		},
	}

	ctx := context.Background()
	for _, tt := range cases {
		batch, err := buildBatchFilter(ctx, tt.conf)
		if err != nil {
			assert.Equal(t, tt.BuildErr, err)
		} else {
			assert.NotNil(t, batch)
			successNumber, filterIds, resultErr := batch.Run(ctx, tt.Data, cache.NewCache())
			assert.Equal(t, tt.ResultErr, resultErr)
			assert.Equal(t, tt.SuccessNumber, successNumber)
			assert.Equal(t, tt.FilterIds, filterIds)
			assert.Equal(t, tt.ExpectedData, tt.Data)
		}
	}
}

func TestBuildSingleFilter(t *testing.T) {
	cases := []struct {
		Id           string
		Weight       int64
		Priority     int64
		FilterData   []interface{}
		BuildErr     error
		Data         interface{}
		ExpectedData interface{}
		Result       bool
		ResultErr    error
	}{
		// err
		{
			Id:         "1",
			Weight:     1,
			Priority:   1,
			FilterData: []interface{}{},
			BuildErr:   errors.New("filter must contain at least two items"),
		},
		{
			Id:         "2",
			Weight:     2,
			Priority:   2,
			FilterData: []interface{}{"1"},
			BuildErr:   errors.New("filter must contain at least two items"),
		},
		{
			Id:       "2",
			Weight:   2,
			Priority: 2,
			FilterData: []interface{}{
				// condition
				"1",
				// executor
				[]interface{}{"name", "=", "李四"},
			},
			BuildErr: errors.New("condition item must contains three element"),
		},
		{
			Id:       "2",
			Weight:   2,
			Priority: 2,
			FilterData: []interface{}{
				// condition
				[]interface{}{"success", "=", 1},
				// executor
				"2",
			},
			BuildErr: errors.New("executor item must contains 3 elements"),
		},
		// success
		{
			Id:       "3",
			Weight:   3,
			Priority: 3,
			FilterData: []interface{}{
				// condition
				[]interface{}{"success", "=", 1},
				// executor
				[]interface{}{"name", "=", "李四"},
			},
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "李四",
			},
			Result: true,
		},
		{
			Id:       "4",
			Weight:   4,
			Priority: 4,
			FilterData: []interface{}{
				// condition
				[]interface{}{"success", "=", 1},
				[]interface{}{"timestamp", ">", 1},
				// executor
				[]interface{}{
					[]interface{}{"name", "=", "李四"},
					[]interface{}{"age", "=", 10},
				},
			},
			Data: map[string]interface{}{},
			ExpectedData: map[string]interface{}{
				"name": "李四",
				"age":  10,
			},
			Result: true,
		},
		{
			Id:       "5",
			Weight:   5,
			Priority: 5,
			FilterData: []interface{}{
				// condition
				[]interface{}{"success", ">", 1},
				// executor
				[]interface{}{"name", "=", "李四"},
			},
			Data:         map[string]interface{}{},
			ExpectedData: map[string]interface{}{},
			Result:       false,
		},
		{
			Id:       "4",
			Weight:   4,
			Priority: 4,
			FilterData: []interface{}{
				// condition
				[]interface{}{"success", "=", 1},
				[]interface{}{"timestamp", "<=", 1},
				// executor
				[]interface{}{
					[]interface{}{"age", "=", 10},
				},
			},
			Data:         map[string]interface{}{},
			ExpectedData: map[string]interface{}{},
			Result:       false,
		},
	}

	ctx := context.Background()
	for _, tt := range cases {
		fil, err := buildSingleFilter(ctx, tt.Id, tt.Weight, tt.Priority, tt.FilterData)
		if err != nil {
			assert.Equal(t, tt.BuildErr, err)
		} else {
			assert.NotNil(t, fil)
			result, resultErr := fil.Run(ctx, tt.Data, cache.NewCache())
			assert.Equal(t, tt.ResultErr, resultErr)
			assert.Equal(t, tt.Result, result)
			assert.Equal(t, tt.ExpectedData, tt.Data)
		}
	}
}

func TestTotalWeight(t *testing.T) {
	var (
		filters = make([]*singleFilter, 0, 10)
		total   = int64(0)
	)

	for i := 0; i < 10; i++ {
		total += int64(i)
		filters = append(filters, &singleFilter{
			weight: int64(i),
		})
	}
	assert.Equal(t, total, filterWeight(filters))
}

func TestPickByWeight(t *testing.T) {
	var (
		filters = make([]*singleFilter, 0, 10)
		total   = int64(0)
	)

	for i := 1; i <= 10; i++ {
		total += int64(i)
		filters = append(filters, &singleFilter{
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

func TestShuffleByWeight(t *testing.T) {
	var (
		filters     = make([]*singleFilter, 0, 10)
		totalWeight = int64(0)
	)

	for i := 1; i <= 10; i++ {
		totalWeight += int64(i)
		filters = append(filters, &singleFilter{
			weight: int64(i),
		})
	}

	total := len(filters)
	for i := 0; i < 100000; i++ {
		shuffleByWeight(filters, totalWeight)
		assert.Equal(t, total, len(filters))
		countMapping := make(map[*singleFilter]struct{})
		for _, f := range filters {
			_, ok := countMapping[f]
			assert.False(t, ok)
			countMapping[f] = struct{}{}
		}
	}
}
