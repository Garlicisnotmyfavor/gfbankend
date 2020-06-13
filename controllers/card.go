package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"
	"github.com/gfbankend/utils"
	_ "github.com/pkg/errors"
	"strconv"
	"time"
	_ "unsafe"
)

type CardController struct {
	beego.Controller
}

/*
func (c *CardController) Get_cardidinfo() {
	//声明一个结构体
	var jsonValue models.CardInfo
	//取得卡的id
	//根据卡的id从数据库中读取信息存储进一个临时结构体temp中
	//从card结构体中的ScoreList，ScoreNum，CouponsList，CouponsNum 用strings.Split()进行分割
	  生成四个数组scoreList，scoreNum，couponsList，couponsNum，把四个数组传进jsonValue
	//for index,value := range scoreList{ //从数据库中寻找value对应的方法 //把那一列的数据添加到jsonValue的ScoreDetails中 }
	//for index,value := range CouponsList {//从数据库中寻找value对应的方法 //数据添加到jsonValue的CouponsDetails中}
	//可以直接append
}
*/

// 查询指定card_id对应的卡片的所有信息
// zjn
// @Title GetCardIDInfo
// @Description 将这张卡片的所有信息传出去
// @Param	id	path	string	true	查询的卡号
// @Success 200	查询成功
// @Failure 400	查询不到对应的卡
// @Failure 401	没处于登录状态，无权限
// @Failure 404	查询不到对应的公司
// @router /card/:id [get]
func (c *CardController) GetCardIDInfo() {
	sess := c.GetSession("userInfo")
	// 由cookie 得不到session说明没登录，无权限
	if sess == nil {
		models.Log.Error("not login: ")
		c.Ctx.ResponseWriter.WriteHeader(401)
		return
	}
	// 获取路由参数
	id := c.Ctx.Input.Param(":id")
	o := orm.NewOrm()
	//设置一个填充了cardid的card结构
	var card models.Card
	var ep models.Enterprise
	card.CardId = id
	// 查询记录
	if err := o.Read(&card); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) // 查不到id对应的卡
		return
	}
	//找到卡后要去找对应的公司的信息
	ep.Name = card.Enterprise
	if err := o.Read(&ep, "name"); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) // 查不到公司的信息
		return
	}
	var cardenter struct {
		Card       models.Card
		Enterprise models.Enterprise
	}
	fmt.Println(card)
	fmt.Println(ep)
	cardenter.Card = card
	cardenter.Enterprise = ep
	////若查到card这一列后，需要找到它的卡的积分或卷的规则
	//var CouponsDetails models.Coupons
	//var ScoreDetails   models.Score
	//ScoreList := strings.Split(card.ScoreList, " ")
	//CouponsList := strings.Split(card.CouponsList, " ")
	//for i, value := range ScoreList {
	//	ScoreDetails.ScoreID = value
	//	if err := o.Read(&ScoreDetails); err != nil {
	//		models.Log.Error("not exist error: ", err)
	//		c.Ctx.ResponseWriter.WriteHeader(403) //找不到这个类型
	//		i = i+1 //尽量修改不用这种方式
	//		i = i-1
	//		return
	//	}
	//	cardinfo.ScoreDetails = append(cardinfo.ScoreDetails, ScoreDetails)
	//}
	//for i, value := range CouponsList {
	//	CouponsDetails.CouponsID = value
	//	if err := o.Read(&CouponsDetails); err != nil {
	//		models.Log.Error("not exist error: ", err)
	//		c.Ctx.ResponseWriter.WriteHeader(403) //找不到这个类型
	//		i = i+1 //尽量修改不用这种方式
	//		i = i-1
	//		return
	//	}
	//	cardinfo.CouponsDetails = append(cardinfo.CouponsDetails, CouponsDetails)
	//}
	c.Ctx.ResponseWriter.WriteHeader(200) //成功
	c.Data["json"] = cardenter
	c.ServeJSON()
}

