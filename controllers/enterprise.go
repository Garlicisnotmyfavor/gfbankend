package controllers

import (
	"encoding/json"
	_"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"
	_ "github.com/pkg/errors"
	//util "github.com/gfbankend/utils"
	//"strconv"
	//"strings"
	//"time"
)

type EnterpriseController struct {
	beego.Controller
}

// @author: zjn
// @Title 
// @Description 显示所有优惠政策
// @Param eid path string true 商家ID
// @Success 200  
// @Failure 404 Fail to read
// @router enterprise/:id [get]
func (c *UserController) AllCarddemo() {
	 
}

// @author: ml
// @Title Register
// @Description  商家注册
// @Param EnterPriseInfo body models.Enterprise true 注册信息
// @Success 200 {object} models.User "OK"
// @Failure 400 解析body错误
// @Failure 406 账号信息格式有误
// @Failure 403 数据库插入错误
// @router enterprise/enroll [post]
func (c *UserController) EnterpriseEnroll() {
	  
}

// @author: zyj
// @Title Login
// @Description 商家登陆
// @Param enterpriseInfo body true account(string)+password(string)+accountType(string)为mail或者phone
// @Success 200 {object} models.User Register successfully
// @Failure 406 数据库查询报错，可能用户所填账号或密码错误
// @Failure 400 信息内容或格式有误
// @router enterprise/login [put]
func (c *EnterpriseController) EnterpriseLogin() {
	o := orm.NewOrm()
	enterprise := models.Enterprise{}
	body := c.Ctx.Input.RequestBody
	var eInfo struct {
		Account    string
		Password   string
		Remember   bool
	}
	if err := json.Unmarshal(body,&eInfo); err != nil {
		models.Log.Error("Unmarshal error: ",err)
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	enterprise.Name = eInfo.Account
	enterprise.Password = eInfo.Password
	if err:= o.Read(&enterprise,"name","password");err!=nil{
		models.Log.Error("login error: auth fail")
		c.Ctx.ResponseWriter.WriteHeader(406)
		return
	}
	c.Data["json"] = enterprise
	// 如果需要记住账号密码
	if eInfo.Remember == true {
		c.Ctx.SetSecureCookie("miller", "account", eInfo.Account)
		c.Ctx.SetSecureCookie("miller", "password", eInfo.Password)
		c.Ctx.SetCookie("remember", "true")
	}
	c.SetSession("enterpriseInfo", enterprise) // 登录成功，设置session
	c.ServeJSON()                  // 传用户对象给前端

}

// @author: zyj
// @Title changePW
// @Description change password
// @Param userInfo body / true 用户信息(需要的是用户ID,原密码,新密码）
// @Success 200 Update successfully
// @Failure 404 数据库无此用户
// @Failure 400 解析body失败
// @Failure 500 数据库更新密码失败
// @router Enterprise/password [put]
func (c *EnterpriseController) EnterpriseChangePW() {
	if c.GetSession("enterpriseInfo") == nil {
		models.Log.Error("no login")
		c.Ctx.ResponseWriter.WriteHeader(401)
		return
	}
	var enterprise struct {
		Id          string
		OldPassword string
		NewPassword string
	}
	body := c.Ctx.Input.RequestBody
	if err := json.Unmarshal(body,&enterprise); err!= nil{
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	o := orm.NewOrm()
	enpri := models.Enterprise{Id:enterprise.Id}
	if err := o.Read(&enpri);err!=nil{
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) // 查不到id对应的用户
		return
	}
	//验证用户输入的密码是否与旧密码一致
	if enterprise.OldPassword != enpri.Password {
		models.Log.Error("wrong old password: ")
		c.Ctx.ResponseWriter.WriteHeader(403) // 原密码错误
		return
	}
	enpri.Password = enterprise.NewPassword
	if _, err := o.Update(&enpri); err != nil {
		models.Log.Error("update error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(500) // 更新数据失败
		return
	}
	//根据旧的userInfo删除对应session
	c.DelSession("enterpriseInfo")
	c.Ctx.ResponseWriter.WriteHeader(200) // 更新成功
}

// @author: lj
// @Title ForgetPW
// @Description Forget password
// @Param userInfo body models.User true 用户信息(需要的是用户ID，邮件）
// @Success 200 successfully
// @Failure 404 数据库无此用户
// @Failure 400 解析body失败
// @router Enterprise/forgetPw [post]
func (c *UserController) EnterpriseForgetPW() {
	 
}

// @author: ml
// @Title NewPassword
// @Description  通过前面忘记密码的过程后，设置新的密码
// @Param userInfo body models.User true 用户信息(需要的是用户ID，新密码）
// @Success 200 Update successfully
// @Failure 404 数据库无此用户
// @Failure 400 解析body失败
// @Failure 406 更新密码失败
// @router Enterprise/ForgetPW/New [put]
func (c *UserController) EnterpriseNewPW() {
	 
}

// @author:zjn
// @Title
// @Description  修改注册的商家信息
// @Param userInfo body models.User true 用户信息(需要的是用户ID，新密码）
// @Success 200 Update successfully
// @Failure 404 数据库无此用户
// @Failure 400 解析body失败
// @Failure 406 更新密码失败
// @router Enterprise/infomodify [put]
func (c *UserController) EnterpriseInfomodify() {
	 
}

// @author:
// @Title NewPassword
// @Description  发布新的优惠政策
// @Param userInfo body models.User true 用户信息(需要的是用户ID，新密码）
// @Success 200 Update successfully
// @Failure 404 数据库无此用户
// @Failure 400 解析body失败
// @Failure 406 更新密码失败
// @router Enterprise/newdemo [put]
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
// @router Enterprise/NewCard [put]
func (c *UserController) EnterpriseNewCard() {
	 
}