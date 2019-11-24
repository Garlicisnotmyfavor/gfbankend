package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"
	"time"
	"fmt"
)

type CardController struct {
	beego.Controller
}

// swagger注解配置
// @Title Get
// @Description get card
// @router /card/:id([0-9]+) [get]
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

// swagger注解配置
// @Title Post
// @Description insert card
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

// swagger注解配置
// @Title Delete
// @Description delete card
// @router /card/:id [delete]
func (c *CardController) Delete() {
	id := c.Ctx.Input.Param(":id")
	//fmt.Println(id)
	o := orm.NewOrm()
	card := models.Card{Id: id}
	if err := o.Read(&card); err == nil {
		count, _ := o.Delete(&card)
		if count == 0 {
			models.Log.Error("delete fail") //删除0个元素，即删除失败，返回状态码403
			c.Ctx.ResponseWriter.WriteHeader(403)
		} else {
			delCard := models.DelCard{CardId: card.Id, UserId: card.UserId, Remark: card.Remark}
			delCard.DelTime = time.Now()
			_, err := o.Insert(&delCard)
			if err != nil {
				models.Log.Error("Insert error: ", err) //被删卡插入垃圾箱失败
				c.Ctx.ResponseWriter.WriteHeader(403)
				return
			}
			c.Ctx.ResponseWriter.WriteHeader(200) //删除成功
		}
	} else {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(200) //card本就不存在，删除不存在的卡当作删除成功
	}
}

//GZH，修改备注
func (c *CardController) Put() {
	// 接收数据
	id := c.Ctx.Input.Param(":id")
	remark := c.GetString("remark")
	if remark == "" {
		// remark参数为空，设置400状态码
		c.Ctx.Output.SetStatus(400)
	} else {
		o := orm.NewOrm()
		card := models.Card{Id: id}
		err := o.Read(&card);
		if err == orm.ErrNoRows {
			c.Ctx.Output.SetStatus(404)
		} else if err == orm.ErrMissPK {
			c.Ctx.Output.SetStatus(404)
		} else {
			// 查到了该记录，进行赋值
			card.Remark = remark
			// 更新记录
			_, _ = o.Update(&card)
			// 设置成功响应
			c.Ctx.Output.SetStatus(200)
			c.Data["json"] = card
		}
	}
	c.ServeJSON()
}

// swagger注解配置
// @Title Get
// @Param Ename query string true "enterprise_name"
// @Description get help message by the given enterprise_name
// @Success 200
// @Failure 404 read error
// @router /card/help/:Ename [get]
func (c *CardController) Help() {
	EName := c.Ctx.Input.Param(":Ename")
	fmt.Println(EName)
	o := orm.NewOrm()
	enterprise := models.Enterprise{Name: EName}
	// 查询记录
	if err := o.Read(&enterprise); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) // 查不到企业名对应的企业
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200) //成功
	c.Data["json"] = enterprise.HelpMsg
	c.ServeJSON()
}
