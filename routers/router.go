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

//func init() {
//	ns := beego.NewNamespace("/v1",
//		beego.NSNamespace("/api/user",
//			beego.NSNamespace("/card",
//				beego.NSRouter("/", &controllers.CardController{}),                      //post 添加卡片
//				beego.NSRouter("/:id", &controllers.CardController{}),                   //get 显示单张卡片详细信息
//				beego.NSRouter("/:id/remark", &controllers.CardController{}),            //put 修改卡片备注
//				beego.NSRouter("/all", &controllers.UserController{}, "get:GetAllCard"), //显示所有卡片
//				beego.NSRouter("/help/:id", &controllers.CardController{}, "get:Help"),  //单张卡片帮助信息
//			),
//			beego.NSRouter("/", &controllers.UserController{}),                               //get 返回用户资料
//			beego.NSRouter("/join", &controllers.UserController{}),                           //post 用户注册
//			beego.NSRouter("/login", &controllers.UserController{}),                          //put 用户登录
//			beego.NSRouter("/logout", &controllers.UserController{}),                         //delete 用户退出登录
//			beego.NSRouter("/password", &controllers.UserController{}, "put:ChangePW"),       //修改密码
//			beego.NSRouter("/feedback", &controllers.UserController{}, "post:Feedback"),      //反馈
//			beego.NSRouter("/garbage", &controllers.UserController{}, "get:GetDel"),          //显示所有被删除卡片
//			beego.NSRouter("/garbage/:id", &controllers.UserController{}, "post:RecoverDel"), //恢复被选中卡片
//		),
//	)
//	beego.AddNamespace(ns)
//}

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/api/user",
			beego.NSInclude(
				&controllers.UserController{},
				&controllers.CardController{},
			),
		),
	)
		beego.AddNamespace(ns)
}
