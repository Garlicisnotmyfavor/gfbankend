package routers

import (
	"github.com/astaxie/beego"
	"github.com/gfbankend/controllers"
)

func init() {
	beego.Router("/api/user/cards/:id", &controllers.CardController{})
	beego.Router("/api/user/cards/", &controllers.CardController{})

}
