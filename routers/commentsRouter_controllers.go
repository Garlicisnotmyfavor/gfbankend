package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/gfbankend/controllers:CardController"] = append(beego.GlobalControllerRouter["github.com/gfbankend/controllers:CardController"],
        beego.ControllerComments{
            Method: "GetCardIDInfo",
            Router: `/card/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gfbankend/controllers:CardController"] = append(beego.GlobalControllerRouter["github.com/gfbankend/controllers:CardController"],
        beego.ControllerComments{
            Method: "Coupons",
            Router: `/card/:id/coupons`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gfbankend/controllers:CardController"] = append(beego.GlobalControllerRouter["github.com/gfbankend/controllers:CardController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/card/:id/delete`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gfbankend/controllers:CardController"] = append(beego.GlobalControllerRouter["github.com/gfbankend/controllers:CardController"],
        beego.ControllerComments{
            Method: "ModifyCardInfo",
            Router: `/card/:id/info`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gfbankend/controllers:CardController"] = append(beego.GlobalControllerRouter["github.com/gfbankend/controllers:CardController"],
        beego.ControllerComments{
            Method: "UseScore",
            Router: `/card/:id/score`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gfbankend/controllers:CardController"] = append(beego.GlobalControllerRouter["github.com/gfbankend/controllers:CardController"],
        beego.ControllerComments{
            Method: "AddCard",
            Router: `/card/add`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetAllCard",
            Router: `/:id:int`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"],
        beego.ControllerComments{
            Method: "NewPW",
            Router: `/ForgetPW/New`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"],
        beego.ControllerComments{
            Method: "Enroll",
            Router: `/enroll`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"],
        beego.ControllerComments{
            Method: "SendCodeInEnroll",
            Router: `/enroll`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"],
        beego.ControllerComments{
            Method: "ForgetPW",
            Router: `/forgetPw`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"],
        beego.ControllerComments{
            Method: "LoginWithCookie",
            Router: `/login`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/login`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/gfbankend/controllers:UserController"],
        beego.ControllerComments{
            Method: "ChangePW",
            Router: `/password`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
