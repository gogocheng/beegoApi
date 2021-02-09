package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["beegoApi/controllers:BaseController"] = append(beego.GlobalControllerRouter["beegoApi/controllers:BaseController"],
		beego.ControllerComments{
			Method:           "ChannelRegion",
			Router:           "/channel/region",
			AllowHTTPMethods: []string{"*"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beegoApi/controllers:BaseController"] = append(beego.GlobalControllerRouter["beegoApi/controllers:BaseController"],
		beego.ControllerComments{
			Method:           "ChannelType",
			Router:           "/channel/type",
			AllowHTTPMethods: []string{"*"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beegoApi/controllers:UserController"] = append(beego.GlobalControllerRouter["beegoApi/controllers:UserController"],
		beego.ControllerComments{
			Method:           "LoginDo",
			Router:           "/login/do",
			AllowHTTPMethods: []string{"*"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beegoApi/controllers:UserController"] = append(beego.GlobalControllerRouter["beegoApi/controllers:UserController"],
		beego.ControllerComments{
			Method:           "SaveRegister",
			Router:           "/register/save",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beegoApi/controllers:VideoController"] = append(beego.GlobalControllerRouter["beegoApi/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "ChannelAdvert",
			Router:           "/channel/advert",
			AllowHTTPMethods: []string{"*"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beegoApi/controllers:VideoController"] = append(beego.GlobalControllerRouter["beegoApi/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "ChannelHostList",
			Router:           "/channel/hot",
			AllowHTTPMethods: []string{"*"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beegoApi/controllers:VideoController"] = append(beego.GlobalControllerRouter["beegoApi/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "ChannelRecommendRegionList",
			Router:           "/channel/recommend/region",
			AllowHTTPMethods: []string{"*"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beegoApi/controllers:VideoController"] = append(beego.GlobalControllerRouter["beegoApi/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "GetChannelRecommendTypeList",
			Router:           "/channel/recommend/type",
			AllowHTTPMethods: []string{"*"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beegoApi/controllers:VideoController"] = append(beego.GlobalControllerRouter["beegoApi/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "ChannelVideo",
			Router:           "/channel/video",
			AllowHTTPMethods: []string{"*"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beegoApi/controllers:VideoController"] = append(beego.GlobalControllerRouter["beegoApi/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "VideoEpisodesList",
			Router:           "/video/episodes/list",
			AllowHTTPMethods: []string{"*"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["beegoApi/controllers:VideoController"] = append(beego.GlobalControllerRouter["beegoApi/controllers:VideoController"],
		beego.ControllerComments{
			Method:           "VideoInfo",
			Router:           "/video/info",
			AllowHTTPMethods: []string{"*"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
