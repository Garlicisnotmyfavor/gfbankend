package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"
	"github.com/pkg/errors"
	"time"
)

type CardController struct {
	beego.Controller
}

//查询指定card_id对应的卡片的所有信息
//zjn
//@Title give_card_all_info
//@Description 将这张卡片的所有信息传出去
//@Param	id	query	string	true	原本的卡号  我还不清楚这里的修改
//@Success 200	{object} models.Card {object} model.Strategy	修改成功，返回新卡片对象和策略对象
//@Failure 400	查询不到对应的卡
//@Failure 404	查询不到对应的优惠策略
//@router  /card/:id	[get]
func (c *CardController) Get_cardidinfo() {
	// 获取路由参数
	id := c.Ctx.Input.Param(":id")
	o := orm.NewOrm()
	//设置一个填充了cardid的card结构
	card := models.Card{CardId: id}
	// 查询记录
	if err := o.Read(&card); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) // 查不到id对应的卡
		return
	}
	//若查到card这一列后，需要找到它的卡的积分或卷的规则，但目前只针对一个策略
	//后面需要用匹配的方式，解析出多个策略
	stra := models.StrategyTable{Strategy: card.Strategy}
	if err := o.Read(&stra); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) // 查不到卡对应的优惠策略
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200) //成功
	
	//怎么传两个字段？ stra
	c.Data["json"] = card
	c.ServeJSON()
}

//添加卡片 在user表里添加此user和card的关联
//zjn
//@Title addCard
//@Description 将这个user的id和卡绑定
//@Param	id	query	string	true	原本的卡号
//@Success 200	{object} models.Card 	返回绑定的卡的大致信息
//@Failure 403	绑定的卡片不存在
//@router  /card/:id/add [get]
func (c *CardController) addCard() {
	//这里没有对比enterprise和cardid
	id := c.Ctx.Input.Param(":id")
	var card models.Card
	card.CardId = id
	//目前的逻辑不需要解析函数
	//if err := card.CardParse(); err != nil {
	//	models.Log.Error("card parse error: ", err)
	//	c.Ctx.ResponseWriter.WriteHeader(406) //非法的用户ID
	//	return
	//}
	o := orm.NewOrm()
	//用创建的新卡号查询是否在数据库中存在
	if err := o.Read(&card); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(403) //卡片不存在
		return
	}
	//这里还没有具体设置user的id
	card.UserId = "01"
	c.Ctx.ResponseWriter.WriteHeader(200) //成功
	//传回这个卡片的具体信息
	c.Data["json"] = card
	c.ServeJSON()
}

//修改卡片的id和公司名——我们认为不需要这个
//ml
//@Title ModifyCardInfo
//@Description 修改卡片的卡号，公司
//@Param	id	query	string	true	原本的卡号
//@Success 200	{object} models.Card 	修改成功，返回新卡片对象
//@Failure 400	body解析错误
//@Failure 404	卡片信息读取错误
//@Failure 500	数据库更新操作错误
//@router  /card/:id/info	[put]
func (c *CardController) ModifyCardInfo() {
	oldCardId := c.Ctx.Input.Param(":id")
	body := c.Ctx.Input.RequestBody
	var newCard, oldCard models.Card
	oldCard.CardId = oldCardId
	//解析body
	if err := json.Unmarshal(body, &newCard); err != nil {
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	o := orm.NewOrm()
	//读取原卡片
	if err := o.Read(oldCard); err != nil {
		models.Log.Error("sql read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404)
		return
	}
	//读取新卡片
	if err := o.Read(newCard); err != nil {
		models.Log.Error("sql read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404)
		return
	}
	//增加新卡片中UserId关联,并取消原卡片的关联
	newCard.UserId = oldCard.UserId
	oldCard.UserId = ""

	_, err1 := o.Update(oldCard)
	if err1 != nil {
		models.Log.Error("update error: ", err1)
		c.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	_, err2 := o.Update(newCard)
	if err2 != nil {
		models.Log.Error("update error: ", err2)
		c.Ctx.ResponseWriter.WriteHeader(500)
		return
	}

	//修改成功，返回成功后的卡片对象
	c.Ctx.ResponseWriter.WriteHeader(200)
	c.Data["json"] = newCard
	c.ServeJSON()
}

//nfc扫码增加积分,兑换免费咖啡，前端传给我们1加积分 
//给前端说一下
//zjn
func (c *CardController) use_score() {}

//对优惠券的操作
//使用优惠卷
//前端返回给我优惠券对象的信息
//zyj
func (c *CardController) coupons() {}

//删除卡片 手动删除选项
func (c *CardController) Delete() {
	id := c.Ctx.Input.Param(":id")
	//fmt.Println(id)
	o := orm.NewOrm()
	card := models.Card{Id: id}
	if err := o.Read(&card); err == nil {
		count, _ := o.Delete(&card)
		if count == 0 {
			models.Log.Error("delete fail") //删除0个元素，即删除失败，返回状态码403
			c.Ctx.ResponseWriter.WriteHeader(403)
		} else {
			delCard := models.DelCard{CardId: card.Id, UserId: card.UserId, Remark: card.Remark}
			delCard.DelTime = time.Now()
			_, err := o.Insert(&delCard)
			if err != nil {
				models.Log.Error("Insert error: ", err) //被删卡插入垃圾箱失败
				c.Ctx.ResponseWriter.WriteHeader(403)
				return
			}
			c.Ctx.ResponseWriter.WriteHeader(200) //删除成功
		}
	} else {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(200) //card本就不存在，删除不存在的卡当作删除成功
	}
}

//GZH，修改备注
//@swagger注解配置
//@Title Put
//@Description edit cards' remark
//@Success 200
//@remark parameter is empty 400
//@Failure 403
//@router  /card/:id/remark  [put]
func (c *CardController) Put() {
	// 接收数据
	id := c.Ctx.Input.Param(":id")
	remark := c.GetString("remark")

	if remark == "" {
		// remark参数为空，设置400状态码
		models.Log.Error("param error: ", errors.New("illegal remark"))
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}

	o := orm.NewOrm()
	card := models.Card{Id: id}
	if err := o.Read(&card); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	// 查到了该记录，进行赋值
	card.Remark = remark
	// 更新记录
	if _, err := o.Update(&card); err != nil {
		// 更新失败
		models.Log.Error("update error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(403)
		return
	}
	// 成功,设置成功响应
	c.Ctx.ResponseWriter.WriteHeader(200)
	c.Data["json"] = card

	c.ServeJSON()
}

// swagger注解配置
// @Title Get
// @Param body query string true "enterprise_name"
// @Description get help message by the given enterprise_name
// @Success 200
// @Failure 404 read error
// @router /card/help/:Ename [get]
func (c *CardController) Help() {
	EName := c.Ctx.Input.Param(":Ename")
	fmt.Println(EName)
	o := orm.NewOrm()
	enterprise := models.Enterprise{Name: EName}
	// 查询记录
	if err := o.Read(&enterprise); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) // 查不到企业名对应的企业
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200) //成功
	c.Data["json"] = enterprise.HelpMsg
	c.ServeJSON()
}
