package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"
	util "github.com/gfbankend/utils"
	//"io/ioutil"
	//"path"
)

type UserController struct {
	beego.Controller
}

//全局的session时长

//显示所有卡片
//检验是否在登陆状态，检验session是否存在，有的话不用前端的id，无的话返回错误操作
// @Title showAllCards
// @Description show all cards
// @Param    body        body     models.Card    true
// @Success 200 Read successfully
// @Failure 404 Fail to read
// @router /:id:int [get]
//zjn
func (c *UserController) GetAllCard() {
	if c.GetSession("userInfo") == nil {
		models.Log.Error("no login")
		c.Ctx.ResponseWriter.WriteHeader(403)
		return
	}
	// 取得用户ID from query
	uid := c.Ctx.Input.Param(":id")
	//储存所有卡片信息
	var cardList []models.Card
	//使用orm接口查询相关信息
	o := orm.NewOrm()
	qt := o.QueryTable("card")
	//取出card表中所有信息，放入cardList中
	_, err := qt.Filter("user_id__exact", uid).All(&cardList)
	if err != nil || len(cardList) == 0 {
		models.Log.Error("read error", err) //读取用户卡片信息失败
		c.Ctx.ResponseWriter.WriteHeader(404)
		return
	}
	//使用json格式传输所有信息
	c.Data["json"] = cardList
	//发送json
	c.ServeJSON()
	c.Ctx.ResponseWriter.WriteHeader(200) //成功读取所有卡片
}

//得到头像
//// @Title Get
//// @Description get current user's profile
//// @Param id query models.User  true
//// @Success 200 get successfully
//// @Failure 404 Fail to find picture
//// @router /:id  [get]
//func (c *UserController) Get() {
//	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", c.Ctx.Request.Header.Get("Origin"))
//	userId := c.Ctx.Input.Param(":id") //获取需要上传的文件文件名
//	filename:=userId+".jpg"
//	//查看是否存在需要的图片
// //   readPath :="D:/"
//    readPath := "/root/gfbankend/User/avatar/"
//	img:= path.Join(readPath,filename)
//	c.Ctx.Output.Header("Content-Type", "image/jpg")
//	c.Ctx.Output.Header("Content-Disposition",fmt.Sprintf("inline; filename=\"%s\"",img))
//	file, err := ioutil.ReadFile(img)
//	if err != nil {
//		models.Log.Error("read error", err) //未找到对应图片
//		c.Ctx.ResponseWriter.WriteHeader(404)
//		return
//	}
//	//c.Ctx.WriteString(string(file))
//	c.Data["avatar"]=string(file)
//}

//更新头像
//// @Title UpAvatar
//// @Description upload avatar
//// @Param id  header models.User  true file
//// @Success 200 upload successfully
//// @Failure 500 Fail to save picture
//// @Failure 502 Fail to close uploading file
//// @router /avatar [post]
//func (c *UserController) UpAvatar() {
//	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", c.Ctx.Request.Header.Get("Origin"))
//    userId:=c.GetString("id")
//	tmpFile, fHeader, err  := c.Ctx.Request.FormFile("avatar")
//	if err != nil{
//		models.Log.Error("read error", err) //读取用户卡片信息失败
//		c.Ctx.ResponseWriter.WriteHeader(400)
//	}
//	//savePath := "D:/" +userId+".jpg"
//	savePath := "/root/gfbankend/User/avatar/" +userId+".jpg"        //设置保存路径
//	beego.Info("Header:", fHeader.Header)     //map[Content-Disposition:[form-data; name="123"; filename="upimage.jpg"] Content-Type:[image/jpeg]]
//	beego.Info("Size:", fHeader.Size)         //114353
//	beego.Info("Filename:", fHeader.Filename) //upimage.jpg
//	if err=c.SaveToFile("avatar", savePath);err !=nil{
//		models.Log.Error("save error", err) //存储图片失败
//		c.Ctx.ResponseWriter.WriteHeader(500)
//	}
//	if err:=tmpFile.Close();err!=nil {
//		models.Log.Error("close error", err) //存储图片失败
//		c.Ctx.ResponseWriter.WriteHeader(502)
//	}                   //关闭上传的文件，不然的话会出现临时文件不能清除的情况
//}

