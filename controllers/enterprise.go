package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"

	//"encoding/json"
	_"fmt"
	"github.com/astaxie/beego"
	//"github.com/astaxie/beego/orm"
	//"github.com/gfbankend/models"
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