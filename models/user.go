package models

import (
	redisClient "beegoApi/service/redis"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

type User struct {
	Id       int
	Name     string
	Password string
	Status   int
	AddTime  int64
	Mobile   string
	Avatar   string
}

type UserInfo struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	AddTime int64  `json:"addTime"`
	Avatar  string `json:"avatar"`
}

func init() {
	orm.RegisterModel(new(User))
}

//判断手机号是否存在
func IsUserMobile(mobile string) bool {
	o := orm.NewOrm()
	user := User{Mobile: mobile}
	err := o.Read(&user, "mobile")
	if err == orm.ErrNoRows {
		return false
	} else if err == orm.ErrMissPK {
		return false
	}
	return true
}

//保存用户信息
func UserSave(mobile string, password string) error {
	o := orm.NewOrm()
	var (
		user User
	)
	user.Name = ""
	user.Mobile = mobile
	user.Password = password
	user.Status = 0
	user.AddTime = time.Now().Unix()
	_, err := o.Insert(&user)
	return err
}

//检查是否登录
func IsMobileLoginIn(mobile string, password string) (int, string) {
	o := orm.NewOrm()
	var user User
	err := o.QueryTable("user").Filter("mobile", mobile).Filter("password", password).One(&user)
	if err == orm.ErrNoRows {
		return 0, ""
	} else if err == orm.ErrMissPK {
		return 0, ""
	}
	return user.Id, user.Name
}

//获取用户信息
func GetUserInfo(uid int) (UserInfo, error) {
	o := orm.NewOrm()
	var user UserInfo
	err := o.Raw("SELECT id,name,add_time,avatar FROM user WHERE id=? LIMIT 1", uid).QueryRow(&user)
	return user, err
}

//redis 缓存 获取用户信息
func RedisGetUserInfo(uid int) (UserInfo, error) {
	var user UserInfo
	conn := redisClient.PoolConnect()
	defer conn.Close()
	//定义redis key
	redisKey := "vedio:uid:" + strconv.Itoa(uid)
	//判断是否存在
	exists, err := redis.Bool(conn.Do("exists", redisKey))
	if exists {
		res, _ := redis.Values(conn.Do("hgetall", redisKey))
		err = redis.ScanStruct(res, &user)
	} else {
		//mysql中取数据
		o := orm.NewOrm()
		err := o.Raw("SELECT id,name,add_time,avatar FROM user WHERE id=? LIMIT 1", uid).QueryRow(&user)
		if err == nil {
			//保存到redis
			_, err := conn.Do("hmset", redis.Args{redisKey}.AddFlat(user)...)
			if err == nil {
				//设置过期时间
				conn.Do("expire", redisKey, 86400)
			}
		}
	}
	return user, err
}
