package models

import (
	"github.com/astaxie/beego/orm"
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