//ML，用户注册时验证码获取
// @Title getRanCodeInRegister
// @Description send random code when user enroll
// @Param	email	body	string	true	用户的邮箱
// @Success 200	string	"生成的验证码"
// @Failure 400 解析body失败
// @Failure 500 发送邮件失败
// @router /enroll [get]
func (c *UserController) SendCodeInEnroll() {
	var email string // this is user's email
	body := c.Ctx.Input.RequestBody
	// get email from body
	if err := json.Unmarshal(body, &email); err != nil {
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	randCode := util.GetRandCode() // get random code
	if err := util.SendEmail(email, randCode); err != nil {
		models.Log.Error("send email error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	c.Data["json"] = randCode
	c.ServeJSON()
	c.Ctx.ResponseWriter.WriteHeader(200)
}

//ML，用户注册
// @Title Register
// @Description user register
// @Param userInfo body models.User  true 用户所填信息
// @Success 200 {object} models.User "OK"
// @Failure 400 解析body错误
// @Failure 406 账号信息格式有误
// @Failure 403 数据库插入错误
// @router /enroll [post]
func (c *UserController) Enroll() {
	o := orm.NewOrm()
	body := c.Ctx.Input.RequestBody
	user := models.User{}
	//Obtain information of the new user
	if err := json.Unmarshal(body, &user); err != nil {
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	//检查用户手机或者邮箱是否为空
	if len(user.Tel) == 0 || len(user.Mail) == 0 {
		models.Log.Error("empty account")
		c.Ctx.ResponseWriter.WriteHeader(406) //非法账号
		return
	}
	//解析得到用户ID
	if err := user.UserParse(); err != nil {
		models.Log.Error("error in parsing user id: ", err)
		c.Ctx.ResponseWriter.WriteHeader(406) //用户ID解析出错
		return
	}
	if _, err := o.Insert(&user); err != nil {
		models.Log.Error("error in insert user: ", err)
		c.Ctx.ResponseWriter.WriteHeader(403) //插入错误
		return
	}
	c.Data["json"] = user
	c.ServeJSON()
	c.Ctx.ResponseWriter.WriteHeader(200) //注册成功
}

// 加入是否选择记住密码，设置session，设置cookie
// @Title Login
// @Description user login
// @Param userInfo body true account(string)+password(string)+accounttype(string)为mail或者phone
// @Success 200 {object} models.User Register successfully
// @Failure 406 数据库查询报错，可能用户所填账号或密码错误
// @Failure 400 信息内容或格式有误
// @router /login [put]
func (c *UserController) Login() {
	o := orm.NewOrm()
	user := models.User{}
	body := c.Ctx.Input.RequestBody
	// temp struct to get userInfo
	var uInfo struct {
		Account     string
		Password    string
		AccountType string
	}
	// 解析前端JSON数据获得账号密码
	if err := json.Unmarshal(body, &uInfo); err != nil {
		models.Log.Error("Unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	var column string
	if uInfo.AccountType == "mail" {
		user.Mail = uInfo.Account
		user.Password = uInfo.Password
		column = "mail"
	} else if uInfo.AccountType == "phone" {
		user.Tel = uInfo.Account
		user.Password = uInfo.Password
		column = "tel"
	} else {
		// 非法用户类型
		models.Log.Error("login error: wrong account type")
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	if err := o.Read(&user, column, "password"); err != nil {
		models.Log.Error("login error: auth fail")
		c.Ctx.ResponseWriter.WriteHeader(403)
		return
	}
	// 信息匹配登录成功
	c.Data["json"] = user
	c.ServeJSON()                  // 传用户对象给前端
	c.SetSession("userInfo", user) // 登录成功，设置session
	c.Ctx.ResponseWriter.WriteHeader(200)
}

// @Title test
// @Description user test
// @Failure 406 数据库查询报错，可能用户所填账号或密码错误
// @Failure 400 信息内容或格式有误
// @router	/cookie/test	[get]

// @Title changePW
// @Description change password
// @Param userInfo body models.User true 用户信息(需要的是用户ID，新密码）
// @Success 200 Update successfully
// @Failure 404 数据库无此用户
// @Failure 400 解析body失败
// @Failure 406 更新密码失败
// @router /password [put]
func (c *UserController) ChangePW() {
	var user models.User
	body := c.Ctx.Input.RequestBody
	if err := json.Unmarshal(body, &user); err != nil {
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	o := orm.NewOrm()
	usr := models.User{Id: user.Id}
	// 查询记录
	if err := o.Read(&usr); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) // 查不到id对应的用户
		return
	}
	//查询成功，更新密码
	usr.Password = user.Password
	if _, err := o.Update(&usr); err != nil {
		models.Log.Error("update error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(500) // 更新数据失败
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200) // 更新成功
}

//忘记密码:用邮件找回，需要正确输入邮件验证码，验证通过后重新设置密码
//zjn
// @Title ForgetPW
// @Description Forget password
// @Param userInfo body models.User true 用户信息(需要的是用户ID，邮件）
// @Success 200 successfully
// @Failure 404 数据库无此用户
// @Failure 400 解析body失败
// @router /forgetPw [post]
func (c *UserController) ForgetPW() {
	//输入相关信息，解析出用户邮箱
	var user models.User
	body := c.Ctx.Input.RequestBody
	if err := json.Unmarshal(body, &user); err != nil {
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	o := orm.NewOrm()
	//查询填入内容是否准确
	if err := o.Read(&user); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) // 查不到对应id的用户
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200) // 身份验证成功，后续进入验证码阶段
}

// @Title NewPassword
// @Description  通过前面忘记密码的过程后，设置新的密码
// @Param userInfo body models.User true 用户信息(需要的是用户ID，新密码）
// @Success 200 Update successfully
// @Failure 404 数据库无此用户
// @Failure 400 解析body失败
// @Failure 406 更新密码失败
// @router /ForgetPW/New [put]
func (c *UserController) NewPW() {
	var user models.User
	body := c.Ctx.Input.RequestBody
	if err := json.Unmarshal(body, &user); err != nil {
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	o := orm.NewOrm()
	usr := models.User{Id: user.Id}
	// 查询记录
	if err := o.Read(&usr); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) // 查不到id对应的用户
		return
	}
	//查询成功，更新密码
	usr.Password = user.Password
	if _, err := o.Update(&usr); err != nil {
		models.Log.Error("update error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(500) // 更新数据失败
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200) // 更新成功
}

//// @Title Feedback
//// @Description send feedback mail
//// @Success 200 Update successfully
////@Failure 500 Fail to send mail
//// @router /password [post]
//func (c *UserController) Feedback() {
//	body:=c.GetString("feedback")
//	body = fmt.Sprintf(body)
//	serverHost := "smtp.163.com"
//	serverPort := 465
//	toEmail:="1725500398@qq.com"
//	fromEmail := "gfbankend@163.com"
//	fromPasswd := "ahz12345"
//	var m *gomail.Message
//	m = gomail.NewMessage()
//	m.SetAddressHeader("From",fromEmail,"ANZ-WORKSHOP")
//	m.SetHeader("To",toEmail)
//	m.SetHeader("Subject", "Feedback")
//	m.SetBody("text/html",body)
//	d:=gomail.NewPlainDialer(serverHost,serverPort,fromEmail,fromPasswd)
//	if err :=d.DialAndSend(m);err!=nil{
//		models.Log.Error("feedback error: ", err)
//		c.Ctx.ResponseWriter.WriteHeader(500) // 发送邮件失败
//		return
//	}
//}
