package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"
	"github.com/go-gomail/gomail"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"
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
// @Failure 400 Fail to unmarshal json
// @Failure 406 Illegal account form
// @Failure 403 Fail to insert
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
	pattern := "^[0-9a-zA-Z]+$" //匹配密码模式，仅允许字母数字
	ok, err := regexp.Match(pattern, []byte(user.Password))
	if !ok {
		models.Log.Error("Wrong Password")
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

/*
*@function:得到6位长的验证码
*@return {[]byte} 验证码
 */
func GetRandCode() []byte {
	var code []byte
	number := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Seed(time.Now().Unix())
	var sb strings.Builder
	size := len(number)
	for i := 0; i < 6; i++ {
		_, _ = fmt.Fprintf(&sb, "%d", number[rand.Intn(size)])
	}
	code = []byte(sb.String())
	return code
}

/*
*@function:发送验证码给target邮箱
*@param {string} 目标邮箱
*@return {[]byte}vcode，{error}err
 */
func SendEmail(target string) ([]byte, error) {
	//产生验证码
	vcode := GetRandCode()
	if len(vcode) != 6 {
		models.Log.Error("Error generating verify code")
		return nil, errors.New("Fail to generate verify code!")
	}
	//邮箱内容
	content := fmt.Sprintf("[ANZ]尊敬的客户' %s '，您本次登录所需的验证码为:%s,请勿向任何人提供您收到的验证码!", target, vcode)
	m := gomail.NewMessage()
	//设置邮件信息
	m.SetAddressHeader("From", "gfbankend@163.com", "ANZ-WORKSHOP") //设置发件人
	m.SetHeader("Subject", "Verify your device")                    //设置主题
	m.SetBody("text/html", content)                                 //设置主体内容
	m.SetHeader("To", m.FormatAddress(target, "收件人"))               //设置收件人
	//连接邮箱服务器并发送邮件
	d := gomail.NewPlainDialer("smtp.163.com", 465, "gfbankend@163.com", "ahz12345")

	if err := d.DialAndSend(m); err != nil {
		log.Println("Fail to send: ", err)
		return nil, err
	}
	return vcode, nil
}

// @Title Login
// @Description user login
// @Success 200 Register successfully
// @Failure 404 Fail to login
// @Failure 400 Fail to unmarshal json
// @router /login [put]
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
	fmt.Println(user)
	//查询用户信息是否与数据库匹配(现在匹配有bug，必须同时邮件、手机、密码）
	err1 := o.Read(&user, "mail", "password") //判断邮箱加密码
	err2 := o.Read(&user, "tel", "password") //判断手机加密码
	if err1 != nil && err2 != nil {
		models.Log.Error("read error: ", err2, err1)
		c.Ctx.ResponseWriter.WriteHeader(404) //读取用户信息错误
	}
	c.Ctx.ResponseWriter.WriteHeader(200) //信息匹配登录成功
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
