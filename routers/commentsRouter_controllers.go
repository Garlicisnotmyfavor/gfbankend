package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

// 初始get网页的操作
func init() {

    beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})
}
