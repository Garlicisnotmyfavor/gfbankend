package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"
	"github.com/go-gomail/gomail"
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
func (c *UserController) Post() {

}

//ML，登录，修改密码可调用ChangePw
func (c *UserController) Put() {

}

//ZJN，显示所有被删卡片
func (c *UserController) GetDel() {

}

//ZJN，恢复指定卡片
func (c *UserController) RecoverDel() {

}

//YZY，修改密码，要求传输一个有用户ID，新密码的json
func (c *UserController) ChangePW() {
	var user models.User
	body := c.Ctx.Input.RequestBody
	if err := json.Unmarshal(body, &user); err != nil {
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	o := orm.NewOrm()
	usr:= models.User{Id: user.Id}
	// 查询记录
	if err := o.Read(&usr); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) // 查不到id对应的用户
		return
	}
	//查询成功，更新密码
	usr.Password=user.Password
	if _,err:=o.Update(&usr);err!=nil{
		models.Log.Error("update error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(500) // 更新数据失败
		return
	}
}

//YZY，反馈
func (c *UserController) Feedback() {
	body:=c.GetString("feedback")
	body = fmt.Sprintf(body)
	serverHost := "smtp.163.com"
	serverPort := 465
	toEmail:="1725500398@qq.com"
	fromEmail := "gfbankend@163.com"
	fromPasswd := "ahz12345"
	var m *gomail.Message
	m = gomail.NewMessage()
	m.SetAddressHeader("From",fromEmail,"ANZ-WORKSHOP")
	m.SetHeader("To",toEmail)
	m.SetHeader("Subject", "Feedback")
	m.SetBody("text/html",body)
	d:=gomail.NewPlainDialer(serverHost,serverPort,fromEmail,fromPasswd)
	if err :=d.DialAndSend(m);err!=nil{
		models.Log.Error("feedback error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(500) // 发送邮件失败
		return
	}
}
