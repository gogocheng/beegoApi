package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Comment struct {
	Id          int
	Content     string
	AddTime     int64
	UserId      int
	Stamp       int
	Status      int
	PraiseCount int
	EpisodesId  int
	VideoId     int
}

func init() {
	orm.RegisterModel(new(Comment))
}

//获取评论列表
func GetCommentList(episodesId int, offset int, limit int) (int64, []Comment, error) {
	o := orm.NewOrm()
	var comments []Comment
	num, _ := o.Raw("select id from comment where status=1 and episodes_id=?", episodesId).QueryRows(&comments)
	_, err := o.Raw("select id,content,add_time,user_id,stamp,praise_count,episodes_id from comment where status=1"+
		" and episodes_id=? order by add_time desc limit ?,?", episodesId, offset, limit).QueryRows(&comments)
	return num, comments, err
}

//格式化时间
func DateFormat(times int64) string {
	video_time := time.Unix(times, 0)
	return video_time.Format("2006-01-02")
}
