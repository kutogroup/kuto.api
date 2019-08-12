package pkg

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

//WahaCache 缓存结构体
type WahaCache struct {
	pool *redis.Pool
}

//NewCache 新建缓存对象
func NewCache(address string, poolSize int, timeout time.Duration) *WahaCache {
	return &WahaCache{
		pool: &redis.Pool{
			MaxIdle: poolSize,
			// MaxActive:   5,
			IdleTimeout: timeout,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", address)
				if err != nil {
					return nil, err
				}

				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
		},
	}
}

//Get 获取缓存对象
func (c *WahaCache) Get(key string) (string, error) {
	conn := c.pool.Get()
	conn.Send("GET", key)
	conn.Flush()
	v, err := conn.Receive()
	if err != nil {
		return "", err
	}

	return redis.String(v, err)
}

//Set 设置缓存对象
func (c *WahaCache) Set(key string, value string) error {
	conn := c.pool.Get()
	conn.Send("SET", key, value)
	conn.Flush()
	_, err := conn.Receive()
	return err
}

//SetByEx 设置缓存对象
func (c *WahaCache) SetByEx(key string, value string, expired int) error {
	conn := c.pool.Get()
	conn.Send("SET", key, value, "EX", expired)
	conn.Flush()
	_, err := conn.Receive()
	return err
}

//Exists 判断缓存是否存在
func (c *WahaCache) Exists(key string) bool {
	exists, _ := redis.Bool(c.pool.Get().Do("EXISTS", key))
	return exists
}

//Del 删除缓存
func (c *WahaCache) Del(key string) error {
	conn := c.pool.Get()
	conn.Send("DEL", key)
	conn.Flush()
	_, err := conn.Receive()
	return err
}
