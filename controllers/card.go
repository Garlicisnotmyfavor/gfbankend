package controllers

import (
	"encoding/json"
	_"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gfbankend/models"
	_ "github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
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
	var card models.Card
	var cardinfo models.CardInfo
	card.CardId = id
	// 查询记录
	if err := o.Read(&card); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400) // 查不到id对应的卡
		return
	}
	//若查到card这一列后，需要找到它的卡的积分或卷的规则
	var CouponsDetails models.Coupons
	var ScoreDetails   models.Score
	ScoreList := strings.Split(card.ScoreList, " ")
	CouponsList := strings.Split(card.CouponsList, " ")
	for i, value := range ScoreList {
		ScoreDetails.ScoreID = value
		if err := o.Read(&ScoreDetails); err != nil {
			models.Log.Error("not exist error: ", err)
			c.Ctx.ResponseWriter.WriteHeader(403) //找不到这个类型
			i = i+1 //尽量修改不用这种方式
			i = i-1
			return
		}
		cardinfo.ScoreDetails = append(cardinfo.ScoreDetails, ScoreDetails)
	}
	for i, value := range CouponsList {
		CouponsDetails.CouponsID = value
		if err := o.Read(&CouponsDetails); err != nil {
			models.Log.Error("not exist error: ", err)
			c.Ctx.ResponseWriter.WriteHeader(403) //找不到这个类型
			i = i+1 //尽量修改不用这种方式
			i = i-1
			return
		}
		cardinfo.CouponsDetails = append(cardinfo.CouponsDetails, CouponsDetails)
	}

	//整合到一个struct里
	cardinfo.CardId = card.CardId      
	cardinfo.UserId = card.UserId      
	cardinfo.CouponsList = card.CouponsList
	cardinfo.CardType = card.CardType     
	cardinfo.Enterprise = card.Enterprise   
	cardinfo.State = card.State       
	cardinfo.City = card.City         
	cardinfo.Money = card.Money        
	cardinfo.ScoreNum = card.ScoreNum     
	cardinfo.ScoreList = card.ScoreList    
	cardinfo.CouponsNum = card.CouponsNum   
	cardinfo.ExpireTime =  card.ExpireTime
	cardinfo.DelTime = card.DelTime     
	cardinfo.CardOrder = card.CardOrder   
	cardinfo.FactoryNum = card.FactoryNum   
	cardinfo.BatchNum = card.BatchNum    
	cardinfo.SerialNum = card.SerialNum    
	//cardinfo.CouponsDetails = CouponsDetails
	//cardinfo.ScoreDetails = ScoreDetails

	c.Ctx.ResponseWriter.WriteHeader(200) //成功
	c.Data["json"] = cardinfo
	c.ServeJSON()
	
}

