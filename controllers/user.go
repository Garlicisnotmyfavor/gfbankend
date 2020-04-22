package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"
	util "github.com/gfbankend/utils"
	//util "github.com/gfbankend/utils"
	//"io/ioutil"
	//"path"
)

type UserController struct {
	beego.Controller
}

//全局的session时长

//显示所有卡片
// 检验是否在登陆状态，检验session是否存在，有的话不用前端的id，无的话返回错误操作
// @Title showAllCards
// @Description show all cards
// @Param    userID        path    string    true	用户ID
// @Success 200 Read successfully
// @Failure 404 Fail to read
// @Failure 401 没登录，无权限
// @router /:id:int [get]
//zjn
func (c *UserController) GetAllCard() {
	if c.GetSession("userInfo") == nil {
		models.Log.Error("no login")
		c.Ctx.ResponseWriter.WriteHeader(401)
		return
	}
	// 取得用户ID from path
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
	var Response struct {
		Status int	`json:"status"`
		Msg	string `json:"msg"`
		Data []models.Card	`json:"data"`
	}
	Response.Status = 200
	Response.Msg = "success"
	Response.Data = cardList
	//使用json格式传输所有信息
	c.Data["json"] = Response
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

//ML，用户注册时验证码获取
// @Title getRanCodeInRegister
// @Description send random code
// @Param	email	body	string	true	"email":xxx
// @Success 200	string	"生成的验证码"
// @Failure 400 解析body失败
// @Failure 500 发送邮件失败
// @router /verify [post]
func (c *UserController) SendCode() {
	var email struct{ Email string } // this is user's email
	body := c.Ctx.Input.RequestBody
	// get email from body
	if err := json.Unmarshal(body, &email); err != nil {
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	randCode := util.GetRandCode() // get random code
	if err := util.SendEmail(email.Email, randCode); err != nil {
		models.Log.Error("send email error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data string `json:"data"`
	}
	response.Code = 200
	response.Msg = "success"
	response.Data = string(randCode)
	// 将验证码加入到用户对应的Session
	c.SetSession("verify", string(randCode))
	c.Data["json"] = response
	c.ServeJSON()
	c.Ctx.ResponseWriter.WriteHeader(200)
}

//ML，用户注册
// @Title Register
// @Description user register
// @Param userInfo body \  true 电话tel+邮箱mail+密码password+验证码verify
// @Success 200 返回值：结构体其中有msg信息，数据data
// @Failure 400 解析json出错，返回值：具体错误信息{"msg":xxx}
// @Failure 406 信息有误，详见返回值：错误信息{"msg":xxx}
// @Failure 403 插入数据库错误，返回值：错误信息{"msg":xxx}
// @router /enroll [post]
func (c *UserController) Enroll() {
	o := orm.NewOrm()
	body := c.Ctx.Input.RequestBody
	var userInfo struct {
		// ID       string
		Tel      string
		Mail     string
		Password string
		//Verify   string
	}
	//Obtain information of the new user
	if err := json.Unmarshal(body, &userInfo); err != nil {
		models.Log.Error("empty account")
		var response struct {
			Msg string `json:"msg"`
		}
		response.Msg = "fail unmarshal json"
		c.Data["json"] = response
		c.ServeJSON()
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	//检查用户手机或者邮箱是否为空
	if len(userInfo.Tel) == 0 || len(userInfo.Mail) == 0 {
		models.Log.Error("empty account")
		var response struct {
			Msg string `json:"msg"`
		}
		response.Msg = "empty mail or phone"
		c.Data["json"] = response
		c.ServeJSON()
		c.Ctx.ResponseWriter.WriteHeader(401) //非法账号
		return
	}

	if len(userInfo.Password) == 0 {
		models.Log.Error("empty password")
		var response struct {
			Msg string `json:"msg"`
		}
		response.Msg = "empty password"
		c.Data["json"] = response
		c.ServeJSON()
		c.Ctx.ResponseWriter.WriteHeader(402) //没输入密码
		return
	}
	// sess := c.GetSession("verify")
	//if sess == nil {
	// 	models.Log.Error("enroll without being verified")
	// 	var response struct {
	// 		Msg string `json:"msg"`
	// 	}
	// 	response.Msg = "no verification"
	// 	c.Data["json"] = response
	// 	c.ServeJSON()
	// 	c.Ctx.ResponseWriter.WriteHeader(406) //没有点击验证码
	// 	return
	// }
	 //vCode := sess.(string)
	 //if vCode != userInfo.Verify {
	 //	models.Log.Error("verify fail")
	 //	var response struct {
	// 		Msg string `json:"msg"`
	 //	}
	// 	response.Msg = "wrong verify code"
	// 	c.Data["json"] = response
	// 	c.ServeJSON()
	// 	c.Ctx.ResponseWriter.WriteHeader(403) //没有点击验证码
	// 	return
	// }
	 //ready to insert new user
	user := models.User{
		// Id:       userInfo.ID,
		Tel:      userInfo.Tel,
		Mail:     userInfo.Mail,
		Password: userInfo.Password,
	}
	//fmt.Println(user)
	if err:= o.Read(&user, "tel", "mail"); err == nil {
		models.Log.Error("enroll error")
		var response struct {
			Msg string `json:"msg"`
		}
		response.Msg = "email or phone already used"
		c.Data["json"] = response
		c.ServeJSON()
		return
	}
	user.UserParse()
	if _, err := o.Insert(&user); err != nil {
		models.Log.Error("error in insert user: ", err)
		var response struct {
			Msg string `json:"msg"`
		}
		response.Msg = "insert fail"
		c.Data["json"] = response
		c.ServeJSON()
		c.Ctx.ResponseWriter.WriteHeader(405) //插入错误
		return
	}
	// response to front-end
	var response struct {
		Msg  string      `json:"msg"`
		Data models.User `json:"data"`
	}
	// 注册成功销毁验证码
	c.DelSession("verify")
	response.Msg = "success"
	response.Data = user
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title LoginWithCookie
// @Description 通过获取客户端的cookie来进行自动登录cookie在正常登录时选择自动登录后获得,失败则消除cookie
// @Success 200 通过cookie完成了自动登录
// @Success 204 cookie失效
// @Fail  406 cookie中的消息过期或者错误，比如另一客户端修改密码后,需要重新登录
// @router /login [get]
func (c *UserController) LoginWithCookie() {
	account, _ := c.Ctx.GetSecureCookie("miller", "account")
	password, _ := c.Ctx.GetSecureCookie("miller", "password")
	accountType := c.Ctx.GetCookie("accounttype")
	remember := c.Ctx.GetCookie("remember")
	if account == "" || password == "" || accountType == "" {
		models.Log.Error("fail to get cookie")
		c.Ctx.ResponseWriter.WriteHeader(204)
		return
	}
	// 由cookie得到用户信息
	var column string
	user := models.User{Password: password}
	if accountType == "mail" {
		user.Mail = account
		column = "mail"
	} else {
		user.Tel = account
		column = "tel"
	}
	o := orm.NewOrm()
	// 登录失败，说明原cookie失效（比如另一客户端修改密码之后），需老实点重新登录
	if err := o.Read(&user, column, "password"); err != nil {
		models.Log.Error("login error: auth fail")
		// 删除原cookie
		c.Ctx.SetCookie("account", "", -1)
		c.Ctx.SetCookie("password", "", -1)
		c.Ctx.SetCookie("accounttype", "", -1)
		c.Ctx.SetCookie("remember", "", -1)
		c.Ctx.ResponseWriter.WriteHeader(406)
		return
	}
	if remember != "true" {
		c.Ctx.SetCookie("account", "", -1)
		c.Ctx.SetCookie("password", "", -1)
		c.Ctx.SetCookie("accounttype", "", -1)
		c.Ctx.SetCookie("remember", "", -1)
	}
	// 自动获取cookie后经过信息检验登录成功,需重新设置该客户端中sessionid对应session
	c.SetSession("userInfo", user)
	c.Data["json"] = user
	c.ServeJSON()
}

// @Title Login
// @Description 登录，返回中由名为bsessionID的cookie，用于查用户是否已登录如果点击了自动登录（即remember为true）还会返回名为account，password，accounttype以及remember的cookie
// @Param	userInfo	body	/	true account(string)+password(string)+accounttype(string)为mail或者phone+remember(是否记住密码bool)
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
		Remember    bool
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
		c.Ctx.ResponseWriter.WriteHeader(406)
		return
	}
	// 信息匹配登录成功
	c.Data["json"] = user
	// 如果需要记住账号密码
	if uInfo.Remember == true {
		c.Ctx.SetSecureCookie("miller", "account", uInfo.Account)
		c.Ctx.SetSecureCookie("miller", "password", uInfo.Password)
		c.Ctx.SetCookie("accounttype", uInfo.AccountType)
		c.Ctx.SetCookie("remember", "true")
	}
	c.SetSession("userInfo", user) // 登录成功，设置session
	c.ServeJSON()                  // 传用户对象给前端
}

// @Title test
// @Description user test
// @Failure 406 数据库查询报错，可能用户所填账号或密码错误
// @Failure 400 信息内容或格式有误
// @router	/cookie/test	[get]

// @Title changePW
// @Description change password
// @Param userInfo body / true 用户信息(需要的是用户ID，原密码，新密码）
// @Success 200 Update successfully
// @Failure 401 没登录，无权限
// @Failure 404 数据库无此用户
// @Failure 403 数据库无此用户
// @Failure 400 解析body失败
// @Failure 406 更新密码失败
// @router /password [put]
func (c *UserController) ChangePW() {
	if c.GetSession("userInfo") == nil {
		models.Log.Error("no login")
		c.Ctx.ResponseWriter.WriteHeader(401)
		return
	}
	var user struct {
		UserId      string
		OldPassword string
		NewPassword string
	}
	body := c.Ctx.Input.RequestBody
	if err := json.Unmarshal(body, &user); err != nil {
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	o := orm.NewOrm()
	usr := models.User{Id: user.UserId}
	// 查询记录
	if err := o.Read(&usr); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) // 查不到id对应的用户
		return
	}
	//验证用户输入的密码是否与旧密码一致
	if user.OldPassword != usr.Password {
		models.Log.Error("wrong old password: ")
		c.Ctx.ResponseWriter.WriteHeader(403) // 原密码错误
		return
	}
	usr.Password = user.NewPassword
	if _, err := o.Update(&usr); err != nil {
		models.Log.Error("update error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(500) // 更新数据失败
		return
	}
	//根据旧的userInfo删除对应session
	c.DelSession("userInfo")
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
	//validMail := true
	// validTel := true
	// validId := true"
	//查询填入内容是否准确
	if err := o.Read(&user); err != nil {
		//validId = false
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) // 查不到对应邮箱或手机号的用户
		return
	}
	// if err := o.Read(&user,"mail"); err != nil {
	// 	validMail = false
	// }
	// if err := o.Read(&user,"tel"); err != nil {
	// 	validTel = false
	// }
	// if validTel==false&&validMail==false {
	// 	models.Log.Error("read error: ", err)
	// 	c.Ctx.ResponseWriter.WriteHeader(404) // 查不到对应邮箱或手机号的用户
	// 	return
	// }
	c.Ctx.ResponseWriter.WriteHeader(200) // 身份验证成功，后续进入验证码阶段
}

// //ML，用户注册时验证码获取
// // @Title getRanCodeInRegister
// // @Description send random code when user enroll
// // @Param	email	body	string	true	用户的邮箱
// // @Success 200	string	"生成的验证码"
// // @Failure 400 解析body失败
// // @Failure 500 发送邮件失败
// // @router /forgetPw/New [put]
// func (c *UserController) SendCodeInNew() {
// 	var email struct{ Email string } // this is user's email
// 	body := c.Ctx.Input.RequestBody
// 	// get email from body
// 	if err := json.Unmarshal(body, &email); err != nil {
// 		models.Log.Error("unmarshal error: ", err)
// 		c.Ctx.ResponseWriter.WriteHeader(400)
// 		return
// 	}
// 	randCode := util.GetRandCode() // get random code
// 	if err := util.SendEmail(email.Email, randCode); err != nil {
// 		models.Log.Error("send email error: ", err)
// 		c.Ctx.ResponseWriter.WriteHeader(500)
// 		return
// 	}
// 	c.Data["json"] = randCode
// 	c.ServeJSON()
// 	c.Ctx.ResponseWriter.WriteHeader(200)
// }

// @Title NewPassword
// @Description  通过前面忘记密码的过程后，设置新的密码
// @Param userInfo body models.User true 用户信息(需要的是用户ID，新密码）
// @Success 200 Update successfully
// @Failure 404 数据库无此用户
// @Failure 400 解析body失败
// @Failure 406 更新密码失败
// @router /forgetPw/New [post]
func (c *UserController) NewPW() {
	var userInfo struct{
		Id string
		NewPassword string
		Verify string
	}
	body := c.Ctx.Input.RequestBody
	if err := json.Unmarshal(body, &userInfo); err != nil {
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	sess := c.GetSession("verify")
	if sess == nil {
		models.Log.Error("set new password without being verified")
		var response struct {
			Msg string `json:"msg"`
		}
		response.Msg = "no verification"
		c.Data["json"] = response
		c.ServeJSON()
		c.Ctx.ResponseWriter.WriteHeader(406) //没有点击验证码
		return
	}
	vCode := sess.(string)
	fmt.Println(vCode,userInfo.Verify)
	if vCode != userInfo.Verify {
		models.Log.Error("verify fail")
		var response struct {
			Msg string `json:"msg"`
		}
		response.Msg = "wrong verify code"
		c.Data["json"] = response
		c.ServeJSON()
		c.Ctx.ResponseWriter.WriteHeader(403) //没有点击验证码
		return
	}
	o := orm.NewOrm()
	usr := models.User{Id: userInfo.Id}
	//fmt.Println(usr)
	// 查询记录
	if err := o.Read(&usr); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) // 查不到id对应的用户
		return
	}
	//查询成功，更新密码
	usr.Password = userInfo.NewPassword
	if _, err := o.Update(&usr, "password"); err != nil {
		models.Log.Error("update error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(500) // 更新数据失败
		return
	}
	c.DelSession("verify")
	//根据旧的userInfo删除旧session
	c.DelSession("userInfo")
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
