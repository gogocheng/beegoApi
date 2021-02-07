package controllers

import (
	"beegoApi/models"
	"github.com/astaxie/beego"
)

type VideoController struct {
	beego.Controller
}

//获取频道页顶部广告
// @router /channel/advert [*]
func (this *VideoController) ChannelAdvert() {
	channelId, _ := this.GetInt("channelId")
	//log.Println(channelId)
	if 0 == channelId {
		this.Data["json"] = ReturnError(4001, "必须指定频道")
	}
	num, video, err := models.GetChannelAdvert(channelId)
	if err == nil {
		this.Data["json"] = ReturnSuccess(0, "success", video, num)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(4004, "请求数据失败，请稍后重试")
		this.ServeJSON()
	}
}
