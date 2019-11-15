package routers

import (
	"github.com/astaxie/beego"
	"gfbankend/controllers"
)

func init() {
	beego.Router("/api/user/card/:id", &controllers.CardController{})
	beego.Router("/api/user/card/", &controllers.CardController{})
	beego.Router("/api/user/", &controllers.UserController{})

}