//添加卡片 在user表里添加此user和card的关联
//userid从cookie，session取得
//zjn
//@Title AddCard
//@Description 将这个user的id和卡绑定,由cookie获取sessionid从而得到当前用户ID
//@Param	id	body	/	true	原本的卡号cardid+企业enterprise
//@Success 200	{object} models.Card 	返回绑定的卡的大致信息
//@Failure 403	绑定的卡片不存在
//@Failure 400	解析错误
//@Failure 401	没处于登录状态，无权限
//@Failure 402	数据不匹配
//@router  /card/add [post]
func (c *CardController) AddCard() {
	sess := c.GetSession("userInfo")
	// 由cookie 得不到session说明没登录，无权限
	if sess == nil {
		models.Log.Error("not login: ")
		c.Ctx.ResponseWriter.WriteHeader(401)
		return
	}
	user := sess.(models.User)
	userId := user.Id
	//这里没有对比enterprise和cardid
	var addinfo struct {
		CardID     string
		Enterprise string
	}
	body := c.Ctx.Input.RequestBody
	//解析body
	if err := json.Unmarshal(body, &addinfo); err != nil {
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	//目前的逻辑不需要解析函数
	//if err := card.CardParse(); err != nil {
	//	models.Log.Error("card parse error: ", err)
	//	c.Ctx.ResponseWriter.WriteHeader(406) //非法的用户ID
	//	return
	//}
	o := orm.NewOrm()
	card := models.Card{}
	card.CardId = addinfo.CardID
	//用创建的新卡号查询是否在数据库中存在
	if err := o.Read(&card); err != nil {
		// models.Log.Error("not exist error: ", err)
		// c.Ctx.ResponseWriter.WriteHeader(404) //卡片不存在
		card.Enterprise = addinfo.Enterprise
		card.CardType = "Integrate"
		card.UserId = userId
		card.State = "empty"
		card.City = "empty"
		if _,err := o.Insert(&card); err!= nil {
			models.Log.Error("insert database error: %s",err)
			c.Ctx.ResponseWriter.WriteHeader(405) //数据库更新失败
			return
		}
		c.Data["json"] = card
		c.ServeJSON()
		return
	}
	if len(card.UserId) != 0 {
		models.Log.Error("card already bind")
		c.Ctx.ResponseWriter.WriteHeader(403) //卡片另有主人
		return
	}
	if card.Enterprise != addinfo.Enterprise {
		models.Log.Error("wrong cardID")
		c.Ctx.ResponseWriter.WriteHeader(402) //输入id和公司名不匹配
		return
	}
	//匹配后建立关联
	card.UserId = userId
	if _, err := o.Update(&card); err != nil {
		models.Log.Error("update database error")
		c.Ctx.ResponseWriter.WriteHeader(405) //数据库更新失败
		return
	}
	// card.UserId = &models.User{Id:"2018091620000"}
	//card.UserId = addinfo.Enterprise
	c.Ctx.ResponseWriter.WriteHeader(200) //成功
	//传回这个卡片的具体信息
	c.Data["json"] = card
	c.ServeJSON()
}


// lj
// @Title getCard
// @Description 领取新的卡片
// @Param CardID body string true 卡号
// @Param Enterprise body string true 公司名称
// @Param CardType body string true 卡片类型
// @Success 200 
// @Failure 400 解析Json出错
// @Failure 401	没处于登录状态，无权限
// @Failure 403	卡片已经被其他用户领取
// @Failure 405	数据库更新失败
// @router  /card/get/:demoId [put]
// 用户领取某种类型的卡片，前端给卡号，卡的类型，企业名称
func (c *CardController) GetCard() {
	sess := c.GetSession("userInfo")
	// 由cookie 得不到session说明没登录，无权限
	if sess == nil {
		models.Log.Error("not login: ")
		c.Ctx.ResponseWriter.WriteHeader(401)
		return
	}
	demoId := c.Ctx.Input.Param(":demoId")
	user := sess.(models.User)
	userId := user.Id
	////定义卡的信息结构体
	//var CardInfo struct {
	//	CardDemoId string `json:"cardDemoId"`
	//}
	//body := c.Ctx.Input.RequestBody
	////解析body
	//if err := json.Unmarshal(body, &CardInfo); err != nil {
	//	models.Log.Error("unmarshal error: ", err)
	//	c.Ctx.ResponseWriter.WriteHeader(400)
	//	return
	//}
	o := orm.NewOrm()
	var demo models.CardDemo
	var card models.Card
	if ID,err := strconv.Atoi(demoId); err==nil {
		demo.Id = ID
	} else {
		models.Log.Error("atoi error: ")
		c.Ctx.ResponseWriter.WriteHeader(402)
		return
	}
	if err := o.Read(&demo); err != nil {
		models.Log.Error("sql read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404)
		return
	}
	//将用户ID与card信息相关联
	card.UserId = userId
	card.Coupons = demo.Coupons
	card.Enterprise = demo.Enterprise
	for {
		id := util.RandStr(13)
		if err := o.Read(&models.Card{CardId: id}) ; err != nil {
			card.CardId = id
			break
		}
	}
	if _, err := o.Insert(&card); err != nil {
		models.Log.Error("insert card error")
		c.Ctx.ResponseWriter.WriteHeader(405) //数据库更新失败
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200) //成功
	//传回这个卡片的具体信息
	c.Data["json"] = card
	c.ServeJSON()
}



//修改卡片的id和公司名--我们认为不需要这个
//ml
//@Title ModifyCardInfo
//@Description 修改卡片的卡号，公司
//@Param	id	path	string	true	原本的卡号
//@Param	cardInfo	body	/	true	新卡信息   CardId(string)+Enterprise(string)
//@Success 200	{object} models.Card 	修改成功，返回新卡片对象
//@Failure 400	body解析错误
//@Failure 401	没处于登录状态，无权限
//@Failure 403	卡号解析错误
//@Failure 404	卡片信息读取错误
//@Failure 500	数据库更新操作错误
//@router  /card/:id/info [put]
func (c *CardController) ModifyCardInfo() {
	sess := c.GetSession("userInfo")
	// 由cookie 得不到session说明没登录，无权限
	if sess == nil {
		models.Log.Error("not login: ")
		c.Ctx.ResponseWriter.WriteHeader(401)
		return
	}
	oldCardId := c.Ctx.Input.Param(":id")
	body := c.Ctx.Input.RequestBody
	var newCard models.Card
	oldCard := models.Card{CardId: oldCardId}
	o := orm.NewOrm()
	//读取原卡片
	if err := o.Read(&oldCard); err != nil {
		models.Log.Error("sql read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404)
		return
	}
	newCard = oldCard
	//解析body
	if err := json.Unmarshal(body, &newCard); err != nil {
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	////读取新卡片
	//if err := o.Read(&newCard); err != nil {
	//	models.Log.Error("sql read error: ", err)
	//	c.Ctx.ResponseWriter.WriteHeader(404)
	//	return
	//}
	////新卡片另有主人
	//if newCard.UserId != "" && newCard.UserId != oldCard.UserId {
	//	models.Log.Error("sql update error: card already have owner")
	//	c.Ctx.ResponseWriter.WriteHeader(409)
	//	return
	//}
	// if err := newCard.CardParse(); err != nil {
	// 	models.Log.Error("parse error: ", err)
	// 	c.Ctx.ResponseWriter.WriteHeader(403)
	// 	return
	// }
	// 增加新卡片中UserId关联,并取消原卡片的关联
	// oldCard.UserId = newCard.UserId
	if _, err := o.Delete(&oldCard); err != nil {
		models.Log.Error("delete oldCard error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	newCard.UserId = oldCard.UserId
	if _, err := o.Insert(&newCard); err != nil {
		models.Log.Error("insert newCard error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(406)
		return
	}
	//if _, err := o.Update(&newCard); err != nil {
	//	models.Log.Error("update error: ", err)
	//	c.Ctx.ResponseWriter.WriteHeader(500)
	//	return
	//}
	//修改成功，返回成功后的卡片对象
	c.Ctx.ResponseWriter.WriteHeader(200)
	c.Data["json"] = newCard
	c.ServeJSON()
}

//nfc扫码增加积分,兑换免费咖啡，前端传给我们1加积分
//给前端说一下
//ml
//@Title UseScore
//@Description 操作卡片的积分
//@Param id body / true CardId(string)+increment(int)
//@Success 200	{object} models.Card 	修改成功，返回新卡片对象
//@Failure 400	body解析错误
//@Failure 401	没处于登录状态，无权限
//@Failure 406	积分信息有误
//@Failure 500	数据库更新操作错误
//@router /card/:id/score [put]
func (c *CardController) UseScore() {
	sess := c.GetSession("userInfo")
	// 由cookie 得不到session说明没登录，无权限
	if sess == nil {
		models.Log.Error("not login: ")
		c.Ctx.ResponseWriter.WriteHeader(401)
		return
	}
	var ScoreInfo struct {
		CardId    string
		Increment int
	}
	ScoreInfo.CardId = c.Ctx.Input.Param(":id")
	body := c.Ctx.Input.RequestBody
	//解析请求体
	if err := json.Unmarshal(body, &ScoreInfo); err != nil {
		models.Log.Error("unmarshal error：", err)
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	card := models.Card{CardId: ScoreInfo.CardId}
	o := orm.NewOrm()
	//根据CardId读取完整的卡片信息
	if err := o.Read(&card); err != nil {
		models.Log.Error("sql read error：", err)
		c.Ctx.ResponseWriter.WriteHeader(404)
		return
	}
	//ScoreList := strings.Split(Card.ScoreList, " ")
	//ScoreNumList := strings.Split(Card.ScoreNum, " ")
	//for i, v := range ScoreList {
	//	if v == ScoreInfo.ScoreId {
	//		tmp, err := strconv.Atoi(ScoreNumList[i])
	//		if err != nil {
	//			models.Log.Error("invalid data：", err)
	//			c.Ctx.ResponseWriter.WriteHeader(406)
	//			return
	//		}
	//		tmp += ScoreInfo.increment
	//		ScoreNumList[i] = strconv.Itoa(tmp)
	//		Card.ScoreNum = strings.Join(ScoreNumList, " ")
	//		hasIncrease = true
	//	}
	//}
	card.Score += ScoreInfo.Increment
	if _, err := o.Update(&card); err != nil {
		models.Log.Error("sql update error：", err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	//	成功更新对应的卡片积分数量
	c.Ctx.ResponseWriter.WriteHeader(200)
	//	返回更新积分后的卡片
	c.Data["json"] = card
	c.ServeJSON()
}

//对优惠券的操作
//使用优惠卷
//前端返回给我优惠券对象的信息以及优惠券的信息
//增加或减少某张卡的某种优惠券
//zyj
//@Title coupons
//@Description 增加或减少某张卡的某种优惠券
//@Param id path string true 卡号
//@Param Increment body int true  优惠券改变的数量，可以为负数
//@Success 200  成功
//@Failure 400	json解析错误
//@Failure 401	没处于登录状态，无权限
//@Failure 403	优惠券不足
//@Failure 404	卡不存在
//@Failure 406	非法数据
//@router  /card/:id/coupons [post]
func (c *CardController) Coupons() {
	sess := c.GetSession("userInfo")
	// 由cookie 得不到session说明没登录，无权限
	if sess == nil {
		models.Log.Error("not login: ")
		c.Ctx.ResponseWriter.WriteHeader(401)
		return
	}
	CardId := c.Ctx.Input.Param(":id")
	o := orm.NewOrm()
	//o.Insert(&models.Card{CardId: "1234567890123456", UserId: "1234567890124", CardType: "MembershipCard", Enterprise: "StarBuck", State: "Sichuan", City: "Chengdu", Money: 100, ExpireTime: time.Now()})
	var increment struct{ Value int }
	card := models.Card{CardId: CardId}
	body := c.Ctx.Input.RequestBody
	//解析请求体
	if err := json.Unmarshal(body, &increment); err != nil {
		models.Log.Error("unmarshal error：", err)
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	if err := o.Read(&card); err != nil {
		models.Log.Error("can't find card: ", err)
		c.Ctx.ResponseWriter.WriteHeader(403) //查找不到相应的id卡进行数据更新
		return
	}
	card.CouponsNum += increment.Value
	fmt.Println(card)
	if _, err := o.Update(&card); err != nil {
		models.Log.Error("can't update card: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) //查找不到相应的id卡进行数据更新
		return
	}
	//返回成功的信息
	c.Ctx.ResponseWriter.WriteHeader(200)
	c.Data["json"] = card
	c.ServeJSON()
}

//zyj
//@Title delete
//@Description 删除卡片
//@Param id path string true 卡号
//@Success 200
//@Failure 400 json解析错误
//@Failure 401 没登录，无权限
//@Failure 403 用户ID不存在
//@Failure 404 卡不存在
//@router  /card/:id/delete [post]
func (c *CardController) Delete() {
	sess := c.GetSession("userInfo")
	// 由cookie 得不到session说明没登录，无权限
	if sess == nil {
		models.Log.Error("not login: ")
		c.Ctx.ResponseWriter.WriteHeader(401)
		return
	}
	id := c.Ctx.Input.Param(":id")
	o := orm.NewOrm()
	card := models.Card{CardId: id}
	if err := o.Read(&card); err != nil {
		models.Log.Error("can't find card: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) //查找不到相应的id卡进行数据更新
		return
	}
	card.DelTime = time.Now()
	if _, err := o.Update(&card); err != nil {
		models.Log.Error("can't update card: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) //查找不到相应的id卡进行数据更新
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200)
	c.Data["json"] = card
	c.ServeJSON()
	return
}

//author: lj
//@Title 查看卡片使用记录
//@Description
//@Param id query string true 卡号
//@Success 200
//@Failure 400 解析Json失败
//@Failure 401 没有登录
//@Failure 404 卡不存在
//@router  /card/:id/cardLog [get]
func (c *CardController) CardLog() {
	sess := c.GetSession("userInfo")
	// 由cookie 得不到session说明没登录，无权限
	if sess == nil {
		models.Log.Error("not login: ")
		c.Ctx.ResponseWriter.WriteHeader(401)
		return
	}
	// 获取路由参数
	id := c.Ctx.Input.Param(":id")
	o := orm.NewOrm()
	//设置一个填充了CardID的CardLog结构
	var cardLog models.CardLog
	cardLog.CardId = id
	body := c.Ctx.Input.RequestBody
	//解析请求体
	if err := json.Unmarshal(body, &cardLog); err != nil {
		models.Log.Error("unmarshal error：", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析Json失败
		return
	}
	// 查询记录
	if err := o.Read(&cardLog); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) // 查不到id对应的卡
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200)
	c.Data["json"] = cardLog
	c.ServeJSON()
	return
}

//author: zyj
//@Title  获取虚拟卡片
//@Description
//@Param id query string true 卡号
//@Success 200
//@Failure
//@router  /card/getnewcard [post]
func (c *CardController) GetNewCard() {

}

// func isIdInUsers(id string) bool {
// 	o := orm.NewOrm()
// 	user := models.User{Id:id}
// 	if err := o.Read(&user); err != nil {
// 		return false;
// 	}
// 	return true
//}

//GZH，修改备注
//@swagger注解配置
//@Title Put
//@Description edit cards' remark
//@Success 200
//@remark parameter is empty 400
//@Failure 403
//@router  /card/:id/remark  [put]
// func (c *CardController) Put() {
// 	// 接收数据
// 	id := c.Ctx.Input.Param(":id")
// 	remark := c.GetString("remark")

// 	if remark == "" {
// 		// remark参数为空，设置400状态码
// 		models.Log.Error("param error: ", errors.New("illegal remark"))
// 		c.Ctx.ResponseWriter.WriteHeader(400)
// 		return
// 	}
// 	o := orm.NewOrm()
// 	//读取原卡片
// 	if err := o.Read(oldCard); err != nil {
// 		models.Log.Error("sql read error: ", err)
// 		c.Ctx.ResponseWriter.WriteHeader(404)
// 		return
// 	}
// 	//读取新卡片
// 	if err := o.Read(newCard); err != nil {
// 		models.Log.Error("sql read error: ", err)
// 		c.Ctx.ResponseWriter.WriteHeader(404)
// 		return
// 	}
// 	//增加新卡片中UserId关联,并取消原卡片的关联
// 	newCard.UserId = oldCard.UserId
// 	oldCard.UserId = ""

// 	_, err1 := o.Update(oldCard)
// 	if err1 != nil {
// 		models.Log.Error("update error: ", err1)
// 		c.Ctx.ResponseWriter.WriteHeader(500)
// 		return
// 	}
// 	_, err2 := o.Update(newCard)
// 	if err2 != nil {
// 		models.Log.Error("update error: ", err2)
// 		c.Ctx.ResponseWriter.WriteHeader(500)
// 		return
// 	}

// 	//修改成功，返回成功后的卡片对象
// 	c.Ctx.ResponseWriter.WriteHeader(200)
// 	c.Data["json"] = newCard
// 	c.ServeJSON()
// }

//nfc扫码增加积分,兑换免费咖啡，前端传给我们1加积分
//给前端说一下
//zjn

//func (c *CardController) use_score() {}

//删除卡片 手动删除选项
//func (c *CardController) Delete() {
//	id := c.Ctx.Input.Param(":id")
//	//fmt.Println(id)
//	o := orm.NewOrm()
//	card := models.Card{Id: id}
//	if err := o.Read(&card); err == nil {
//		count, _ := o.Delete(&card)
//		if count == 0 {
//			models.Log.Error("delete fail") //删除0个元素，即删除失败，返回状态码403
//			c.Ctx.ResponseWriter.WriteHeader(403)
//		} else {
//			delCard := models.DelCard{CardId: card.Id, UserId: card.UserId, Remark: card.Remark}
//			delCard.DelTime = time.Now()
//			_, err := o.Insert(&delCard)
//			if err != nil {
//				models.Log.Error("Insert error: ", err) //被删卡插入垃圾箱失败
//				c.Ctx.ResponseWriter.WriteHeader(403)
//				return
//			}
//			c.Ctx.ResponseWriter.WriteHeader(200) //删除成功
//		}
//	} else {
//		models.Log.Error("read error: ", err)
//		c.Ctx.ResponseWriter.WriteHeader(200) //card本就不存在，删除不存在的卡当作删除成功
//	}
//}
//
////GZH，修改备注
////@swagger注解配置
////@Title Put
////@Description edit cards' remark
////@Success 200
////@remark parameter is empty 400
////@Failure 403
////@router  /card/:id/remark  [put]
//func (c *CardController) Put() {
//	// 接收数据
//	id := c.Ctx.Input.Param(":id")
//	remark := c.GetString("remark")
//
//	if remark == "" {
//		// remark参数为空，设置400状态码
//		models.Log.Error("param error: ", errors.New("illegal remark"))
//		c.Ctx.ResponseWriter.WriteHeader(400)
//		return
//	}
//
//	o := orm.NewOrm()
//	card := models.Card{Id: id}
//	if err := o.Read(&card); err != nil {
//		models.Log.Error("read error: ", err)
//		c.Ctx.ResponseWriter.WriteHeader(400)
//		return
//	}
//	// 查到了该记录，进行赋值
//	card.Remark = remark
//	// 更新记录
//	if _, err := o.Update(&card); err != nil {
//		// 更新失败
//		models.Log.Error("update error: ", err)
//		c.Ctx.ResponseWriter.WriteHeader(403)
//		return
//	}
//	// 成功,设置成功响应
//	c.Ctx.ResponseWriter.WriteHeader(200)
//	c.Data["json"] = card
//
//	c.ServeJSON()
//}
//
//// swagger注解配置
//// @Title Get
//// @Param body query string true "enterprise_name"
//// @Description get help message by the given enterprise_name
//// @Success 200
//// @Failure 404 read error
//// @router /card/help/:Ename [get]
//func (c *CardController) Help() {
//	EName := c.Ctx.Input.Param(":Ename")
//	fmt.Println(EName)
//	o := orm.NewOrm()
//	enterprise := models.Enterprise{Name: EName}
//	// 查询记录
//	if err := o.Read(&enterprise); err != nil {
//		models.Log.Error("read error: ", err)
//		c.Ctx.ResponseWriter.WriteHeader(404) // 查不到企业名对应的企业
//		return
//	}
//	c.Ctx.ResponseWriter.WriteHeader(200) //成功
//	c.Data["json"] = enterprise.HelpMsg
//	c.ServeJSON()
//}
//
