package routers

import (
	"beegoApi/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	//用户
	beego.Include(&controllers.UserController{})
	//视频
	beego.Include(&controllers.VideoController{})

	beego.Include(&controllers.BaseController{})
}