//添加卡片 在user表里添加此user和card的关联
//zjn
//@Title AddCard
//@Description 将这个user的id和卡绑定
//@Param	id	query	string	true	原本的卡号
//@Success 200	{object} models.Card 	返回绑定的卡的大致信息
//@Failure 403	绑定的卡片不存在
//@Failure 400	解析错误
//@Failure 402	数据不匹配
//@router  /card/add [post]
func (c *CardController) AddCard() {
	//这里没有对比enterprise和cardid
	var addinfo struct {
		CardId    string
		Enterprise	string
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
	card.CardId = addinfo.CardId
	//用创建的新卡号查询是否在数据库中存在
	if err := o.Read(&card); err != nil {
		models.Log.Error("not exist error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(403) //卡片不存在
		return
	}
	if card.Enterprise != addinfo.Enterprise {
		c.Ctx.ResponseWriter.WriteHeader(402) //输入id和公司名不匹配
		return
	}
	//匹配后建立关联
	//这里还没有具体设置user的id
	card.UserId = "0000000000000"
	//card.UserId = addinfo.Enterprise
	c.Ctx.ResponseWriter.WriteHeader(200) //成功
	//传回这个卡片的具体信息
	c.Data["json"] = card
	c.ServeJSON()
}

//修改卡片的id和公司名--我们认为不需要这个
//ml
//@Title ModifyCardInfo
//@Description 修改卡片的卡号，公司
//@Param	id	query	string	true	原本的卡号
//@Param	cardInfo	body	/	true	新卡信息(卡号CardId+公司Enterprise)
//@Success 200	{object} models.Card 	修改成功，返回新卡片对象
//@Failure 400	body解析错误
//@Failure 404	卡片信息读取错误
//@Failure 500	数据库更新操作错误
//@router  /card/:id/info [put]
func (c *CardController) ModifyCardInfo() {
	oldCardId := c.Ctx.Input.Param(":id")
	body := c.Ctx.Input.RequestBody
	var newCard models.Card
	oldCard := models.Card{CardId: oldCardId}
	//解析body
	if err := json.Unmarshal(body, &newCard); err != nil {
		models.Log.Error("unmarshal error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(400)
		return
	}
	o := orm.NewOrm()
	//读取原卡片
	if err := o.Read(&oldCard); err != nil {
		models.Log.Error("sql read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404)
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
	//增加新卡片中UserId关联,并取消原卡片的关联
	newCard.UserId = oldCard.UserId
	if _, err := o.Insert(&newCard); err != nil {
		models.Log.Error("insert error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(406)
		return
	}
	oldCard.UserId = ""
	if _, err := o.Update(&oldCard); err != nil {
		models.Log.Error("update error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(500)
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
//@Failure 406	积分信息有误
//@Failure 500	数据库更新操作错误
//@router /card/:id/score [put]
func (c *CardController) UseScore() {
	var ScoreInfo struct {
		CardId    string
		Increment int
	}
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
	oldScore, ok := strconv.Atoi(card.Score)
	if ok != nil {
		models.Log.Error("score update error ")
		c.Ctx.ResponseWriter.WriteHeader(406)
		return
	}
	card.Score  = strconv.Itoa(oldScore + ScoreInfo.Increment)
	//	卡片更新错误，可能是数据库出错
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
//@Param id query
//@Param CouponsID,Increment Body / true id(string)+CouponsID(string)+increment(int)  卡号，优惠券号，增量(zeng'l增量可以为负数)
//@Success 200  成功
//@Failure 400/403/404/406	json解析错误/优惠券不足/卡不存在/非法数据
//@router  /card/:id/coupons [post]
func (c *CardController) Coupons() {
	CardId := c.Ctx.Input.Param(":id")
	var info struct{CardID string `json:"-"`;CouponsID string;Increment int}
	info.CardID = CardId
	var card models.Card
	body := c.Ctx.Input.RequestBody
	if err:= json.Unmarshal(body,&info); err != nil{
		models.Log.Error("unmarshal error：", err)
		c.Ctx.ResponseWriter.WriteHeader(400) //解析json错误
		return
	}
	card.CardId = info.CardID
	o := orm.NewOrm()
	if err := o.Read(&card); err != nil {
		models.Log.Error("read error: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) // 查不到id对应的卡
		return
	}
	couponsList := strings.Split(card.CouponsList, " ")
	couponsNumList := strings.Split(card.CouponsNum, " ")
	for i, value := range couponsList {
		if value == info.CouponsID {
			var temp int
			temp, err := strconv.Atoi(couponsNumList[i])
			if err != nil {
				models.Log.Error("invalid data: ", err)
				c.Ctx.ResponseWriter.WriteHeader(406) //非法数据
				return
			}
			temp += info.Increment
			if temp<0 {
				models.Log.Error("not enough coupons")
				c.Ctx.ResponseWriter.WriteHeader(403) //优惠券不足
				return
			}
			couponsNumList[i] = strconv.Itoa(temp)
		}
	}
	newCouponsNum := strings.Join(couponsNumList, " ")
	card.CouponsNum = newCouponsNum
	if _, err := o.Update(&card); err != nil {
		models.Log.Error("can't update card: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) //查找不到相应的id卡进行数据更新
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200) //成功
	// fmt.Println(card)  //用于postman测试，上线工作后记得注释掉
	c.Data["json"] = card
	c.ServeJSON()
}

//zyj
//@Title delete
//@Description 删除卡片
//@Param id query string true 卡号
//@Success 200
//@Failure 400/404	json解析错误/卡不存在
//@router  /card/:id/delete [post]
func (c *CardController) Delete() {
	id := c.Ctx.Input.Param(":id")
	o := orm.NewOrm()
	card := models.Card{CardId: id}
	if err := o.Read(&card); err != nil {	
		models.Log.Error("can't find card: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) //查找不到相应的id卡进行数据更新
		return
	} 
	card.DelTime = time.Now()
	if _,err := o.Update(&card); err != nil{
		models.Log.Error("can't update card: ", err)
		c.Ctx.ResponseWriter.WriteHeader(404) //查找不到相应的id卡进行数据更新
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(200)
	c.Data["json"] = card
	c.ServeJSON()
	return
}

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
