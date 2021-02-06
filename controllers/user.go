package controllers

import (
	"beegoApi/models"
	"github.com/astaxie/beego"
	"regexp"
)

type UserController struct {
	beego.Controller
}

//用户注册
// @router /register/save [post]
func (this *UserController) SaveRegister() {
	var (
		mobile   string
		password string
		err      error
	)
	mobile = this.GetString("mobile")
	password = this.GetString("password")
	if "" == mobile {
		this.Data["json"] = ReturnError(4001, "手机号不能为空")
		this.ServeJSON()
	}
	//手机号格式校验
	isorno, _ := regexp.MatchString(`^1(3|4|5|6|7|8)[0-9]\d{8}$`, mobile)
	if !isorno {
		this.Data["json"] = ReturnError(4002, "手机号格式不正确")
		this.ServeJSON()
	}
	if "" == password {
		this.Data["json"] = ReturnError(4001, "密码为空")
		this.ServeJSON()
	}
	//判断手机号是否注册
	status := models.IsUserMobile(mobile)
	if status {
		//已注册
		this.Data["json"] = ReturnError(4005, "此手机号已经注册")
		this.ServeJSON()
	} else {
		err = models.UserSave(mobile, MD5V(password))
		if err == nil {
			this.Data["json"] = ReturnSuccess(0, "注册成功", nil, 0)
			this.ServeJSON()
		} else {
			this.Data["json"] = ReturnError(5000, err)
		}
	}
}
