package iredis

import (
	"gopkg.in/redis.v5"
)

type RedisCfg struct {
	MasterName    string   `json:"master_name"`
	SentinelAddrs []string `json:"sentinel_addrs"`
	Host          string   `json:"host"`
	Password      string   `json:"password"`
	DB            int      `json:"db"`
	MaxIdle       int      `json:"max_idle"`
}

type Client struct {
	RedisCfg
	*redis.Client
}

func NewClient(cfg RedisCfg) (client *Client, err error) {
	client = &Client{
		RedisCfg: cfg,
	}

	if cfg.MaxIdle <= 0 {
		cfg.MaxIdle = 30
	}

	if len(cfg.SentinelAddrs) != 0 {
		client.Client = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    cfg.MasterName,
			SentinelAddrs: cfg.SentinelAddrs,
			Password:      cfg.Password,

			DB:       cfg.DB,
			PoolSize: cfg.MaxIdle,
		})
	} else {
		client.Client = redis.NewClient(&redis.Options{
			Addr:     cfg.Host,
			Password: cfg.Password,

			DB:       cfg.DB,
			PoolSize: cfg.MaxIdle,
		})
	}

	err = client.Ping().Err()
	return
}
