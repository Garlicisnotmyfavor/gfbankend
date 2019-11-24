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
	"io/ioutil"
	"path"
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

// @Title Get
// @Description get current user's profile
// @Param id models.User.id  true
// @Success 200 get successfully
// @Failure 404 Fail to find picture
// @router / [get]
func (c *UserController) Get() {
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", c.Ctx.Request.Header.Get("Origin"))
	userId := c.Ctx.Input.Param(":id") //获取需要上传的文件文件名
	filename:=userId+".jpg"
	//查看是否存在需要的图片
 //   readPath :="D:/"
    readPath := "/root/gfbankend/User/avatar/"
	img:= path.Join(readPath,filename)
	c.Ctx.Output.Header("Content-Type", "image/jpg")
	c.Ctx.Output.Header("Content-Disposition",fmt.Sprintf("inline; filename=\"%s\"",img))
	file, err := ioutil.ReadFile(img)
	if err != nil {
		models.Log.Error("read error", err) //未找到对应图片
		c.Ctx.ResponseWriter.WriteHeader(404)
		return
	}
	//c.Ctx.WriteString(string(file))
	c.Data["avatar"]=string(file)
}

// @Title UpAvatar
// @Description upload avatar
// @Param id  models.User.Id avatar file
// @Success 200 upload successfully
// @Failure 500 Fail to save picture
// @Failure 502 Fail to close uploading file
// @router /avatar UpAvatar[post]
func (c *UserController) UpAvatar() {
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", c.Ctx.Request.Header.Get("Origin"))
    userId:=c.GetString("id")
	tmpFile, fHeader, err  := c.Ctx.Request.FormFile("avatar")
	if err != nil{
		models.Log.Error("read error", err) //读取用户卡片信息失败
		c.Ctx.ResponseWriter.WriteHeader(400)
	}
	//savePath := "D:/" +userId+".jpg"
	savePath := "/root/gfbankend/User/avatar/" +userId+".jpg"        //设置保存路径
	beego.Info("Header:", fHeader.Header)     //map[Content-Disposition:[form-data; name="123"; filename="upimage.jpg"] Content-Type:[image/jpeg]]
	beego.Info("Size:", fHeader.Size)         //114353
	beego.Info("Filename:", fHeader.Filename) //upimage.jpg
	if err=c.SaveToFile("avatar", savePath);err !=nil{
		models.Log.Error("save error", err) //存储图片失败
		c.Ctx.ResponseWriter.WriteHeader(500)
	}
	if err:=tmpFile.Close();err!=nil {
		models.Log.Error("close error", err) //存储图片失败
		c.Ctx.ResponseWriter.WriteHeader(502)
	}                   //关闭上传的文件，不然的话会出现临时文件不能清除的情况
}

//ML，用户注册
// @Title Register
// @Description user register
// @Param user body models.User UserInfo true
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
// @Failure 401 Fail to login
// @Failure 400 illegal account form
// @router /login [put]
func (c *UserController) Put() {
	o := orm.NewOrm()
	user := models.User{}

	body:=c.Ctx.Input.RequestBody
	var uInfo map[string]string

	//解析前端JSON数据获得账号密码
	if err:=json.Unmarshal(body,&uInfo);err!=nil{
		models.Log.Error("Unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}

	account:=[]byte(uInfo["account"])
	user.Password=uInfo["password"]

	//正则表达式匹配模式
	pattern1:="^[0-9]+$" //匹配手机号
	pattern2 := "^[a-z0-9A-Z]+[- | a-z0-9A-Z . _]+@([a-z0-9A-Z]+(-[a-z0-9A-Z]+)?.)+[a-z]{2,}$" //匹配邮箱

	//判断是用户使用的是邮箱还是手机
	isPhone,_:=regexp.Match(pattern1,account)
	isMail,_:=regexp.Match(pattern2,account)
	if isPhone{
		user.Tel=string(account)
	}else if isMail{
		user.Mail=string(account)
	}else{
		models.Log.Error("illegal account")
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}

	err1:=o.Read(&user,"mail","password")
	err2:=o.Read(&user,"tel","password")
	//用户信息错误
	if err1!=nil&&err2!=nil{
		models.Log.Error("read error",err1)
		c.Ctx.ResponseWriter.WriteHeader(401)//登录失败
		return
	}
	//信息匹配登录成功
	c.Ctx.ResponseWriter.WriteHeader(200)
}

//ZJN，显示所有被删卡片
func (c *UserController) GetDel() {

}

//ZJN，恢复指定卡片
func (c *UserController) RecoverDel() {

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
// @Title Feedback
// @Description send feedback mail
// @Param    body        body         true
// @Success 200 Update successfully
//@Failure 500 Fail to send mail
// @router /password [post]
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
