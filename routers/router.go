// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
// 配置路由映射,以上注解必须添加
package routers

import (
	"github.com/astaxie/beego"
	"github.com/gfbankend/controllers"
)

func init() {
	beego.Router("/api/user/card/:id", &controllers.CardController{})
	beego.Router("/api/user/card/", &controllers.CardController{})
	beego.Router("/api/user/", &controllers.UserController{})
}

//func init() {
//	beego.Router("/api/user/card/:id", &controllers.CardController{})
//	beego.Router("/api/user/card/", &controllers.CardController{})
//	beego.Router("/api/user/", &controllers.UserController{})
//}

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/api/user",
			beego.NSNamespace("/card",
				beego.NSRouter("/", &controllers.CardController{}),
				beego.NSRouter("/:id", &controllers.CardController{}),
				beego.NSRouter("/all", &controllers.UserController{}, "get:GetAllCard"),
				beego.NSRouter("/help/:id", &controllers.CardController{}, "get:Help"),
			),
			beego.NSRouter("/", &controllers.UserController{}),       //get 返回用户资料
			beego.NSRouter("/join", &controllers.UserController{}),   //post
			beego.NSRouter("/login", &controllers.UserController{}),  //put
			beego.NSRouter("/logout", &controllers.UserController{}), //delete
			beego.NSRouter("/password", &controllers.UserController{}, "put:ChangePW"),
			beego.NSRouter("/feedback", &controllers.UserController{}, "post:Feedback"),
			beego.NSRouter("/garbage", &controllers.UserController{}, "get:GetDel"),
			beego.NSRouter("/garbage/:id", &controllers.UserController{}, "post:RecoverDel"),
		),
	)
	beego.AddNamespace(ns)
}
