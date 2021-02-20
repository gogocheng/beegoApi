package controllers

import (
	"beegoApi/models"
	"github.com/astaxie/beego"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

//用户相关
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

//用户登录
// @router /login/do [*]
func (this *UserController) LoginDo() {
	mobile := this.GetString("mobile")
	password := this.GetString("password")
	if "" == mobile {
		this.Data["json"] = ReturnError(4001, "手机号不能为空")
		this.ServeJSON()
	}
	if "" == password {
		this.Data["json"] = ReturnError(4001, "密码不能为空")
		this.ServeJSON()
	}
	uid, name := models.IsMobileLoginIn(mobile, MD5V(password))
	if uid != 0 {
		m := make(map[string]interface{})
		m["uid"] = uid
		m["username"] = name
		this.Data["json"] = ReturnSuccess(0, "登录成功", m, 1)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(4004, "手机号或密码不正确")
		this.ServeJSON()
	}

}

//批量发送消息
// @router /send/message [*]
func (this *UserController) SendMessageDo() {
	uids := this.GetString("uids")
	content := this.GetString("content")
	if "" == uids {
		this.Data["json"] = ReturnError(4001, "请填写接收用户")
		this.ServeJSON()
	}
	if "" == content {
		this.Data["json"] = ReturnError(4002, "内容不能为空")
		this.ServeJSON()
	}
	messageId, err := models.SendMessageDo(content)
	if err == nil {
		uidConfig := strings.Split(uids, ",")
		for _, v := range uidConfig {
			userId, _ := strconv.Atoi(v)
			models.SendMessageUser(userId, messageId)
		}
		this.Data["json"] = ReturnSuccess(0, "success", "", 1)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(5000, "发送失败，请联系客服")
		this.ServeJSON()
	}
}

//上传视频文件
// @router user/upload [post]
func (this *UserController) UploadVideo() {
	uid := this.GetString("uid")
	if "" == uid {
		this.Data["json"] = ReturnError(4001, "请先登录")
		this.ServeJSON()
	}
	//文件流
	file, header, _ := this.GetFile("file")
	//转换为二进制
	b, _ := ioutil.ReadAll(file)
	//生成文件名
	filename := strings.Split(header.Filename, ".")
	filename[0] = GetVideoName(uid)
	////保存路径
	var fileDir = "./static/upload/" + filename[0] + "." + filename[1]
	log.Println(fileDir)
	//播放地址
	//var playUrl = fileDir
	err := ioutil.WriteFile(fileDir, b, 0777)
	if err == nil {
		m := make(map[string]interface{})
		m["play_url"] = fileDir
		this.Data["json"] = ReturnSuccess(0, "success", m, 1)
		this.ServeJSON()
	} else {
		log.Println(err)
		this.Data["json"] = ReturnError(4004, "上传失败")
		this.ServeJSON()
	}
}
