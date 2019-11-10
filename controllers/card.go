package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"
)

type CardController struct {
	beego.Controller
}

type jsonStruct struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (c *CardController) Get() {
	// 获取路由参数
	id := c.Ctx.Input.Param(":id")
	o := orm.NewOrm()
	card := models.Card{Id: id}
	// 查询记录
	err := o.Read(&card)
	// 初始化返回
	json := &jsonStruct{-1, "not found", nil} // 推荐修改json这个变量名，因为json是包名
	if err == nil {
		// err为空，填充数据
		json.Code = 0
		json.Msg = "success"
		json.Data = card
	}
	c.Data["json"] = json
	c.ServeJSON()
}

func (c *CardController) Post() {
	var card models.Card
	body := c.Ctx.Input.RequestBody
	if err := json.Unmarshal(body, &card); err != nil {
		models.Log.Error("unmarshal error: ", err)
	}
	o := orm.NewOrm()
	if _, err := o.Insert(&card); err != nil {
		models.Log.Error("insert error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(403)
	}
	c.Ctx.ResponseWriter.WriteHeader(200)
	//c.ServeJSON()
}
