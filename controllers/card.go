package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"
	"time"
)

type CardController struct {
	beego.Controller
}

// swagger注解配置
// @Title Get
// @Description get card
// @router /card/:id [get]
func (c *CardController) Get() {
	// 获取路由参数
	id := c.Ctx.Input.Param(":id")
	o := orm.NewOrm()
	card := models.Card{Id: id}
	// 查询记录
	if err := o.Read(&card); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) // 查不到id对应的卡
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200) //成功
	c.Data["json"] = card
	c.ServeJSON()
}

// @Title write by gu yu hao
// @Description 插入一张卡片
// @Param card body  json{"Id":string,"UserId":string,"Remark":string} true "只需要给出卡Id以及用户Id和备注，即可解析出额外信息"
// @Success 200 插入卡片成功
// @Failure 400 解析json错误
// @Failure 403 插入失败
// @Failure 406 非法的用户ID
// @router /card [post]
func (c *CardController) Post() {
	var card models.Card
	body := c.Ctx.Input.RequestBody
	if err := json.Unmarshal(body, &card); err != nil {
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	if err := card.CardParse(); err != nil {
		models.Log.Error("card parse error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(406) //非法的用户ID
		return
	}
	o := orm.NewOrm()
	if _, err := o.Insert(&card); err != nil {
		models.Log.Error("insert error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(403) //插入错误
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200) //成功
}

func (c *CardController) Delete() {
	id := c.Ctx.Input.Param(":id")
	o := orm.NewOrm()
	card := models.Card{Id: id}
	if err := o.Read(&card); err == nil {
		count, _ := o.Delete(&card)
		if count == 0 {
			models.Log.Error("delete fail")
		} else {
			delCard := models.DelCard{CardId: card.Id, UserId: card.UserId, Remark: card.Remark}

			delCard.DelTime = time.Now()
			_, err := o.Insert(&delCard)
			if err != nil {
				models.Log.Error("Insert error: ", err)
				c.Ctx.ResponseWriter.WriteHeader(403)
				return
			}
			c.Ctx.ResponseWriter.WriteHeader(200)
		}
	} else {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(200) //card本就不存在
	}
}

//GZH，修改备注
func (c *CardController) Put() {

}

//GYH,显示卡片帮助信息
func (c *CardController) Help() {

}
