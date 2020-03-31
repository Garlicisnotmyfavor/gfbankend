package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"
	//"io/ioutil"
	//"path"
)

type UserController struct {
	beego.Controller
}

//显示所有卡片，修改输出的信息，不需要详细到卷
// @Title showAllCards
// @Description show all cards
// @Param    body        body     models.Card    true
// @Success 200 Read successfully
// @Failure 404 Fail to read
// @router / [get]
//zjn
func (c *UserController) GetAllCard() {
	//储存所有卡片信息
	var cardList []models.Card
	//使用orm接口查询相关信息
	o := orm.NewOrm()
	qt := o.QueryTable("card")
	//取出card表中所有信息，放入cardList中
	_, err := qt.All(&cardList)
	if err != nil {
		models.Log.Error("read error", err) //读取用户卡片信息失败
		c.Ctx.ResponseWriter.WriteHeader(403)
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200) //成功读取所有卡片
	//使用json格式传输所有信息
	c.Data["json"] = cardList
	//发送json
	c.ServeJSON()
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

//ML，用户注册 修改
// @Title Register
// @Description user register
// @Param user body models.User  true UserInfo
// @Success 200 {object} models.User "OK"
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
	if len(user.Tel) == 0 || len(user.Mail) == 0 {
		models.Log.Error("empty account")
		c.Ctx.ResponseWriter.WriteHeader(406) //非法账号
		return
	}
	//解析得到用户ID
	if err := user.UserParse(&user); err != nil {
		models.Log.Error("error in parsing user id: ", err)
		c.Ctx.ResponseWriter.WriteHeader(406) //用户解析出错
		return
	}
	if _, err := o.Insert(&user); err != nil {
		models.Log.Error("error in insert user: ", err)
		c.Ctx.ResponseWriter.WriteHeader(403) //插入错误
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200) //注册成功
}

//ML，登录，修改密码可调用ChangePw

// @Title Login
// @Description user login
// @Param userInfo body / true "account(string) + password(string) + accountType（string)为mail或者phone"
// @Success 200 Register successfully
// @Failure 403 Fail to login
// @Failure 400 Fail to unmarshal body
// @Failure 406 Illegal accountType
// @router /login [put]
func (c *UserController) Put() {
	o := orm.NewOrm()
	user := models.User{}

	body:=c.Ctx.Input.RequestBody
	var uInfo struct{
		account string
		password string
		accountType string
	}
	//解析前端JSON数据获得账号密码
	if err:=json.Unmarshal(body, &uInfo);err!=nil {
		models.Log.Error("Unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	var column string
	if uInfo.accountType == "mail" {
		user.Mail = uInfo.account
		user.Password = uInfo.password
		column = "mail"
	} else if uInfo.accountType == "phone" {
		user.Tel = uInfo.account
		user.Password = uInfo.password
		column = "tel"
	} else {
		//非法用户类型
		models.Log.Error("login error: wrong account type")
		c.Ctx.ResponseWriter.WriteHeader(406)
		return
	}
	if err := o.Read(&user, column, "password"); err != nil {
		models.Log.Error("login error: auth fail")
		c.Ctx.ResponseWriter.WriteHeader(403)
		return
	}
	//信息匹配登录成功
	c.Ctx.ResponseWriter.WriteHeader(200)
}

// @Title changePW
// @Description change password
// @Param    body        body     models.User    true
// @Success 200 Update successfully
// @Failure 404 Fail to read
//@Failure 400 Fail to unmarshal json
//@Failure 500 Fail to update
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

//忘记密码:用邮件找回，需要正确输入邮件验证码，验证通过后重新设置密码
//zjn
func (c *UserController) ForgetPW() {

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
