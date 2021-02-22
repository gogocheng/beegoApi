package models

import (
	redisClient "beegoApi/service/redis"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	"strconv"
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

//获取视频详情
func GetVideoInfo(videoId int) (Video, error) {
	o := orm.NewOrm()
	var video Video
	err := o.Raw("select * from video where id=? limit 1", videoId).QueryRow(&video)
	return video, err
}

//redis 获取视频详情
func RedisGetVideoInfo(videoId int) (Video, error) {
	var video Video
	conn := redisClient.PoolConnect()
	defer conn.Close()
	//定义redis key
	redisKey := "vedio:id:" + strconv.Itoa(videoId)
	//判断是否存在
	exists, err := redis.Bool(conn.Do("exists", redisKey))
	if exists {
		res, _ := redis.Values(conn.Do("hgetall", redisKey))
		err = redis.ScanStruct(res, &video)
	} else {
		//mysql中取数据
		o := orm.NewOrm()
		err := o.Raw("select * from video where id=? limit 1", videoId).QueryRow(&video)
		if err == nil {
			//保存到redis
			_, err := conn.Do("hmset", redis.Args{redisKey}.AddFlat(video)...)
			if err == nil {
				//设置过期时间
				conn.Do("expire", redisKey, 86400)
			}
		}
	}
	return video, err
}

//获取视频剧集列表
func GetVideoEpisodesList(videoId int) (int64, []Episodes, error) {
	o := orm.NewOrm()
	var episodes []Episodes
	num, err := o.Raw("select id,title,add_time,num,play_url,comment from video_episodes where video_id=? order by num asc", videoId).QueryRows(&episodes)
	return num, episodes, err
}

//redis 缓存  获取视频剧集列表
func RedisVideoEpisodesList(videoId int) (int64, []Episodes, error) {
	var (
		episodes []Episodes
		num      int64
		err      error
	)
	conn := redisClient.PoolConnect()
	defer conn.Close()
	//定义redis key
	redisKey := "vedio:episodes:videoId:" + strconv.Itoa(videoId)
	//判断是否存在
	exists, err := redis.Bool(conn.Do("exists", redisKey))
	if exists {
		//获取list长度
		num, err = redis.Int64(conn.Do("llen", redisKey))
		if err == nil {
			values, _ := redis.Values(conn.Do("lrange", redisKey, "0", "-1"))
			var episodesInfo Episodes
			//遍历获取
			for _, v := range values {
				err = json.Unmarshal(v.([]byte), &episodesInfo)
				if err == nil {
					episodes = append(episodes, episodesInfo)
				}
			}
		}
	} else {
		//mysql中取数据
		o := orm.NewOrm()
		num, err = o.Raw("select id,title,add_time,num,play_url,comment from video_episodes where video_id=? order by num asc", videoId).QueryRows(&episodes)
		if err == nil {
			//保存到redis  转化成json保存到list
			for _, v := range episodes {
				jsonValue, err := json.Marshal(v)
				if err == nil {
					conn.Do("rpush", redisKey, jsonValue)
				}
			}
			//设置过期时间
			conn.Do("expire", redisKey, 86400)
		}
	}
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

//redis 改造频道排行榜
func RedisGetChannelTop(channelId int) (int64, []VideoData, error) {
	var (
		videos []VideoData
		num    int64
	)
	conn := redisClient.PoolConnect()
	defer conn.Close()
	//定义redis key
	redisKey := "vedio:topchannel:channelId:" + strconv.Itoa(channelId)
	//判断是否存在
	exists, err := redis.Bool(conn.Do("exists", redisKey))
	if exists {
		num = 0
		//redis有序集合中获取
		res, _ := redis.Values(conn.Do("zrevarange", redisKey, "0", "10", "WITHSCORES"))
		for k, v := range res {
			fmt.Println(string(v.([]byte)))
			//获取id
			if k%2 == 0 {
				videoId, err := strconv.Atoi(string(v.([]byte)))
				videoInfo, err := RedisGetVideoInfo(videoId)
				if err == nil {
					var videoDataInfo VideoData
					videoDataInfo.Id = videoInfo.Id
					videoDataInfo.Img = videoInfo.Img
					videoDataInfo.Img1 = videoInfo.Img1
					videoDataInfo.IsEnd = videoInfo.IsEnd
					videoDataInfo.SubTitle = videoInfo.SubTitle
					videoDataInfo.Title = videoInfo.Title
					videoDataInfo.AddTime = videoInfo.AddTime
					videoDataInfo.Comment = videoInfo.Comment
					videoDataInfo.EpisodesCount = videoInfo.EpisodesCount
					videos = append(videos, videoDataInfo)
					num++
				}
			}
		}
	} else {
		//mysql中取数据
		o := orm.NewOrm()
		num, err = o.Raw("select id,title,sub_title,img,img1,add_time,episodes_count,is_end from video where status=1"+
			" and channel_id=? order by comment desc limit 10", channelId).QueryRows(&videos)
		if err == nil {
			//保存到redis    用redis有序集合
			for _, v := range videos {
				conn.Do("zadd", redisKey, v.Comment, v.Id)
			}
			//设置过期时间一个月
			conn.Do("expire", redisKey, 86400*30)
		}
	}
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

//redis  改造类型排行榜
func RedisGetTypeTop(typeId int) (int64, []VideoData, error) {
	var (
		videos []VideoData
		num    int64
	)
	conn := redisClient.PoolConnect()
	defer conn.Close()
	//定义redis key
	redisKey := "vedio:toptype:typeId:" + strconv.Itoa(typeId)
	//判断是否存在
	exists, err := redis.Bool(conn.Do("exists", redisKey))
	if exists {
		num = 0
		//redis有序集合中获取
		res, _ := redis.Values(conn.Do("zrevarange", redisKey, "0", "10", "WITHSCORES"))
		fmt.Println(num)
		fmt.Println(22)
		for k, v := range res {
			fmt.Println(string(v.([]byte)))
			fmt.Println(11)
			//获取id
			if k%2 == 0 {
				videoId, _ := strconv.Atoi(string(v.([]byte)))
				videoInfo, err := RedisGetVideoInfo(videoId)
				if err == nil {
					var videoDataInfo VideoData
					videoDataInfo.Id = videoInfo.Id
					videoDataInfo.Title = videoInfo.Title
					videoDataInfo.SubTitle = videoInfo.SubTitle
					videoDataInfo.Img = videoInfo.Img
					videoDataInfo.Img1 = videoInfo.Img1
					videoDataInfo.AddTime = videoInfo.AddTime
					videoDataInfo.EpisodesCount = videoInfo.EpisodesCount
					videoDataInfo.IsEnd = videoInfo.IsEnd
					videos = append(videos, videoDataInfo)
					num++
				}
			}
		}
	} else {
		//mysql中取数据
		o := orm.NewOrm()
		num, err = o.Raw("select id,title,sub_title,img,img1,add_time,episodes_count,is_end from video where status=1"+
			" and channel_id=? order by comment desc limit 10", typeId).QueryRows(&videos)
		if err == nil {
			//保存到redis    用redis有序集合
			for _, v := range videos {
				conn.Do("zadd", redisKey, v.Comment, v.Id)
			}
			//设置过期时间一个月
			conn.Do("expire", redisKey, 86400*30)
		}
	}
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
