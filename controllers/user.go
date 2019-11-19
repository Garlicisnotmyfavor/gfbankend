package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"
	"regexp"
)

type UserController struct {
	beego.Controller
}

// @Title showAllCards
// @Description show all cards
// @Param    body        body     models.Card    true
// @Success 200 Read successfully
// @Failure 404 Fail to read
// @router / [get]
func (c *UserController) GetAllCard() {
	var cardList []orm.Params                 //存储所有卡片信息
	sql := fmt.Sprintf(`select * from card;`) //需要卡的table名
	o := orm.NewOrm()
	//根据sql指令将table中所有卡信息读入到carList中
	_, err := o.Raw(sql).Values(&cardList)
	if err != nil {
		models.Log.Error("read error", err) //读取用户卡片信息失败
		c.Ctx.ResponseWriter.WriteHeader(404)
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200) //成功读取所有卡片
	c.Data["json"] = cardList
	//发送json
	c.ServeJSON()
}

//YZY，返回用户资料
func (c *UserController) Get() {

}

//ML，用户注册
// @Title Register
// @Description user register
// @Success 200 Register successfully
// @Failure 404 Fail to register
// @router /join [post]
func (c *UserController) Post() {
	o := orm.NewOrm()
	body := c.Ctx.Input.RequestBody
	user := models.User{}
	//Obtain information of the new user
	if err := json.Unmarshal(body, &user); err != nil {
		models.Log.Error("Unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	//检查用户手机或者邮箱是否为空
	if len(user.Tel) == 00 || len(user.Mail) == 0 {
		models.Log.Error("Wrong Phone or Email")
		c.Ctx.ResponseWriter.WriteHeader(406) //非法用户ID
		return
	}
	//正则表达式匹配密码是否合法
	pattern := "$[0-9a-zA-Z]+^"
	ok, err := regexp.Match(pattern, []byte(user.Password))
	if !ok {
		models.Log.Error("wrong Password")
		c.Ctx.ResponseWriter.WriteHeader(406) //非法用户密码
		return
	}
	//解析得到用户ID
	err = user.UserParse()
	if err != nil {
		models.Log.Error("error in parsing user id: ", err)
		c.Ctx.ResponseWriter.WriteHeader(406) //用户解析出错
		return
	}
	_, err = o.Insert(&user)
	if err != nil {
		models.Log.Error("error in insert user: ", err)
		c.Ctx.ResponseWriter.WriteHeader(403) //插入错误
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200) //注册成功
}

//ML，登录，修改密码可调用ChangePw
func (c *UserController) Put() {
	o := orm.NewOrm()
	body := c.Ctx.Input.RequestBody
	user := models.User{}
	//取得用户信息
	if err := json.Unmarshal(body, &user); err != nil {
		models.Log.Error("Unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	//发送验证码确认登录





	//查询用户是否与在数据库中的信息匹配
	if err := o.Read(&user); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) //用户读取错误
	}

}

//ZJN，显示所有被删卡片
func (c *UserController) GetDel() {

}

//ZJN，恢复指定卡片
func (c *UserController) RecoverDel() {

}

//YZY，修改密码
func (c *UserController) ChangePW() {

}

//YZY，反馈
func (c *UserController) Feedback() {

}
