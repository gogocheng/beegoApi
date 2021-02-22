package controllers

import (
	redisClient "beegoApi/service/redis"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
)

type RedisDemoController struct {
	beego.Controller
}

// @router /redis/demo [*]
func (this *RedisDemoController) Demo() {
	c := redisClient.PoolConnect()
	defer c.Close()
	//_, err := c.Do("set", "username", "gogochen")
	//if err == nil {
	//	//设置过期时间
	//	c.Do("expire", "username", 1000)
	//}
	r, err := redis.String(c.Do("get", "username"))
	if err == nil {
		fmt.Println(1)
		fmt.Println(r)
		//获取剩余过期时间
		ttl, _ := redis.Int64(c.Do("ttl", "username"))
		fmt.Println(ttl)

	} else {
		fmt.Println(2)
		fmt.Println(err)
	}
}
