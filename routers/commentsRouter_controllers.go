package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["beegoApi/controllers:UserController"] = append(beego.GlobalControllerRouter["beegoApi/controllers:UserController"],
		beego.ControllerComments{
			Method:           "SaveRegister",
			Router:           `/register/save`,
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
