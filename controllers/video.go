package controllers

import (
	"beegoApi/models"
	"github.com/astaxie/beego"
)

//视频
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

//获取热播列表
// @router /channel/hot [*]
func (this *VideoController) ChannelHostList() {
	channelId, _ := this.GetInt("channelId")
	if 0 == channelId {
		this.Data["json"] = ReturnError(4001, "必须指定频道")
		this.ServeJSON()
	}
	num, videos, err := models.ChannelHostList(channelId)
	if err == nil {
		this.Data["json"] = ReturnSuccess(0, "success", videos, num)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(4004, "没有相关内容")
		this.ServeJSON()
	}
}

//根据频道地区获取推荐的视频列表
// @router /channel/recommend/region [*]
func (this *VideoController) ChannelRecommendRegionList() {
	channelId, _ := this.GetInt("channelId")
	regionId, _ := this.GetInt("regionId")
	if 0 == channelId {
		this.Data["json"] = ReturnError(4001, "必须指定频道")
		this.ServeJSON()
	}
	if 0 == regionId {
		this.Data["json"] = ReturnError(4001, "必须指定地区")
		this.ServeJSON()
	}
	num, videos, err := models.GetChannelRecommendRegionList(channelId, regionId)
	if err == nil {
		this.Data["json"] = ReturnSuccess(0, "success", videos, num)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(4004, "没有相关内容")
		this.ServeJSON()
	}
}

//根据频道类型获取推荐视频
// @router /channel/recommend/type [*]
func (this *VideoController) GetChannelRecommendTypeList() {
	channelId, _ := this.GetInt("channelId")
	typeId, _ := this.GetInt("typeId")
	if 0 == channelId {
		this.Data["json"] = ReturnError(4001, "必须指定频道")
		this.ServeJSON()
	}
	if 0 == typeId {
		this.Data["json"] = ReturnError(4002, "必须指定频道类型")
		this.ServeJSON()
	}
	num, videos, err := models.GetChannelRecommendTypeList(channelId, typeId)
	if err == nil {
		this.Data["json"] = ReturnSuccess(0, "success", videos, num)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(4004, "没有相关内容")
		this.ServeJSON()
	}

}

//根据传入参数获取视频列表
// @router /channel/video [*]
func (this *VideoController) ChannelVideo() {
	channelId, _ := this.GetInt("channelId")
	//地区
	regionId, _ := this.GetInt("regionId")
	//类型
	typeId, _ := this.GetInt("typeId")
	//状态
	end := this.GetString("end")
	//排序
	sort := this.GetString("sort")
	//页码信息
	limit, _ := this.GetInt("limit")
	offset, _ := this.GetInt("offset")

	if 0 == channelId {
		this.Data["json"] = ReturnError(4001, "必须指定频道")
		this.ServeJSON()
	}
	if 0 == limit {
		limit = 12
	}
	num, videos, err := models.GetChannelVideoList(channelId, regionId, typeId, end, sort, limit, offset)
	if err == nil {
		this.Data["json"] = ReturnSuccess(0, "success", videos, num)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(4004, "没有相关内容")
		this.ServeJSON()
	}

}

//获取视频详情
// @router /video/info [*]
func (this *VideoController) VideoInfo() {
	videoId, _ := this.GetInt("videoId")
	if 0 == videoId {
		this.Data["json"] = ReturnError(4001, "请指定视频")
		this.ServeJSON()
	}
	video, err := models.GetVideoInfo(videoId)
	if err == nil {
		this.Data["json"] = ReturnSuccess(0, "success", video, 1)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(4004, "请求数据失败，请稍后重试~")
		this.ServeJSON()
	}

}

//获取视频聚集列表
// @router /video/episodes/list [*]
func (this *VideoController) VideoEpisodesList() {
	videoId, _ := this.GetInt("videoId")
	if videoId == 0 {
		this.Data["json"] = ReturnError(4001, "必须指定视频ID")
		this.ServeJSON()
	}
	num, episodes, err := models.GetVideoEpisodesList(videoId)
	if err == nil {
		this.Data["json"] = ReturnSuccess(0, "success", episodes, num)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(4004, "请求数据失败，请稍后重试~")
		this.ServeJSON()
	}
}
