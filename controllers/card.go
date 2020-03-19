package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"
	"github.com/pkg/errors"
	"time"
	"strings"
	"strconv"
)

type CardController struct {
	beego.Controller
}

//查询指定card_id对应的卡片的所有信息
//zyc
func (c *CardController) Get_cardidinfo() {
	// 获取路由参数
	id := c.Ctx.Input.Param(":id")
	o := orm.NewOrm()
	card := models.Card{CardId: id}
	// 查询记录
	if err := o.Read(&card); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) // 查不到id对应的卡
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200) //成功
	//一张卡的信息可能不是card表就可以表出来的
	c.Data["json"] = card
	c.ServeJSON()
}

//添加卡片 在user表里添加此user和card的关联
//zjn
func (c *CardController) Post() {
	var card models.Card
	body := c.Ctx.Input.RequestBody
	if err := json.Unmarshal(body, &card); err != nil {
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	if err := card.CardParse(); err != nil {
		models.Log.Error("card parse error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(406) //非法的用户ID
		return
	}
	o := orm.NewOrm()
	if _, err := o.Insert(&card); err != nil {
		models.Log.Error("insert error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(403) //插入错误
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200) //成功
}

//修改卡片的id和公司名--我们认为不需要这个
//ml
func (c *CardController) modify_card(){}

//nfc扫码增加积分,兑换免费咖啡，前端传给我们1加积分 
//给前端说一下
//zjn
func (c *CardController) use_score(){}

//对优惠券的操作
//使用优惠卷
//前端返回给我优惠券对象的信息以及优惠券的信息
//zyj

// if err := o.Read(&card); err != nil {
// 	models.Log.Error("read error: ", err)
// 	c.Ctx.ResponseWriter.WriteHeader(404) // 查不到id对应的卡
// 	return
// }
func (c *CardController) coupons() {
	var info struct{CardID string;CouponsID string;Increment int}
	var card models.Card
	body := c.Ctx.Input.RequestBody
	if err:= json.Unmarshal(body,&body); err != nil{
		models.Log.Error("unmarshal error：", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return 
	}
	card.CardId = info.CardID
	o := orm.NewOrm()
	if err:= o.Read(&card); err != nil{
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) // 查不到id对应的卡
		return
	}
	couponsList := strings.Split(card.CouponsList," ")
	couponsNumList := strings.Split(card.CouponsNum," ")
	for i ,value  := range couponsList{
		if value==info.CouponsID{
			var temp int
			if temp,err := strconv.Atoi(couponsNumList[i]);err!=nil{
				models.Log.Error("invalid data: ",err)
				c.Ctx.ResponseWriter.WriteHeader(406) //非法数据
				return 
			}
			temp += info.Increment
			couponsNumList[i] = strconv.Itoa(temp)
		}
	}
	newCouponsNum := strings.Join(couponsNumList," ")
	card.CouponsNum = newCouponsNum
	if _ , err := o.Update(&card);err!=nil{
		models.Log.Error("invalid data: ",err)
		c.Ctx.ResponseWriter.WriteHeader(404) //查找不到相应的id卡进行数据更新
		return 
	}
}

//生成新的代表优惠券数量的字符串
func makeCouponsNum(card *models.Card,couponsID string,increment int) {

}


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
