package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"gfbankend/routers"
)

func main() {
	logs.Informational("App start")
	beego.Run()
}
