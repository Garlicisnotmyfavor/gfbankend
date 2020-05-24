package controllers

import (
	"encoding/json"
	_ "fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"
	_ "github.com/pkg/errors"
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
func (c *UserController) AllCarddemo() {
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
		Manager models.Manager `json:"manager"`
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
func (c *EnterpriseController) EnterpriseLogin() {
	o := orm.NewOrm()
	manager := models.Manager{}
	body := c.Ctx.Input.RequestBody
	var eInfo struct {
		Account    string	`json:"account"`
		Password   string	`json:"password"`
		Remember   bool		`json:"remember"`
	}
	if err := json.Unmarshal(body,&eInfo); err != nil {
		models.Log.Error("Unmarshal error: ",err)
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	manager.ID = eInfo.Account
	manager.Password = eInfo.Password
	if err:= o.Read(&manager,"id","password");err!=nil{
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
	c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	c.SetSession("managerInfo", manager) // 登录成功，设置session
	c.ServeJSON()                  // 传用户对象给前端
}
// @author: zyj
// @Title Login
// @Description 商家登陆
// @Param enterpriseInfo body true account(string)+password(string)+accountType(string)为mail或者phone
// @Success 200 {object} models.User Register successfully
// @Failure 406 数据库查询报错，可能用户所填账号或密码错误
// @Failure 400 信息内容或格式有误
// @router /enterprise/login [options]
func (c *EnterpriseController) checkAcross() {
	c.Ctx.ResponseWriter.WriteHeader(200)
}
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
		Id          string
		OldPassword string
		NewPassword string
	}
	body := c.Ctx.Input.RequestBody
	if err := json.Unmarshal(body,&manager); err!= nil{
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	o := orm.NewOrm()
	man := models.Manager{ID:manager.Id}
	if err := o.Read(&man);err!=nil{
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
// @router /enterprise/forgetpw [post]
func (c *UserController) EnterpriseForgetPW() {
	var Request struct {
		ID string
		Phone string
	}
	body := c.Ctx.Input.RequestBody
	if err := json.Unmarshal(body,&Request); err!= nil {
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	manager := models.Manager{ID: Request.ID, Phone:Request.Phone}
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
// @router /enterprise/forgetpw/new [put]
func (c *UserController) EnterpriseNewPW() {
	body := c.Ctx.Input.RequestBody
	var Request struct {
		Phone string
		Password string
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
	if _,err := o.Update(&manager); err != nil {
		models.Log.Error("NewPW: fail to update", err)
		c.Ctx.ResponseWriter.WriteHeader(406)
		return
	}
	c.Data["json"] = manager
	c.ServeJSON()
}

// @author:zjn
// @Title enterprise information modify
// @Description  修改注册的商家注册信息
// @Param enterpriseInfo body models.Enterprise true 重新提交的商家注册信息
// @Success 200 Update成功
// @Failure 404 数据库无此商铺
// @Failure 400 解析body失败
// @Failure 406 更新商铺信息失败
// @router /enterprise/infomodify [put]
func (c *UserController) EnterpriseInfomodify() {
	body := c.Ctx.Input.RequestBody
	var enterprise models.Enterprise
	var newenterprise models.Enterprise
	if err := json.Unmarshal(body, &enterprise); err != nil {
		models.Log.Error("Enterprise enroll: wrong json")
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	o := orm.NewOrm()
	if err := o.Read(&newenterprise); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) // 查不到对于商铺的信息
		return
	}
	if _, err := o.Update(&enterprise); err != nil {
		models.Log.Error("Enterprise enroll: fail to update")
		c.Ctx.ResponseWriter.WriteHeader(406)
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200)
}

// @author:
// @Title NewPassword
// @Description  发布新的优惠政策
// @Param userInfo body models.User true 用户信息(需要的是用户ID，新密码）
// @Success 200 Update successfully
// @Failure 404 数据库无此用户
// @Failure 400 解析body失败
// @Failure 406 更新密码失败
// @router /enterprise/newdemo [put]
func (c *UserController) EnterpriseNewDemo() {

}

// @author:
// @Title NewPassword
// @Description  发布新的卡片
// @Param  cardInfo body models.CardDemo true 用户类型,数量
// @Success 200 Update successfully
// @Failure 404 数据库无此用户
// @Failure 400 解析body失败
// @Failure 406 更新密码失败
// @router /enterprise/NewCard [put]
func (c *UserController) EnterpriseNewCard() {

}
