package controllers

import (
	"encoding/json"
	_ "fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"
	_ "github.com/pkg/errors"
	"os"
)

type EnterpriseController struct {
	beego.Controller
}

// @author: zjn
// @Title show all card type
// @Description 显示所有优惠政策
// @Param id	path	string	true 商家ID
// @Success 200  
// @Failure 404 Fail to read enterpriseId
// @router /enterprise/:id [get]
func (c *UserController) AllCardDemo() {
	//查看session的操作
	//if c.GetSession("userInfo") == nil {
	//	models.Log.Error("no login")
	//	c.Ctx.ResponseWriter.WriteHeader(401)
	//	return
	//}
	// 取得用户ID from path
	id := c.Ctx.Input.Param(":id")
	//储存所有卡片类型信息
	var carddemoList []models.CardDemo
	//使用orm接口查询相关信息
	o := orm.NewOrm()
	qt := o.QueryTable("carddemo")
	//取出carddemo表中所有信息，放入carddemoList中
	_, err := qt.Filter("enterprise__exact", id).All(&carddemoList)
	if err != nil || len(carddemoList) == 0 {
		models.Log.Error("read error", err)
		c.Ctx.ResponseWriter.WriteHeader(404)
		return
	}
	//使用json格式传输所有信息
	c.Data["json"] = carddemoList
	//发送json
	c.ServeJSON()
	c.Ctx.ResponseWriter.WriteHeader(200)
}

