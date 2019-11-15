package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_"gfbankend/routers"
)

func main() {
	if beego.DEV == "dev" {
		beego.SetStaticPath("/swagger", "swagger")
	}
	logs.Informational("App start")
	beego.Run("101.37.27.155")
}
