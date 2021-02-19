package routers

import (
	"beegoApi/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Include(&controllers.BaseController{})
	//用户
	beego.Include(&controllers.UserController{})
	//视频
	beego.Include(&controllers.VideoController{})
	//评论
	beego.Include(&controllers.CommentController{})

}
