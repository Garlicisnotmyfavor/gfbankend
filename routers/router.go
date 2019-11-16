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
			),
			beego.NSRouter("/", &controllers.UserController{}),
		),
	)
	beego.AddNamespace(ns)
}
