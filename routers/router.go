package routers

import (
	"github.com/astaxie/beego"
	"github.com/gfbankend/controllers"
)

func init() {
	beego.Router("/api/cards/:id", &controllers.CardController{})
	beego.Router("/api/cards/", &controllers.CardController{})
}
