package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"
)

type UserController struct {
	beego.Controller
}

//显示所有卡片，修改输出的信息
// @Title showAllCards
// @Description show all cards
// @Param    body        body     models.Card    true
// @Success 200 Read successfully
// @Failure 404 Fail to read
// @router / [get]
func (c *UserController) Get() {
	var cardList []orm.Params                 //存储所有卡片信息
	sql := fmt.Sprintf(`select * from card;`) //需要卡的table名
	o := orm.NewOrm()
	//根据sql指令将table中所有卡信息读入到carList中
	_, err := o.Raw(sql).Values(&cardList)
	if err != nil {
		models.Log.Error("read error", err) //读取用户卡片信息失败
		c.Ctx.ResponseWriter.WriteHeader(404)
	}
	c.Ctx.ResponseWriter.WriteHeader(200) //成功读取所有卡片
	c.Data["json"] = cardList
	//发送json
	c.ServeJSON()
}
