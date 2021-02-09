package models

import "github.com/astaxie/beego/orm"

type Video struct {
	Id                 int
	Title              string
	SubTitle           string
	AddTime            int64
	Img                string
	Img1               string
	EpisodesCount      int
	IsEnd              int
	ChannelId          int
	Status             int
	RegionId           int
	TypeId             int
	EpisodesUpdateTime int64
	Comment            int
	UserId             int
	IsRecommend        int
}

type VideoData struct {
	Id            int
	Title         string
	SubTitle      string
	AddTime       int64
	Img           string
	Img1          string
	EpisodesCount int
	IsEnd         int
	Comment       int
}

func init() {
	orm.RegisterModel(new(Video))
}

//获取热播列表
func ChannelHostList(channelId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var (
		videos []VideoData
	)
	num, err := o.Raw("select id,title,sub_title,add_time,img,img1,episodes_count,is_end from video where status=1"+
		" and channel_id=? order by episodes_update_time desc limit 9", channelId).QueryRows(&videos)
	return num, videos, err

}

//获取频道页地区推荐视频列表
func GetChannelRecommendRegionList(channelId int, regionId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var (
		videos []VideoData
	)
	num, err := o.Raw("select id,title,sub_title,add_time,episodes_count,is_end from video where status=1 and channel_id=?"+
		" and is_recommend=1 and region_id=? order by episodes_update_time desc limit 9", channelId, regionId).QueryRows(&videos)
	return num, videos, err
}

//获取频道页不同类型视频
func GetChannelRecommendTypeList(channelId int, typeId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var (
		videos []VideoData
	)
	num, err := o.Raw("select id,title,sub_title,add_time,episodes_count,is_end from video where status=1 and channel_id=?"+
		" and is_recommend=1 and type_id=? order by episodes_update_time desc limit 9", channelId, typeId).QueryRows(&videos)
	return num, videos, err
}

//按条件获取视频列表  orm
func GetChannelVideoList(channelId int, regionId int, typeId int, end string, sort string, limit int, offset int) (int64,
	[]orm.Params, error) {
	o := orm.NewOrm()
	var (
		videos []orm.Params
	)
	qs := o.QueryTable("video")
	qs = qs.Filter("channel_id", channelId)
	qs = qs.Filter("status", 1)
	//判断
	if 0 < regionId {
		qs = qs.Filter("region_id", regionId)
	}
	if 0 < typeId {
		qs = qs.Filter("type_id", typeId)
	}
	if "n" == end {
		qs = qs.Filter("is_end", 0)
	} else if "y" == end {
		qs = qs.Filter("is_end", 1)
	}
	//排序
	if "episodesUpdateTime" == sort {
		//倒序
		qs = qs.OrderBy("-episodes_update_time")
	} else if "comment" == sort {
		qs = qs.OrderBy("-comment")
	} else if "addTime" == sort {
		qs = qs.OrderBy("-add_time")
	} else {
		qs = qs.OrderBy("-add_time")
	}
	nums, _ := qs.Values(&videos, "id", "title", "sub_title", "img", "img1", "episodes_count", "is_end")
	qs = qs.Limit(limit, offset)
	_, err := qs.Values(&videos, "id", "title", "sub_title", "img", "img1", "episodes_count", "is_end")
	return nums, videos, err
}
