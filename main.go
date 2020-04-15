package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/gfbankend/routers"
)

var sessionID string = "bsessionID"

func main() {
	if beego.DEV == "dev" {
		beego.SetStaticPath("/swagger", "swagger")
	}
	logs.Informational("App start")
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = sessionID
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600 * 24
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = 3600
	beego.Run()
}
