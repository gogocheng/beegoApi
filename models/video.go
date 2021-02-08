package models

import "github.com/astaxie/beego/orm"

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
