package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_"github.com/gfbankend/routers"
)

func main() {
	if beego.DEV == "dev" {
		beego.SetStaticPath("/swagger", "swagger")
	}
	logs.Informational("App start")
	beego.Run()
}
