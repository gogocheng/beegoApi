package controllers

import (
	"beegoApi/models"
	"github.com/astaxie/beego"
)

//评论
type CommentController struct {
	beego.Controller
}

type CommentInfo struct {
	Id           int             `json:"id"`
	Content      string          `json:"content"`
	AddTime      int64           `json:"addTime"`
	AddTimeTitle string          `json:"addTimeTitle"`
	UserId       int             `json:"userId"`
	Stamp        int             `json:"stamp"`
	PraiseCount  int             `json:"praiseCount"`
	UserInfo     models.UserInfo `json:"userinfo"`
	EpisodesId   int             `json:"episodesId"`
}

//获取评论列表
// @router /comment/list [*]
func (this *CommentController) List() {
	//获取剧集数
	episodesId, _ := this.GetInt("episodesId")
	//获取页码信息
	limit, _ := this.GetInt("limit")
	offset, _ := this.GetInt("offset")
	if 0 == episodesId {
		this.Data["json"] = ReturnError(4001, "必须指定视频剧集")
		this.ServeJSON()
	}
	if 0 == limit {
		limit = 12
	}
	num, comments, err := models.GetCommentList(episodesId, offset, limit)
	if err == nil {
		var data []CommentInfo
		var commentInfo CommentInfo
		for _, v := range comments {
			commentInfo.Id = v.Id
			commentInfo.Content = v.Content
			commentInfo.AddTime = v.AddTime
			commentInfo.AddTimeTitle = DateFormat(v.AddTime)
			commentInfo.UserId = v.UserId
			commentInfo.Stamp = v.Stamp
			commentInfo.PraiseCount = v.PraiseCount
			//用户信息
			commentInfo.UserInfo, _ = models.RedisGetUserInfo(v.UserId)
			data = append(data, commentInfo)
		}
		this.Data["json"] = ReturnSuccess(0, "success", data, num)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(4004, "没有相关内容")
		this.ServeJSON()
	}
}

//发表 评论
// @router /comment/save [*]
func (this *CommentController) Save() {
	content := this.GetString("content")
	uid, _ := this.GetInt("uid")
	//剧集id
	episodesId, _ := this.GetInt("episodesId")
	videoId, _ := this.GetInt("videoId")

	if content == "" {
		this.Data["json"] = ReturnError(4001, "内容不能为空")
		this.ServeJSON()
	}
	if 0 == uid {
		this.Data["json"] = ReturnError(4002, "请先登录")
		this.ServeJSON()
	}
	if 0 == episodesId {
		this.Data["json"] = ReturnError(4003, "必须指定剧集")
		this.ServeJSON()
	}
	if 0 == videoId {
		this.Data["json"] = ReturnError(4004, "必须指定视频")
		this.ServeJSON()
	}
	err := models.SaveComment(content, uid, episodesId, videoId)
	if err == nil {
		this.Data["json"] = ReturnSuccess(0, "success", "", 1)
		this.ServeJSON()
	} else {
		this.Data["json"] = ReturnError(5000, err)
		this.ServeJSON()
	}

}
