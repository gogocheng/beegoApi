package redisClient

import (
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"time"
)

//直接链接
func Connect() redis.Conn {
	pool, _ := redis.Dial("tcp", beego.AppConfig.String("redisdb"))
	return pool
}

//连接池
func PoolConnect() redis.Conn {
	pool := &redis.Pool{
		MaxIdle:     1,  //最大空闲连接数
		MaxActive:   10, //最大连接数
		IdleTimeout: 100 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", beego.AppConfig.String("redisdb"))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	return pool.Get()
}
