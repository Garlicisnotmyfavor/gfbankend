package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/gfbankend/controllers:CardController"] = append(beego.GlobalControllerRouter["github.com/gfbankend/controllers:CardController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/card`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gfbankend/controllers:CardController"] = append(beego.GlobalControllerRouter["github.com/gfbankend/controllers:CardController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/card/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetAllCard",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/join`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/login`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
