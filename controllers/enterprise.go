package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"

	//"encoding/json"
	_"fmt"
	"github.com/astaxie/beego"
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
// @Title show all card type
// @Description 显示所有优惠政策
// @Param id	path	string	true 商家ID
// @Success 200  
// @Failure 404 Fail to read enterpriseId
// @router enterprise/:id [get]
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
// @Param EnterPriseInfo body models.Enterprise true 注册信息
// @Success 200 {object} models.Enterprise "OK"
// @Failure 400 解析body错误
// @Failure 406 账号信息格式有误
// @Failure 403 数据库插入错误
// @router enterprise/enroll [post]
func (c *UserController) EnterpriseEnroll() {
	  body := c.Ctx.Input.RequestBody
	  var enterprise models.Enterprise
	  if err := json.Unmarshal(body, &enterprise); err != nil {
	  	models.Log.Error("Enterprise enroll: wrong json")
	  	c.Ctx.ResponseWriter.WriteHeader(400)
		  return
	  }
	  // parse to get id
	  if err := enterprise.EnterpriseParse(); err != nil {
		  models.Log.Error("Enterprise enroll: fail to parse")
		  c.Ctx.ResponseWriter.WriteHeader(406)
		  return
	  }
	  o := orm.NewOrm()
	  if _, err := o.Insert(&enterprise); err != nil {
	  	models.Log.Error("Enterprise enroll: fail to insert")
	  	c.Ctx.ResponseWriter.WriteHeader(406)
	  	return
	  }
	  c.Data["json"] = enterprise
	  c.ServeJSON()
}

// @author: zyj
// @Title Login
// @Description 商家登陆
// @Param enterpriseInfo body true account(string)+password(string)+accountType(string)为mail或者phone
// @Success 200 {object} models.User Register successfully
// @Failure 406 数据库查询报错，可能用户所填账号或密码错误
// @Failure 400 信息内容或格式有误
// @router enterprise/login [put]
func (c *UserController) EnterpriseLogin() {
	
}

// @author: zyj
// @Title changePW
// @Description change password
// @Param userInfo body models.User true 用户信息(需要的是用户ID，新密码）
// @Success 200 Update successfully
// @Failure 404 数据库无此用户
// @Failure 400 解析body失败
// @Failure 406 更新密码失败
// @router Enterprise/password [put]
func (c *UserController) EnterpriseChangePW() {
	 
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
// @Param enterpriseInfo body models.Enterprise true 用户信息(需要的是用户ID，新密码）
// @Success 200 Update successfully
// @Failure 404 数据库无此用户
// @Failure 400 解析body失败
// @Failure 406 更新密码失败
// @router Enterprise/ForgetPW/New [put]
func (c *UserController) EnterpriseNewPW() {

}

// @author:zjn
// @Title enterprise information modify
// @Description  修改注册的商家注册信息
// @Param enterpriseInfo body models.Enterprise true 重新提交的商家注册信息
// @Success 200 Update成功
// @Failure 404 数据库无此商铺
// @Failure 400 解析body失败
// @Failure 406 更新商铺信息失败
// @router Enterprise/infomodify [put]
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