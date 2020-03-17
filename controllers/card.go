package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"
)

type CardController struct {
	beego.Controller
}

//查询指定id对应的卡片——卡片的类型信息，所有信息
func (c *CardController) Get() {
	// 获取路由参数
	id := c.Ctx.Input.Param(":id")
	fmt.Println(id)
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

//添加卡片 在user表里添加此user和card的关联
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

//修改卡片的id和公司名——我们认为不需要这个

//nfc扫码操作积分
func (c *CardController) Get1(){}

//对优惠券的操作
func (c *CardController) Get2(){}

//删除卡片 ——超过时限的一些卡
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
