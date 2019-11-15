package main

import (
	_ "gfbankend/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	if beego.DEV == "dev" {
		beego.SetStaticPath("/swagger", "swagger")
	}
	logs.Informational("App start")
	beego.Run()
}
