package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

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

// 剧集
type Episodes struct {
	Id            int
	Title         string
	AddTime       int64
	Num           int
	PlayUrl       string
	Comment       int
	AliyunVideoId string
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

func GetVideoInfo(videoId int) (Video, error) {
	o := orm.NewOrm()
	var video Video
	err := o.Raw("select * from video where id=? limit 1", videoId).QueryRow(&video)
	return video, err
}

//获取视频剧集列表
func GetVideoEpisodesList(videoId int) (int64, []Episodes, error) {
	o := orm.NewOrm()
	var episodes []Episodes
	num, err := o.Raw("select id,title,add_time,num,play_url,comment from video_episodes where video_id=? order by num asc", videoId).QueryRows(&episodes)
	return num, episodes, err
}

//频道排行榜
func GetChannelTop(channelId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var (
		videos []VideoData
	)
	num, err := o.Raw("select id,title,sub_title,img,img1,add_time,episodes_count,is_end from video where status=1"+
		" and channel_id=? order by comment desc limit 10", channelId).QueryRows(&videos)
	return num, videos, err
}

//类型排行榜
func GetTypeTop(typeId int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var (
		videos []VideoData
	)
	num, err := o.Raw("select id,title,sub_title,img,img1,add_time,episodes_count,is_end from video where status=1"+
		" and type_id=? order by comment desc limit 10", typeId).QueryRows(&videos)
	return num, videos, err
}

//获取我的视频
func GetUserVideo(uid int) (int64, []VideoData, error) {
	o := orm.NewOrm()
	var videos []VideoData
	num, err := o.Raw("select id,title,sub_title,img,img1,add_time,episodes_count, is_end from video where user_id=?"+
		" order by add_time desc", uid).QueryRows(&videos)
	return num, videos, err
}

//保存视频信息
func SaveVideo(title string, subTitle string, channelId int, regionId int, typeId int, playUrl string, user_id int, aliyunVideoId string) error {
	o := orm.NewOrm()
	var video Video
	time := time.Now().Unix()
	video.Title = title
	video.SubTitle = subTitle
	video.AddTime = time
	video.Img = ""
	video.Img1 = ""
	video.EpisodesCount = 1
	video.IsEnd = 1
	video.ChannelId = channelId
	video.Status = 1
	video.RegionId = regionId
	video.TypeId = typeId
	video.EpisodesUpdateTime = time
	video.Comment = 0
	video.UserId = user_id
	videoId, err := o.Insert(&video)
	if err == nil {
		if aliyunVideoId != "" {
			playUrl = ""
		}
		_, err = o.Raw("insert into video_episodes (title,add_time,num,video_id,play_url,status,comment,aliyun_video_id) values (?,?,?,?,?,?,?,?)", subTitle, time, 1, videoId, playUrl, 1, 0, aliyunVideoId).Exec()
		//fmt.Println(err)
	}
	return err
}
