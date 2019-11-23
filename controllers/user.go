package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"
	"github.com/go-gomail/gomail"
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
	}
	c.Ctx.ResponseWriter.WriteHeader(200) //成功读取所有卡片
	c.Data["json"] = cardList
	//发送json
	c.ServeJSON()
}

//YZY，返回用户资料
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


func (c *UserController) UpAvatar() {
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", c.Ctx.Request.Header.Get("Origin"))
    userId:=c.GetString("id")
	tmpFile, fHeader, err  := c.Ctx.Request.FormFile("avatar")
	if err != nil{
		models.Log.Error("read error", err) //读取用户卡片信息失败
		c.Ctx.ResponseWriter.WriteHeader(400)
	}
	savePath := "D:/" +userId+".jpg"
	//savePath := "/root/gfbankend/User/avatar/" +userId+".jpg"        //设置保存路径
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
