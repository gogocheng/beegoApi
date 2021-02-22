package controllers

import (
	"beegoApi/models"
	"github.com/astaxie/beego"
)

//排行榜
type TopController struct {
	beego.Controller
}

//根据频道获取排行榜
// @router /channel/top [*]
func (this *TopController) ChannelTop() {
	//频道id
	channelId, _ := this.GetInt("channelId")
	if 0 == channelId {
		this.Data["json"] = ReturnError(4001, "必须指定频道")
		this.ServeJSON()
	}
	num, videos, err := models.RedisGetChannelTop(channelId)
	if err == nil {
		this.Data["json"] = ReturnSuccess(0, "success", videos, num)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(4004, "没有相关内容")
		this.ServeJSON()
	}
}

//根据类型获取排行榜
// @router /type/top [*]
func (this *TopController) TypeTop() {
	typeId, _ := this.GetInt("typeId")
	if 0 == typeId {
		this.Data["json"] = ReturnError(4001, "必须指定类型")
		this.ServeJSON()
	}
	num, videos, err := models.RedisGetTypeTop(typeId)
	if err == nil {
		this.Data["json"] = ReturnSuccess(0, "success", videos, num)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(4004, "没有相关内容")
		this.ServeJSON()
	}
}