// @author: ml
// @Title Register
// @Description  商家注册
// @Param EnterPriseInfo body models.Enterprise true 注册信息(商家名称name+营业执照号码license_id+地址Addr+商家类型type+是否为本地is_local+管理人名称manager_name+管理人身份证manager_id+手机号码phone+密码password
// @Success 200 {object} models.Enterprise "OK"
// @Failure 400 信息有误
// @Failure 406 数据库加入错误
// @router /enterprise/enroll [post]
func (c *UserController) EnterpriseEnroll() {
	var Request struct {
		// Enterprise Info
		Name      string `json:"name"`
		LicenseId string `json:"license_id"`
		Addr      string `json:"addr"`
		Type      string `json:"type"`
		IsLocal   bool   `json:"is_local"`
		// Manager Info
		ManagerName string `json:"manager_name"`
		ManagerID   string `json:"manager_id"`
		Phone       string `json:"phone"`
		Password    string `json:"password"`
	}
	body := c.Ctx.Input.RequestBody
	if err := json.Unmarshal(body, &Request); err != nil {
		models.Log.Error("Enterprise enroll: wrong json")
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	enterprise := models.Enterprise{
		Name:      Request.Name,
		LicenseId: Request.LicenseId,
		Addr:      Request.Addr,
		Type:      Request.Type,
		IsLocal:   Request.IsLocal,
	}
	manager := models.Manager{
		Enterprise: Request.Name,
		Name:       Request.ManagerName,
		ID:         Request.ManagerID,
		Phone:      Request.Phone,
		Password:   Request.Password,
	}
	// parse to get id
	if err := enterprise.EnterpriseParse(); err != nil {
		models.Log.Error("Enterprise enroll: fail to parse", err)
		c.Ctx.ResponseWriter.WriteHeader(406)
		return
	}
	o := orm.NewOrm()
	if _, err := o.Insert(&manager); err != nil {
		models.Log.Error("Enterprise enroll: fail to insert", err)
		c.Ctx.ResponseWriter.WriteHeader(406)
		return
	}
	if _, err := o.Insert(&enterprise); err != nil {
		// 防止出现管理员插入成功，而商家插入失败
		_, _ = o.Delete(&manager)
		models.Log.Error("Enterprise enroll: fail to insert", err)
		c.Ctx.ResponseWriter.WriteHeader(406)
		return
	}
	var Response struct {
		Enterprise models.Enterprise `json:"enterprise"`
		Manager    models.Manager    `json:"manager"`
	}
	Response.Manager = manager
	Response.Enterprise = enterprise
	c.Data["json"] = Response
	c.ServeJSON()
}

// @author: zyj
// @Title Login
// @Description 商家登陆
// @Param enterpriseInfo body true account(string)+password(string)+remember(bool)
// @Success 200 {object} models.User Register successfully
// @Failure 406 数据库查询报错，可能用户所填账号或密码错误
// @Failure 400 信息内容或格式有误
// @router /enterprise/login [put]
// 要返回管理员管理的企业的信息
func (c *EnterpriseController) EnterpriseLogin() {
	o := orm.NewOrm()
	manager := models.Manager{}
	body := c.Ctx.Input.RequestBody
	var eInfo struct {
		Account  string `json:"account"`
		Password string `json:"password"`
		Remember bool   `json:"remember"`
	}
	if err := json.Unmarshal(body, &eInfo); err != nil {
		models.Log.Error("Unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	manager.ID = eInfo.Account
	manager.Password = eInfo.Password
	if err := o.Read(&manager, "id", "password"); err != nil {
		models.Log.Error("login error: auth fail")
		c.Ctx.ResponseWriter.WriteHeader(406)
		return
	}
	c.Data["json"] = manager
	// 如果需要记住账号密码
	if eInfo.Remember == true {
		c.Ctx.SetSecureCookie("miller", "account", eInfo.Account)
		c.Ctx.SetSecureCookie("miller", "password", eInfo.Password)
		c.Ctx.SetCookie("remember", "true")
	}
	// originHeader := c.Ctx.Input.Header("Origin")
	// c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	c.SetSession("managerInfo", manager) // 登录成功，设置session
	c.ServeJSON()                        // 传用户对象给前端
}

// @author: zyj
// @Title Check
// @Description Check across
// @Success 200
// @router /enterprise/login [options]
//func (c *EnterpriseController) CheckAcross() {
//	c.Ctx.ResponseWriter.WriteHeader(200)
//}

// @author: zyj
// @Title changePW
// @Description change password
// @Param enterpriseInfo body  true 用户信息(需要的是用户ID,原密码,新密码）
// @Success 200 Update successfully
// @Failure 404 数据库无此用户
// @Failure 400 解析body失败
// @Failure 500 数据库更新密码失败
// @router /enterprise/password [put]
func (c *EnterpriseController) EnterpriseChangePW() {
	if c.GetSession("managerInfo") == nil {
		models.Log.Error("no login")
		c.Ctx.ResponseWriter.WriteHeader(401)
		return
	}
	var manager struct {
		Id          string `json:"id"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	body := c.Ctx.Input.RequestBody
	if err := json.Unmarshal(body, &manager); err != nil {
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	o := orm.NewOrm()
	man := models.Manager{ID: manager.Id}
	if err := o.Read(&man); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) // 查不到id对应的用户
		return
	}
	//验证用户输入的密码是否与旧密码一致
	if manager.OldPassword != man.Password {
		models.Log.Error("wrong old password: ")
		c.Ctx.ResponseWriter.WriteHeader(403) // 原密码错误
		return
	}
	man.Password = manager.NewPassword
	if _, err := o.Update(&man); err != nil {
		models.Log.Error("update error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(500) // 更新数据失败
		return
	}
	//根据旧的userInfo删除对应session
	c.DelSession("managerInfo")
	c.Ctx.ResponseWriter.WriteHeader(200) // 更新成功
}

// @author: lj
// @Title ForgetPW
// @Description Forget password
// @Param enterpriseInfo body models.Enterprise true 用户信息(需要的是用户ID，Phone）
// @Success 200 successfully
// @Failure 400 解析body失败
// @Failure 404 ID错误
// @Failure 405 Phone错误
// @router /enterprise/password [post]
func (c *UserController) EnterpriseForgetPW() {
	var Request struct {
		ID    string	`json:"id"`
		Phone string	`json:"phone"`
	}
	body := c.Ctx.Input.RequestBody
	if err := json.Unmarshal(body, &Request); err != nil {
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	manager := models.Manager{ID: Request.ID, Phone: Request.Phone}
	o := orm.NewOrm()
	if err := o.Read(&manager); err != nil {
		models.Log.Error("NewPW: fail to read", err)
		c.Ctx.ResponseWriter.WriteHeader(404) //查找不到对应的ID
		return
	}
	c.Data["json"] = manager
	c.ServeJSON()
}

// @author: ml
// @Title NewPassword
// @Description  通过前面忘记密码的过程后，设置新的密码
// @Param enterpriseInfo body models.Enterprise true 用户信息(需要的是用户ID，新密码）
// @Success 200 Update successfully
// @Failure 404 数据库无此用户
// @Failure 400 解析body失败
// @Failure 406 更新密码失败
// @router /enterprise/password/new [put]
func (c *UserController) EnterpriseNewPW() {
	body := c.Ctx.Input.RequestBody
	var Request struct {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(body, &Request); err != nil {
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	manager := models.Manager{Phone: Request.Phone}
	o := orm.NewOrm()
	if err := o.Read(manager, "phone"); err != nil {
		models.Log.Error("NewPW: fail to read", err)
		c.Ctx.ResponseWriter.WriteHeader(404)
		return
	}
	manager.Password = Request.Password
	if _, err := o.Update(&manager); err != nil {
		models.Log.Error("NewPW: fail to update", err)
		c.Ctx.ResponseWriter.WriteHeader(406)
		return
	}
	c.Data["json"] = manager
	c.ServeJSON()
}

// zyj
// ***************
// @author:zjn
// @Title enterprise information modify
// @Description  修改注册的商家注册信息
// @Param enterpriseInfo body  重新提交的商家注册信息,body内部包含两个部分信息,一个为manager信息,一个为enterprise信息
// @Success 200 Update Successfully
// @Failure 404 数据库无此商铺
// @Failure 400 解析body失败
// @Failure 406 更新商铺信息失败
// @Failure 407 创建base64文件失败
// @Failure 408 关闭base64文件失败
// @router /enterprise/modifyInfo [put]
// 加上修改管理员信息
func (c *UserController) EnterpriseInfoModify() {
	body := c.Ctx.Input.RequestBody
	var newInfo struct {
		Enterprise models.Enterprise `json:"enterprise"`
		Manager    models.Manager    `json:"manager"`
		Base64     string            `json:"base64"`
	}
	if err := json.Unmarshal(body, &newInfo); err != nil {
		models.Log.Error("wrong json")
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	enterprise := newInfo.Enterprise
	enterpriseId := newInfo.Enterprise.Id
	manager := newInfo.Manager
	base64 := newInfo.Base64
	path := "../static/base64/" + enterpriseId + ".txt"
	if f,err := os.Create(path); err != nil {
		models.Log.Error(" fail to create the file",err)
		c.Ctx.ResponseWriter.WriteHeader(407)
		return
	} else {
		content := []byte(base64)
		if _,err := f.Write(content); err != nil {
			models.Log.Error(" fail to write base64 to the file",err)
			c.Ctx.ResponseWriter.WriteHeader(407)
			return
		}
		if err := f.Close() ; err != nil {
			models.Log.Error(" fail to close the file",err)
			c.Ctx.ResponseWriter.WriteHeader(408)
			return
		}
	}
	o := orm.NewOrm()
	if _, err := o.Update(&enterprise); err != nil {
		models.Log.Error(" fail to update enterprise",err)
		c.Ctx.ResponseWriter.WriteHeader(406)
		return
	}
	if _, err := o.Update(&manager); err != nil {
		models.Log.Error(" fail to update manager")
		c.Ctx.ResponseWriter.WriteHeader(406)
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200)
}

// zyj
// @author:
// @Title return the enterprise information
// @Description  返回商家具体信息以及管理员信息 打包在一个json内，分别用enterprise和managerList两个key取得
// @Param enterpriseId path string true 商家ID
// @Success 200 Return Successfully
// @Failure 404 数据库无此商家
// @Failure 404 读取管理员失败
// @Failure 400 解析body失败
// @router /enterprise/:id [get]
func (c *UserController) EnterpriseInfo(){
	eid := c.Ctx.Input.Param(":id")
	var enterprise models.Enterprise
	enterprise.Id = eid
	o := orm.NewOrm()
	if err := o.Read(enterprise); err != nil {
		models.Log.Error("read enterprise error ", err)
		c.Ctx.ResponseWriter.WriteHeader(404)
		return
	}
	var managerList []models.Manager
	qt := o.QueryTable("manager")
	_, err := qt.Filter("enterprise__exact", eid).All(&managerList)
	if err != nil {
		models.Log.Error("read manager error", err) //读取用户卡片信息失败
		c.Ctx.ResponseWriter.WriteHeader(405)
		return
	}
	var ret struct{
		Enterprise models.Enterprise `json:"enterprise"`
		ManagerList []models.Manager `json:"managerList"`
	}
	ret.Enterprise = enterprise
	ret.ManagerList = managerList
	c.Data["json"] = ret
	c.ServeJSON()
}

// zyj
// @author:
// @Title put New Card
// @Description  发布新的卡片
// @Param  cardDemo body models.CardDemo true 卡片基本信息
// @Param  backgroundBase64 body string true Base64Code
// @Success 200 put successfully
// @Failure 400 解析body失败
// @Failure 405 插入数据失败
// @router /enterprise/card [put]
func (c *UserController) EnterpriseNewDemo() {
	body := c.Ctx.Input.RequestBody
	var cardInfo struct {
		CardDemo   models.CardDemo   `json:"cardDemo"`
		Base64     string            `json:"backgroundBase64"`
	}
	if err := json.Unmarshal(body, &cardInfo); err != nil {
		models.Log.Error("wrong json")
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	o := orm.NewOrm()
	if _,err := o.Insert(&cardInfo.CardDemo); err!= nil {
		models.Log.Error("insert database error: %s",err)
		c.Ctx.ResponseWriter.WriteHeader(405) //数据库更新失败
		return
	}
	path := "../static/base64/" + string(cardInfo.CardDemo.ID) + ".txt"
	if f,err := os.Create(path); err != nil {
		models.Log.Error(" fail to create the file",err)
		c.Ctx.ResponseWriter.WriteHeader(407)
		return
	} else {
		content := []byte(cardInfo.Base64)
		if _,err := f.Write(content); err != nil {
			models.Log.Error(" fail to write base64 to the file",err)
			c.Ctx.ResponseWriter.WriteHeader(407)
			return
		}
		if err := f.Close() ; err != nil {
			models.Log.Error(" fail to close the file",err)
			c.Ctx.ResponseWriter.WriteHeader(408)
			return
		}
	}
	c.Ctx.ResponseWriter.WriteHeader(200)
}

// // @author:
// // @Title NewPassword
// // @Description  发布新的卡片
// // @Param  cardInfo body models.CardDemo true 用户类型,数量
// // @Success 200 Update successfully
// // @Failure 404 数据库无此用户
// // @Failure 400 解析body失败
// // @Failure 406 更新密码失败
// // @router /enterprise/card [put]
// // 前端给予活动的标题、卡片的类型、活动的内容、展示在卡片详细信息页面的背景图、活动描述，有效期，后端增加CardDemo
// func (c *UserController) EnterpriseNewCard() {

// }

// @author: zjn
// @Title addUser
// @Description  商家增加一个某张已发售卡片的用户
// @Param  
// @Success 
// @Failure 
// @router 
func (c *UserController) AddUser() {

}

// @author: zjn
// @Title deleteUser
// @Description  商家删除一个某张已发售卡片的用户
// @Param  
// @Success 
// @Failure 
// @router 
func (c *UserController) DeleteUser() {

}

// @author: zjn
// @Title readUser
// @Description  商家查询某张已发售卡片的用户
// @Param  
// @Success 
// @Failure 
// @router 
// 前端给卡的类型，后端根据卡的类型到card表单里面找该种卡的所有userId，积分，以及拥有卡的时间
func (c *UserController) ReadUser() {

}

// ml
// @Title readActivity
// @Description  查询活动
// @Param enterprise body string true 企业名称
// @Param card_type body string true 卡片类型
// @Success 200 
// @Failure 400 请求格式出错
// @Failure 503 读取数据库发生错误（服务器端可能有问题）
// @router /enterprise/activity [put]
func (c *UserController) ReadActivity() {
	body := c.Ctx.Input.RequestBody
	// 通过商家、卡片类型查询相关活动
	var Req struct {
		Enterprise string `json:"enterprise"`
		CardType   string `json:"card_type"`
	}
	if err := json.Unmarshal(body, &Req); err != nil {
		models.Log.Error("readActivity unmarshall error:", err)
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	var Activities []models.CardDemo
	qt := orm.NewOrm().QueryTable("card_demo")
	cond := orm.NewCondition()
	//  查询企业名称和类型都匹配的数据
	cond1 := cond.And("enterprise__iexact", Req.Enterprise).And("card_type__iexact", Req.CardType)
	if _, err := qt.SetCond(cond1).All(&Activities);
		err != nil {
		models.Log.Error("readActivity query error:", err)
		c.Ctx.ResponseWriter.WriteHeader(404)
		return
	}
	// 设置回复
	var Resp struct {
		Activity interface{} `json:"activity"`
	}
	Resp.Activity = Activities
	c.Data["json"] = Resp
	c.ServeJSON()
}

// @author:ml
// @Title readAllActivities
// @Description  查看所有优惠活动
// @Success 200 请求成功，返回所以活动
// @Failure 503 读取数据库出错(可能服务器端数据库出错)
// @router /enterprise/activity [get]
// 返回所有cardDemo
func (c *UserController) ReadAllActivities() {
	var Resp struct {
		Activities	[]models.CardDemo	`json:"activities"`
	}
	qt := orm.NewOrm().QueryTable("card_demo")
	if _, err := qt.All(&Resp.Activities); err != nil {
		models.Log.Error("ReadAllActivities query error:",err)
		c.Ctx.ResponseWriter.WriteHeader(503)
		return
	}
	c.Data["json"] = Resp
	c.ServeJSON()
}
